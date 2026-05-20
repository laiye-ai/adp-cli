package cmd

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/laiye-ai/adp-cli/internal/api"
	"github.com/laiye-ai/adp-cli/internal/config"
	"github.com/laiye-ai/adp-cli/internal/errors"
	"github.com/laiye-ai/adp-cli/internal/file"
	"github.com/laiye-ai/adp-cli/internal/i18n"
	"github.com/laiye-ai/adp-cli/internal/telemetry"
)

// checkAPIResponse inspects the response body for a non-success "code" field
func checkAPIResponse(result map[string]interface{}, context string) {
	if result == nil {
		return
	}
	code, ok := result["code"].(string)
	if !ok || code == "" || code == "success" {
		return
	}
	msg, _ := result["message"].(string)
	if msg == "" {
		msg = code
	}
	cliErr := errors.NewCLIError(
		fmt.Sprintf("API error (%s): %s", code, msg),
		errors.ErrorTypeAPI, errors.ExitGeneralError, false,
		"", result,
	)
	telemetry.SetError(cliErr)
	formatterOut.ExitWithError(cliErr)
}

// BatchJob represents a single job in a batch
type BatchJob struct {
	Index int
	Name  string // display name, also used for result file naming
}

// BatchResult represents the result of processing a single job
type BatchResult struct {
	Input  string                 `json:"input"`
	Status string                 `json:"status"` // "success" or "failed"
	Data   map[string]interface{} `json:"data,omitempty"`
	Error  string                 `json:"error,omitempty"`
}

// BatchSummary represents the summary of a batch processing run
type BatchSummary struct {
	Total     int                    `json:"total"`
	Success   int                    `json:"success"`
	Failed    int                    `json:"failed"`
	OutputDir string                 `json:"output_dir,omitempty"`
	Files     []BatchSummaryFileItem `json:"files"`
}

