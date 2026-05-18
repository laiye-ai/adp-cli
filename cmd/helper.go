package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/laiye-ai/adp-cli/internal/api"
	"github.com/laiye-ai/adp-cli/internal/config"
	"github.com/laiye-ai/adp-cli/internal/errors"
	"github.com/laiye-ai/adp-cli/internal/file"
	"github.com/laiye-ai/adp-cli/internal/i18n"
	"github.com/laiye-ai/adp-cli/internal/telemetry"
)

const (
	freeLimit = 1
	paidLimit = 2
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

// validateConcurrency checks and adjusts concurrency based on user payment status.
func validateConcurrency(client *api.Client, concurrency int) int {
	if concurrency <= freeLimit {
		return freeLimit
	}

	isPaid, err := client.IsPaidUser()
	if err != nil {
		formatterOut.PrintWarning(i18n.T("error_invalid_concurrency"))
		return freeLimit
	}

	if isPaid {
		if concurrency > paidLimit {
			formatterOut.PrintWarning(i18n.T("error_invalid_concurrency"))
			return paidLimit
		}
		return concurrency
	}

	if concurrency > freeLimit {
		formatterOut.PrintWarning(i18n.T("error_not_paid_user"))
	}
	return freeLimit
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
			errors.ExitGeneralError, false,
			"Run 'adp config set' to configure your API key and base URL.", nil,
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

// processLocalFiles processes local files for parse, extract, or human-review
func processLocalFiles(pathStr, appID string, asyncMode, noWait bool, exportPath string, timeout, concurrency, maxRetry int, mode string) {
	client, _ := initClientWithConfig(mode + " local")
	concurrency = validateConcurrency(client, concurrency)

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
			errors.ExitParameterError, false,
			"Check the file path and try again.", nil,
		)
		formatterOut.ExitWithError(cliErr)
	}

	jobs := make([]BatchJob, len(validFiles))
	for i, filePath := range validFiles {
		jobs[i] = BatchJob{Index: i, Name: filepath.Base(filePath)}
	}

	if asyncMode && noWait {
		submitFn := func(job BatchJob) (string, error) {
			filePath := validFiles[job.Index]
			fileName := job.Name
			if mode == "parse" {
				return client.ParseAsync("", appID, filePath, "", fileName)
			}
			return client.ExtractAsync("", appID, filePath, "", fileName, nil)
		}
		batchSubmit(jobs, concurrency, submitFn, exportPath)
		return
	}

	processFn := func(job BatchJob) (map[string]interface{}, error) {
		filePath := validFiles[job.Index]
		fileName := job.Name
		var result map[string]interface{}
		if mode == "parse" {
			if asyncMode {
				taskID, err := client.ParseAsync("", appID, filePath, "", fileName)
				if err != nil {
					return nil, err
				}
				formatterOut.PrintInfo(fmt.Sprintf("%s (task_id: %s)", fileName, taskID))
				result, err = client.WaitForTask(taskID, client.QueryParseTask, timeout, 2)
			} else {
				result, err = client.ParseSync("", appID, filePath, "", fileName)
			}
		} else {
			disableCollaboration := (mode == "human-review")
			if asyncMode {
				taskID, err := client.ExtractAsync("", appID, filePath, "", fileName, nil)
				if err != nil {
					return nil, err
				}
				formatterOut.PrintInfo(fmt.Sprintf("%s (task_id: %s)", fileName, taskID))
				result, err = client.WaitForTask(taskID, client.QueryExtractTask, timeout, 2)
			} else {
				result, err = client.ExtractSync("", appID, filePath, "", fileName, nil, disableCollaboration)
			}
		}
		if err != nil {
			return nil, err
		}
		checkAPIResponse(result, mode)
		return result, nil
	}

	results := batchProcess(jobs, concurrency, maxRetry, processFn, exportPath)
	exitWithBatchResults(results)
}