// BatchSummaryFileItem represents a single file entry in the summary
type BatchSummaryFileItem struct {
	Input  string `json:"input"`
	Output string `json:"output"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// ProcessFunc is the function signature for processing a single job
type ProcessFunc func(job BatchJob) (map[string]interface{}, error)

// SubmitFunc is the function signature for submitting a single async job (returns taskID)
type SubmitFunc func(job BatchJob) (string, error)

// SubmitResult represents the result of submitting a single async job
type SubmitResult struct {
	Path   string `json:"path"`
	TaskID string `json:"task_id,omitempty"`
	Error  string `json:"error,omitempty"`
}

// batchProcess runs jobs concurrently with retry support and writes results to individual files
func batchProcess(jobs []BatchJob, concurrency, maxRetry int, processFn ProcessFunc, exportPath string) []BatchResult {
	total := len(jobs)
	results := make([]BatchResult, total)

	formatterOut.PrintInfo(fmt.Sprintf(i18n.T("processing_files"), total))

	if concurrency < 1 {
		concurrency = 1
	}

	isBatch := total > 1

	// For batch mode, prepare output directory upfront
	var outputDir string
	var fileHandler *file.FileHandler
	if isBatch {
		fileHandler = file.New()
		outputDir = exportPath
		if outputDir == "" {
			outputDir = fmt.Sprintf("adp_results_%s", time.Now().Format("20060102_150405"))
		}
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			formatterOut.PrintError(fmt.Sprintf("Failed to create output directory: %s", err.Error()))
			// Fall back to in-memory only; outputBatchSummary will print to stdout
			outputDir = ""
		} else {
			formatterOut.PrintSuccess(fmt.Sprintf(i18n.T("results_exported_to"), absPath(outputDir)))
		}
	}

	jobCh := make(chan BatchJob, total)
	var wg sync.WaitGroup

	for w := 0; w < concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobCh {
				formatterOut.PrintProgress(job.Index+1, total, job.Name)

				data, err := retryWithBackoff(func() (map[string]interface{}, error) {
					return processFn(job)
				}, maxRetry)

				if err != nil {
					results[job.Index] = BatchResult{
						Input:  job.Name,
						Status: "failed",
						Error:  err.Error(),
					}
					formatterOut.PrintError(fmt.Sprintf(i18n.T("failed_to_process"), job.Name, err.Error()))
					// Write error file immediately
					if isBatch && outputDir != "" {
						outputFileName := sanitizeFileName(job.Name) + ".error.json"
						outputPath := filepath.Join(outputDir, outputFileName)
						errData := map[string]interface{}{
							"input":  job.Name,
							"status": "failed",
							"error":  err.Error(),
						}
						if wErr := fileHandler.WriteJSONOutput(errData, outputPath); wErr != nil {
							formatterOut.PrintError(fmt.Sprintf("Failed to write error for %s: %s", job.Name, wErr.Error()))
						}
					}
				} else {
					results[job.Index] = BatchResult{
						Input:  job.Name,
						Status: "success",
						Data:   data,
					}
					// Write result file immediately
					if isBatch && outputDir != "" {
						outputFileName := sanitizeFileName(job.Name) + ".json"
						outputPath := filepath.Join(outputDir, outputFileName)
						if wErr := fileHandler.WriteJSONOutput(data, outputPath); wErr != nil {
							formatterOut.PrintError(fmt.Sprintf("Failed to write result for %s: %s", job.Name, wErr.Error()))
						}
					}
				}
			}
		}()
	}

	for _, job := range jobs {
		jobCh <- job
	}
	close(jobCh)
	wg.Wait()

	// Output summary
	outputBatchSummary(results, exportPath, outputDir, total)

	return results
}

// batchSubmit runs async submit jobs concurrently and outputs all task IDs as a single JSON array
func batchSubmit(jobs []BatchJob, concurrency int, submitFn SubmitFunc, exportPath string) []SubmitResult {
	total := len(jobs)
	results := make([]SubmitResult, total)

	formatterOut.PrintInfo(fmt.Sprintf(i18n.T("submitting_tasks"), total))

	if concurrency < 1 {
		concurrency = 1
	}

	jobCh := make(chan BatchJob, total)
	var wg sync.WaitGroup

	for w := 0; w < concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobCh {
				formatterOut.PrintProgress(job.Index+1, total, job.Name)

				taskID, err := submitFn(job)
				if err != nil {
					results[job.Index] = SubmitResult{
						Path:  job.Name,
						Error: err.Error(),
					}
					formatterOut.PrintError(fmt.Sprintf(i18n.T("failed_to_process"), job.Name, err.Error()))
				} else {
					results[job.Index] = SubmitResult{
						Path:   job.Name,
						TaskID: taskID,
					}
					formatterOut.PrintInfo(fmt.Sprintf("%s (task_id: %s)", job.Name, taskID))
				}
			}
		}()
	}

	for _, job := range jobs {
		jobCh <- job
	}
	close(jobCh)
	wg.Wait()

	// Output all results as a single JSON array
	fileHandler := file.New()
	if exportPath != "" {
		if err := fileHandler.WriteJSONOutput(results, exportPath); err != nil {
			formatterOut.PrintError(fmt.Sprintf("Failed to write task IDs: %s", err.Error()))
			formatterOut.PrintJSON(results)
		} else {
			formatterOut.PrintSuccess(fmt.Sprintf(i18n.T("results_exported_to"), absPath(exportPath)))
		}
	} else {
		formatterOut.PrintJSON(results)
	}

	return results
}

// outputBatchSummary handles writing the summary after individual result files have been written
func outputBatchSummary(results []BatchResult, exportPath, outputDir string, total int) {
	fileHandler := file.New()
	isBatch := total > 1

	if !isBatch {
		// Single file: preserve existing behavior
		if len(results) > 0 && results[0].Status == "success" {
			if exportPath != "" {
				fileHandler.WriteJSONOutput(results[0].Data, exportPath)
				formatterOut.PrintSuccess(fmt.Sprintf(i18n.T("results_exported_to"), exportPath))
			} else {
				formatterOut.PrintJSON(results[0].Data)
			}
		} else if len(results) > 0 {
			formatterOut.PrintJSON(map[string]interface{}{
				"input":  results[0].Input,
				"status": "failed",
				"error":  results[0].Error,
			})
		}
		return
	}

	// Batch mode: individual files already written, just build and write summary
	if outputDir == "" {
		// outputDir creation failed earlier, fallback to stdout
		formatterOut.PrintJSON(collectResultData(results))
		return
	}

	summary := BatchSummary{
		Total:     len(results),
		OutputDir: absPath(outputDir),
		Files:     make([]BatchSummaryFileItem, 0, len(results)),
	}

	for _, r := range results {
		var outputFileName string
		if r.Status == "success" {
			outputFileName = sanitizeFileName(r.Input) + ".json"
			summary.Success++
		} else {
			outputFileName = sanitizeFileName(r.Input) + ".error.json"
			summary.Failed++
		}

		item := BatchSummaryFileItem{
			Input:  r.Input,
			Output: outputFileName,
			Status: r.Status,
		}
		if r.Status == "failed" {
			item.Error = r.Error
		}
		summary.Files = append(summary.Files, item)
	}

	// Write _summary.json
	summaryPath := filepath.Join(outputDir, "_summary.json")
	if err := fileHandler.WriteJSONOutput(summary, summaryPath); err != nil {
		formatterOut.PrintError(fmt.Sprintf("Failed to write summary: %s", err.Error()))
	}

	// Print summary to stdout for Agent consumption
	formatterOut.PrintJSON(summary)
}

// exitWithBatchResults determines the exit code based on batch results
func exitWithBatchResults(results []BatchResult) {
	if len(results) == 0 {
		return
	}

	successCount := 0
	failCount := 0
	for _, r := range results {
		if r.Status == "success" {
			successCount++
		} else {
			failCount++
		}
	}

	if failCount == 0 {
		// All success, exit 0 (default)
		return
	}
	if successCount == 0 {
		// All failed
		flushTelemetry()
		os.Exit(errors.ExitGeneralError)
	}
	// Partial failure
	flushTelemetry()
	os.Exit(errors.ExitPartialFailure)
}

// retryWithBackoff retries a function with exponential backoff
func retryWithBackoff(fn func() (map[string]interface{}, error), maxRetry int) (map[string]interface{}, error) {
	var lastErr error
	for attempt := 0; attempt <= maxRetry; attempt++ {
		result, err := fn()
		if err == nil {
			return result, nil
		}
		lastErr = err

		// Only retry on retryable errors
		if attempt < maxRetry {
			cliErr := errors.ClassifyException(err, "")
			if !cliErr.Retryable {
				return nil, err
			}
			backoff := time.Duration(math.Pow(2, float64(attempt))) * time.Second
			formatterOut.PrintWarning(fmt.Sprintf("Retrying (%d/%d) after %v: %s", attempt+1, maxRetry, backoff, err.Error()))
			time.Sleep(backoff)
		}
	}
	return nil, lastErr
}

// collectResultData collects successful result data for fallback stdout output
func collectResultData(results []BatchResult) interface{} {
	var data []map[string]interface{}
	for _, r := range results {
		entry := map[string]interface{}{
			"input":  r.Input,
			"status": r.Status,
		}
		if r.Status == "success" {
			entry["data"] = r.Data
		} else {
			entry["error"] = r.Error
		}
		data = append(data, entry)
	}
	return data
}

// sanitizeFileName makes a string safe for use as a filename
func sanitizeFileName(name string) string {
	// Replace path separators and other problematic characters
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	return replacer.Replace(name)
}

// absPath returns the absolute path, falling back to the original on error
func absPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abs
}

// loadTasksFromFile reads tasks from a JSON file produced by --no-wait
func loadTasksFromFile(filePath string) ([]SubmitResult, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read task file: %w", err)
	}

	var items []SubmitResult
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse task file: %w", err)
	}

	var valid []SubmitResult
	for _, item := range items {
		if item.TaskID != "" {
			valid = append(valid, item)
		}
	}

	if len(valid) == 0 {
		return nil, fmt.Errorf("no task IDs found in file: %s", filePath)
	}

	return valid, nil
}

// ensureMaxConcurrency loads concurrency from config, refreshes from API if needed, and returns it
func ensureMaxConcurrency(client *api.Client) int {
	cfg, err := config.Load()
	if err != nil {
		return 1
	}
	if config.ShouldRefreshConcurrency(cfg) {
		concurrency, err := client.GetUserConcurrencyLimit()
		if err == nil {
			config.SetMaxConcurrency(concurrency)
			cfg.MaxConcurrency = concurrency
		}
	}
	if cfg.MaxConcurrency == 0 {
		return 1
	}
	return cfg.MaxConcurrency
}

// initClientWithConfig loads config and creates an API client, exiting on error
func initClientWithConfig(mode string) (*api.Client, *config.Config) {
	cfg, err := config.Load()
	if err != nil {
		cliErr := errors.ClassifyException(err, mode)
		formatterOut.ExitWithError(cliErr)
	}

	if !config.IsConfigured(cfg) {
		cliErr := errors.NewCLIError(
			i18n.T("error_not_configured"),
			errors.ErrorTypeSystem,
			errors.ExitGeneralError,
			false,
			"Run 'adp config set' to configure your API key and base URL.",
			nil,
		)
		formatterOut.ExitWithError(cliErr)
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		cliErr := errors.ClassifyException(err, mode)
		formatterOut.ExitWithError(cliErr)
	}

	return client, cfg
}

// processLocalFiles processes local files in batch mode
func processLocalFiles(pathStr, appID string, asyncMode, noWait bool, exportPath string, timeout, maxRetry int, mode string) {
	client, _ := initClientWithConfig(mode + " local")
	maxConcurrency := ensureMaxConcurrency(client)

	fileHandler := file.New()
	files, err := fileHandler.GetFilesFromPath(pathStr)
	if err != nil {
		cliErr := errors.ClassifyException(err, mode+" local")
		formatterOut.ExitWithError(cliErr)
	}

	validFiles, _ := fileHandler.ValidateFiles(files)
	if len(validFiles) == 0 {
		cliErr := errors.NewCLIError(
			i18n.T("no_valid_files"),
			errors.ErrorTypeParam,
			errors.ExitParameterError,
			false,
			"Check the file path and try again.",
			nil,
		)
		formatterOut.ExitWithError(cliErr)
	}

	// Build batch jobs
	jobs := make([]BatchJob, len(validFiles))
	for i, filePath := range validFiles {
		jobs[i] = BatchJob{Index: i, Name: filepath.Base(filePath)}
	}

	effectiveConcurrency := min(maxConcurrency, len(jobs))
	if len(jobs) > 1 {
		log.Debug().Int("max_concurrency", maxConcurrency).Int("file_count", len(jobs)).Int("effective_concurrency", effectiveConcurrency).Msg("Batch concurrency")
	}

	// noWait mode: submit only, output task IDs
	if asyncMode && noWait {
		submitFn := func(job BatchJob) (string, error) {
			filePath := validFiles[job.Index]
			fileName := job.Name
			if mode == "parse" {
				return client.ParseAsync("", appID, filePath, "", fileName)
			}
			return client.ExtractAsync("", appID, filePath, "", fileName, nil)
		}
		batchSubmit(jobs, effectiveConcurrency, submitFn, exportPath)
		return
	}

	processFn := func(job BatchJob) (map[string]interface{}, error) {
		filePath := validFiles[job.Index]
		fileName := job.Name
		if mode == "parse" {
			if asyncMode {
				taskID, err := client.ParseAsync("", appID, filePath, "", fileName)
				if err != nil {
					return nil, err
				}
				formatterOut.PrintInfo(fmt.Sprintf("%s (task_id: %s)", fileName, taskID))
				return client.WaitForTask(taskID, client.QueryParseTask, timeout, 2)
			}
			return client.ParseSync("", appID, filePath, "", fileName)
		}
		// extract or human-review
		disableCollaboration := (mode == "human-review")
		if asyncMode {
			taskID, err := client.ExtractAsync("", appID, filePath, "", fileName, nil)
			if err != nil {
				return nil, err
			}
			formatterOut.PrintInfo(fmt.Sprintf("%s (task_id: %s)", fileName, taskID))
			return client.WaitForTask(taskID, client.QueryExtractTask, timeout, 2)
		}
		result, err := client.ExtractSync("", appID, filePath, "", fileName, nil, disableCollaboration)
		if err != nil {
			return nil, err
		}
		checkAPIResponse(result, mode)
		return result, nil
	}

	results := batchProcess(jobs, effectiveConcurrency, maxRetry, processFn, exportPath)
	exitWithBatchResults(results)
}

// processURLs processes URLs in batch mode
func processURLs(urls []string, appID string, asyncMode, noWait bool, exportPath string, timeout, maxRetry int, mode string) {
	client, _ := initClientWithConfig(mode + " url")
	maxConcurrency := ensureMaxConcurrency(client)

	var validURLs []string
	for _, url := range urls {
		if !file.IsValidURL(url) {
			formatterOut.PrintWarning(fmt.Sprintf(i18n.T("invalid_url_format"), url))
		} else {
			validURLs = append(validURLs, url)
		}
	}

	if len(validURLs) == 0 {
		cliErr := errors.NewCLIError(
			i18n.T("no_valid_files"),
			errors.ErrorTypeParam,
			errors.ExitParameterError,
			false,
			"Check the URL format and try again.",
			nil,
		)
		formatterOut.ExitWithError(cliErr)
	}

	// Build batch jobs — use url_001 style naming for result files
	jobs := make([]BatchJob, len(validURLs))
	for i, url := range validURLs {
		jobs[i] = BatchJob{Index: i, Name: fmt.Sprintf("url_%03d", i+1)}
		_ = url // url is captured via validURLs[job.Index] in processFn/submitFn
	}

	effectiveConcurrency := min(maxConcurrency, len(jobs))
	if len(jobs) > 1 {
		log.Debug().Int("max_concurrency", maxConcurrency).Int("file_count", len(jobs)).Int("effective_concurrency", effectiveConcurrency).Msg("Batch concurrency")
	}

	// noWait mode: submit only, output task IDs
	if asyncMode && noWait {
		submitFn := func(job BatchJob) (string, error) {
			url := validURLs[job.Index]
			if mode == "parse" {
				return client.ParseAsync(url, appID, "", "", "")
			}
			return client.ExtractAsync(url, appID, "", "", "", nil)
		}
		batchSubmit(jobs, effectiveConcurrency, submitFn, exportPath)
		return
	}

	processFn := func(job BatchJob) (map[string]interface{}, error) {
		url := validURLs[job.Index]
		if mode == "parse" {
			if asyncMode {
				taskID, err := client.ParseAsync(url, appID, "", "", "")
				if err != nil {
					return nil, err
				}
				formatterOut.PrintInfo(fmt.Sprintf("%s (task_id: %s)", job.Name, taskID))
				return client.WaitForTask(taskID, client.QueryParseTask, timeout, 2)
			}
			return client.ParseSync(url, appID, "", "", "")
		}
		disableCollaboration := (mode == "human-review")
		if asyncMode {
			taskID, err := client.ExtractAsync(url, appID, "", "", "", nil)
			if err != nil {
				return nil, err
			}
			formatterOut.PrintInfo(fmt.Sprintf("%s (task_id: %s)", job.Name, taskID))
			return client.WaitForTask(taskID, client.QueryExtractTask, timeout, 2)
		}
		result, err := client.ExtractSync(url, appID, "", "", "", nil, disableCollaboration)
		if err != nil {
			return nil, err
		}
		checkAPIResponse(result, mode)
		return result, nil
	}

	results := batchProcess(jobs, effectiveConcurrency, maxRetry, processFn, exportPath)
	exitWithBatchResults(results)
}

// processBase64 processes base64 strings in batch mode
func processBase64(b64Strings []string, appID, fileName string, asyncMode, noWait bool, exportPath string, timeout, maxRetry int, mode string) {
	client, _ := initClientWithConfig(mode + " base64")
	maxConcurrency := ensureMaxConcurrency(client)

	// Build batch jobs
	jobs := make([]BatchJob, len(b64Strings))
	for i := range b64Strings {
		name := fileName
		if len(b64Strings) > 1 {
			name = fmt.Sprintf("%s_%d", fileName, i+1)
		}
		jobs[i] = BatchJob{Index: i, Name: name}
	}

	effectiveConcurrency := min(maxConcurrency, len(jobs))
	if len(jobs) > 1 {
		log.Debug().Int("max_concurrency", maxConcurrency).Int("file_count", len(jobs)).Int("effective_concurrency", effectiveConcurrency).Msg("Batch concurrency")
	}

	// noWait mode: submit only, output task IDs
	if asyncMode && noWait {
		submitFn := func(job BatchJob) (string, error) {
			b64 := b64Strings[job.Index]
			currentFileName := job.Name
			if mode == "parse" {
				return client.ParseAsync("", appID, "", b64, currentFileName)
			}
			return client.ExtractAsync("", appID, "", b64, currentFileName, nil)
		}
		batchSubmit(jobs, effectiveConcurrency, submitFn, exportPath)
		return
	}

	processFn := func(job BatchJob) (map[string]interface{}, error) {
		b64 := b64Strings[job.Index]
		currentFileName := job.Name
		if mode == "parse" {
			if asyncMode {
				taskID, err := client.ParseAsync("", appID, "", b64, currentFileName)
				if err != nil {
					return nil, err
				}
				formatterOut.PrintInfo(fmt.Sprintf("%s (task_id: %s)", currentFileName, taskID))
				return client.WaitForTask(taskID, client.QueryParseTask, timeout, 2)
			}
			return client.ParseSync("", appID, "", b64, currentFileName)
		}
		disableCollaboration := (mode == "human-review")
		if asyncMode {
			taskID, err := client.ExtractAsync("", appID, "", b64, currentFileName, nil)
			if err != nil {
				return nil, err
			}
			formatterOut.PrintInfo(fmt.Sprintf("%s (task_id: %s)", currentFileName, taskID))
			return client.WaitForTask(taskID, client.QueryExtractTask, timeout, 2)
		}
		result, err := client.ExtractSync("", appID, "", b64, currentFileName, nil, disableCollaboration)
		if err != nil {
			return nil, err
		}
		checkAPIResponse(result, mode)
		return result, nil
	}

	results := batchProcess(jobs, effectiveConcurrency, maxRetry, processFn, exportPath)
	exitWithBatchResults(results)
}

// queryTasks queries multiple tasks in batch mode
func queryTasks(taskIDs []string, watch bool, taskFile, exportPath string, timeout int, mode string) {
	// Load tasks from file if --file is specified
	var taskNames map[string]string // taskID -> path (display name)
	if taskFile != "" {
		tasks, err := loadTasksFromFile(taskFile)
		if err != nil {
			cliErr := errors.ClassifyException(err, mode+" query")
			formatterOut.ExitWithError(cliErr)
		}
		taskNames = make(map[string]string, len(tasks))
		taskIDs = make([]string, 0, len(tasks))
		for _, t := range tasks {
			taskIDs = append(taskIDs, t.TaskID)
			taskNames[t.TaskID] = t.Path
		}
	}

	client, _ := initClientWithConfig(mode + " query")
	maxConcurrency := ensureMaxConcurrency(client)

	if len(taskIDs) == 1 {
		// Single task query: preserve existing behavior
		querySingleTask(client, taskIDs[0], watch, exportPath, timeout, mode)
		return
	}

	// Multi task query: use batch processing
	effectiveConcurrency := min(maxConcurrency, len(taskIDs))
	if len(taskIDs) > 1 {
		log.Debug().Int("max_concurrency", maxConcurrency).Int("task_count", len(taskIDs)).Int("effective_concurrency", effectiveConcurrency).Msg("Batch concurrency")
	}

	var queryFunc func(string) (map[string]interface{}, error)
	if mode == "parse" {
		queryFunc = client.QueryParseTask
	} else {
		queryFunc = client.QueryExtractTask
	}

	jobs := make([]BatchJob, len(taskIDs))
	for i, id := range taskIDs {
		name := id
		if taskNames != nil {
			if p, ok := taskNames[id]; ok && p != "" {
				name = p
			}
		}
		jobs[i] = BatchJob{Index: i, Name: name}
	}

	processFn := func(job BatchJob) (map[string]interface{}, error) {
		taskID := taskIDs[job.Index]
		if watch {
			return client.WaitForTask(taskID, queryFunc, timeout, 2)
		}
		return queryFunc(taskID)
	}

	results := batchProcess(jobs, effectiveConcurrency, 0, processFn, exportPath)
	exitWithBatchResults(results)
}

// querySingleTask queries a single task
func querySingleTask(client *api.Client, taskID string, watch bool, exportPath string, timeout int, mode string) {
	var queryFunc func(string) (map[string]interface{}, error)
	if mode == "parse" {
		queryFunc = client.QueryParseTask
	} else {
		queryFunc = client.QueryExtractTask
	}

	var result map[string]interface{}
	var err error

	if watch {
		result, err = client.WaitForTask(taskID, queryFunc, timeout, 2)
		if err != nil {
			cliErr := errors.ClassifyException(err, mode+" query")
			formatterOut.ExitWithError(cliErr)
		}
		formatterOut.PrintSuccess(i18n.T("task_completed"))
	} else {
		result, err = queryFunc(taskID)
		if err != nil {
			cliErr := errors.ClassifyException(err, mode+" query")
			formatterOut.ExitWithError(cliErr)
		}
	}

	if exportPath != "" {
		fileHandler := file.New()
		fileHandler.WriteJSONOutput(result, exportPath)
		formatterOut.PrintSuccess(fmt.Sprintf(i18n.T("results_exported_to"), exportPath))
	} else {
		formatterOut.PrintJSON(result)
	}
}

// resolveURLInput checks if the input is a local file path containing URLs or a direct URL.
func resolveURLInput(input string) []string {
	info, err := os.Stat(input)
	if err == nil && !info.IsDir() {
		fh := file.New()
		urls, err := fh.ReadURLListFile(input)
		if err != nil {
			cliErr := errors.NewCLIError(
				fmt.Sprintf("Failed to read URL list file: %v", err),
				errors.ErrorTypeParam,
				errors.ExitParameterError,
				false,
				"Check the file path and format.",
				nil,
			)
			formatterOut.ExitWithError(cliErr)
		}
		if len(urls) == 0 {
			cliErr := errors.NewCLIError(
				fmt.Sprintf("No valid URLs found in file: %s", input),
				errors.ErrorTypeParam,
				errors.ExitParameterError,
				false,
				"Each line should contain a URL starting with http:// or https://",
				nil,
			)
			formatterOut.ExitWithError(cliErr)
		}
		return urls
	}
	return []string{input}
}