// processURLs processes URLs for parse, extract, or human-review
func processURLs(urls []string, appID string, asyncMode, noWait bool, exportPath string, timeout, concurrency, maxRetry int, mode string) {
	client, _ := initClientWithConfig(mode + " url")
	concurrency = validateConcurrency(client, concurrency)

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
			errors.ExitParameterError, false,
			"Check the URL format and try again.", nil,
		)
		formatterOut.ExitWithError(cliErr)
	}

	jobs := make([]BatchJob, len(validURLs))
	for i, url := range validURLs {
		jobs[i] = BatchJob{Index: i, Name: fmt.Sprintf("url_%03d", i+1)}
		_ = url
	}

	if asyncMode && noWait {
		submitFn := func(job BatchJob) (string, error) {
			url := validURLs[job.Index]
			if mode == "parse" {
				return client.ParseAsync(url, appID, "", "", "")
			}
			return client.ExtractAsync(url, appID, "", "", "", nil)
		}
		batchSubmit(jobs, concurrency, submitFn, exportPath)
		return
	}

	processFn := func(job BatchJob) (map[string]interface{}, error) {
		url := validURLs[job.Index]
		var result map[string]interface{}
		var err error
		if mode == "parse" {
			if asyncMode {
				taskID, err := client.ParseAsync(url, appID, "", "", "")
				if err != nil {
					return nil, err
				}
				formatterOut.PrintInfo(fmt.Sprintf("%s (task_id: %s)", job.Name, taskID))
				result, err = client.WaitForTask(taskID, client.QueryParseTask, timeout, 2)
			} else {
				result, err = client.ParseSync(url, appID, "", "", "")
			}
		} else {
			disableCollaboration := (mode == "human-review")
			if asyncMode {
				taskID, err := client.ExtractAsync(url, appID, "", "", "", nil)
				if err != nil {
					return nil, err
				}
				formatterOut.PrintInfo(fmt.Sprintf("%s (task_id: %s)", job.Name, taskID))
				result, err = client.WaitForTask(taskID, client.QueryExtractTask, timeout, 2)
			} else {
				result, err = client.ExtractSync(url, appID, "", "", "", nil, disableCollaboration)
			}
		}
		if err != nil {
			return nil, err
		}
		checkAPIResponse(result, mode)
		return result, nil
	}

	results := batchProcess(jobs, concurrency, maxRetry, processFn, exportPath)
	exitWithBatchResults(results)
}

// processBase64 processes base64 strings for parse, extract, or human-review
func processBase64(b64Strings []string, appID, fileName string, asyncMode, noWait bool, exportPath string, timeout, concurrency, maxRetry int, mode string) {
	client, _ := initClientWithConfig(mode + " base64")
	concurrency = validateConcurrency(client, concurrency)

	jobs := make([]BatchJob, len(b64Strings))
	for i := range b64Strings {
		name := fileName
		if len(b64Strings) > 1 {
			name = fmt.Sprintf("%s_%d", fileName, i+1)
		}
		jobs[i] = BatchJob{Index: i, Name: name}
	}

	if asyncMode && noWait {
		submitFn := func(job BatchJob) (string, error) {
			b64 := b64Strings[job.Index]
			currentFileName := job.Name
			if mode == "parse" {
				return client.ParseAsync("", appID, "", b64, currentFileName)
			}
			return client.ExtractAsync("", appID, "", b64, currentFileName, nil)
		}
		batchSubmit(jobs, concurrency, submitFn, exportPath)
		return
	}

	processFn := func(job BatchJob) (map[string]interface{}, error) {
		b64 := b64Strings[job.Index]
		currentFileName := job.Name
		var result map[string]interface{}
		var err error
		if mode == "parse" {
			if asyncMode {
				taskID, err := client.ParseAsync("", appID, "", b64, currentFileName)
				if err != nil {
					return nil, err
				}
				formatterOut.PrintInfo(fmt.Sprintf("%s (task_id: %s)", currentFileName, taskID))
				result, err = client.WaitForTask(taskID, client.QueryParseTask, timeout, 2)
			} else {
				result, err = client.ParseSync("", appID, "", b64, currentFileName)
			}
		} else {
			disableCollaboration := (mode == "human-review")
			if asyncMode {
				taskID, err := client.ExtractAsync("", appID, "", b64, currentFileName, nil)
				if err != nil {
					return nil, err
				}
				formatterOut.PrintInfo(fmt.Sprintf("%s (task_id: %s)", currentFileName, taskID))
				result, err = client.WaitForTask(taskID, client.QueryExtractTask, timeout, 2)
			} else {
				result, err = client.ExtractSync("", appID, "", b64, currentFileName, nil, disableCollaboration)
			}
		}
		if err != nil {
			return nil, err
		}
		checkAPIResponse(result, mode)
		return result, nil
	}

	results := batchProcess(jobs, concurrency, maxRetry, processFn, exportPath)
	exitWithBatchResults(results)
}

// queryTasks queries task status for parse, extract, or human-review
func queryTasks(taskIDs []string, watch bool, taskFile, exportPath string, timeout, concurrency int, mode string) {
	var taskNames map[string]string
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

	if len(taskIDs) == 1 {
		querySingleTask(client, taskIDs[0], watch, exportPath, timeout, mode)
		return
	}

	concurrency = validateConcurrency(client, concurrency)

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

	results := batchProcess(jobs, concurrency, 0, processFn, exportPath)
	exitWithBatchResults(results)
}

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
				errors.ExitParameterError, false,
				"Check the file path and format.", nil,
			)
			formatterOut.ExitWithError(cliErr)
		}
		if len(urls) == 0 {
			cliErr := errors.NewCLIError(
				fmt.Sprintf("No valid URLs found in file: %s", input),
				errors.ErrorTypeParam,
				errors.ExitParameterError, false,
				"Each line should contain a URL starting with http:// or https://", nil,
			)
			formatterOut.ExitWithError(cliErr)
		}
		return urls
	}
	return []string{input}
}