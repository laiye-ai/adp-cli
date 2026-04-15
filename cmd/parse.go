package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/cobra"
	"github.com/laiye-ai/adp-cli/internal/api"
	"github.com/laiye-ai/adp-cli/internal/config"
	"github.com/laiye-ai/adp-cli/internal/errors"
	"github.com/laiye-ai/adp-cli/internal/file"
	"github.com/laiye-ai/adp-cli/internal/i18n"
)

// parseCmd represents the parse command group
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: i18n.T("parse_description"),
	Long:  i18n.T("parse_description"),
}

// parseLocalCmd represents the parse local command
var parseLocalCmd = &cobra.Command{
	Use:   "local",
	Short: i18n.T("parse_local_title"),
	Long:  i18n.T("parse_local_title"),
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		appID, _ := cmd.Flags().GetString("app-id")
		asyncMode, _ := cmd.Flags().GetBool("async")
		export, _ := cmd.Flags().GetString("export")
		timeout, _ := cmd.Flags().GetInt("timeout")
		concurrency, _ := cmd.Flags().GetInt("concurrency")

		processLocalFiles(path, appID, asyncMode, export, timeout, concurrency, "parse")
	},
}

// parseURLCmd represents the parse URL command
var parseURLCmd = &cobra.Command{
	Use:   "url",
	Short: i18n.T("parse_url_title"),
	Long:  i18n.T("parse_url_title"),
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		urlArg := args[0]
		appID, _ := cmd.Flags().GetString("app-id")
		asyncMode, _ := cmd.Flags().GetBool("async")
		export, _ := cmd.Flags().GetString("export")
		timeout, _ := cmd.Flags().GetInt("timeout")
		concurrency, _ := cmd.Flags().GetInt("concurrency")

		urls := resolveURLInput(urlArg)
		processURLs(urls, appID, asyncMode, export, timeout, concurrency, "parse")
	},
}

// parseBase64Cmd represents the parse base64 command
var parseBase64Cmd = &cobra.Command{
	Use:   "base64",
	Short: i18n.T("parse_base64_title"),
	Long:  i18n.T("parse_base64_title"),
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appID, _ := cmd.Flags().GetString("app-id")
		asyncMode, _ := cmd.Flags().GetBool("async")
		export, _ := cmd.Flags().GetString("export")
		timeout, _ := cmd.Flags().GetInt("timeout")
		fileName, _ := cmd.Flags().GetString("file-name")
		concurrency, _ := cmd.Flags().GetInt("concurrency")

		processBase64(args, appID, fileName, asyncMode, export, timeout, concurrency, "parse")
	},
}

// parseQueryCmd represents the parse query command
var parseQueryCmd = &cobra.Command{
	Use:   "query",
	Short: i18n.T("parse_query_title"),
	Long:  i18n.T("parse_query_title"),
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskID := args[0]
		watch, _ := cmd.Flags().GetBool("watch")
		export, _ := cmd.Flags().GetString("export")
		timeout, _ := cmd.Flags().GetInt("timeout")

		queryTask(taskID, watch, export, timeout, "parse")
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)

	parseCmd.AddCommand(parseLocalCmd)
	parseCmd.AddCommand(parseURLCmd)
	parseCmd.AddCommand(parseBase64Cmd)
	parseCmd.AddCommand(parseQueryCmd)

	// parse local flags
	parseLocalCmd.Flags().String("app-id", "", i18n.T("option_app_id_parse"))
	parseLocalCmd.Flags().Bool("async", false, i18n.T("option_async"))
	parseLocalCmd.Flags().String("export", "", i18n.T("option_export"))
	parseLocalCmd.Flags().Int("timeout", 900, i18n.T("option_timeout"))
	parseLocalCmd.Flags().Int("concurrency", 1, i18n.T("option_concurrency"))
	parseLocalCmd.MarkFlagRequired("app-id")

	// parse url flags
	parseURLCmd.Flags().String("app-id", "", i18n.T("option_app_id_parse"))
	parseURLCmd.Flags().Bool("async", false, i18n.T("option_async"))
	parseURLCmd.Flags().String("export", "", i18n.T("option_export"))
	parseURLCmd.Flags().Int("timeout", 900, i18n.T("option_timeout"))
	parseURLCmd.Flags().Int("concurrency", 1, i18n.T("option_concurrency"))
	parseURLCmd.MarkFlagRequired("app-id")

	// parse base64 flags
	parseBase64Cmd.Flags().String("app-id", "", i18n.T("option_app_id_parse"))
	parseBase64Cmd.Flags().Bool("async", false, i18n.T("option_async"))
	parseBase64Cmd.Flags().String("export", "", i18n.T("option_export"))
	parseBase64Cmd.Flags().Int("timeout", 900, i18n.T("option_timeout"))
	parseBase64Cmd.Flags().String("file-name", "document", i18n.T("option_file_name"))
	parseBase64Cmd.Flags().Int("concurrency", 1, i18n.T("option_concurrency"))
	parseBase64Cmd.MarkFlagRequired("app-id")

	// parse query flags
	parseQueryCmd.Flags().Bool("watch", false, i18n.T("option_watch"))
	parseQueryCmd.Flags().String("export", "", i18n.T("option_export"))
	parseQueryCmd.Flags().Int("timeout", 900, i18n.T("option_watch_timeout"))
}

const (
	freeLimit = 1
	paidLimit = 2
)

// validateConcurrency checks and adjusts concurrency based on user payment status.
// Free users max=1, paid users max=2. Warns on invalid values and degrades gracefully.
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

	// Free user trying to use > 1
	if concurrency > freeLimit {
		formatterOut.PrintWarning(i18n.T("error_not_paid_user"))
	}
	return freeLimit
}

func processLocalFiles(pathStr, appID string, asyncMode bool, exportPath string, timeout, concurrency int, mode string) {
	cfg, err := config.Load()
	if err != nil {
		cliErr := errors.ClassifyException(err, mode+" local")
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
		cliErr := errors.ClassifyException(err, mode+" local")
		formatterOut.ExitWithError(cliErr)
	}

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
			errors.ExitParameterError,
			false,
			"Check the file path and try again.",
			nil,
		)
		formatterOut.ExitWithError(cliErr)
	}

	formatterOut.PrintInfo(fmt.Sprintf(i18n.T("processing_files"), len(validFiles)))

	if concurrency < 1 {
		concurrency = 1
	}

	type fileJob struct {
		index    int
		filePath string
		fileName string
	}
	jobs := make(chan fileJob, len(validFiles))
	results := make([]map[string]interface{}, len(validFiles))
	var wg sync.WaitGroup

	// Start N workers
	for w := 0; w < concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				var result map[string]interface{}
				var err error

				if mode == "parse" {
					if asyncMode {
						taskID, e := client.ParseAsync("", appID, job.filePath, "", job.fileName)
						if e != nil {
							err = e
						} else {
							formatterOut.PrintProgress(job.index+1, len(validFiles), fmt.Sprintf("%s TaskID:%s", job.fileName, taskID))
							result, err = client.WaitForTask(taskID, client.QueryParseTask, timeout, 2)
						}
					} else {
						formatterOut.PrintProgress(job.index+1, len(validFiles), job.fileName)
						result, err = client.ParseSync("", appID, job.filePath, "", job.fileName)
					}
				} else {
					if asyncMode {
						taskID, e := client.ExtractAsync("", appID, job.filePath, "", job.fileName, nil)
						if e != nil {
							err = e
						} else {
							formatterOut.PrintProgress(job.index+1, len(validFiles), fmt.Sprintf("%s TaskID:%s", job.fileName, taskID))
							result, err = client.WaitForTask(taskID, client.QueryExtractTask, timeout, 2)
						}
					} else {
						formatterOut.PrintProgress(job.index+1, len(validFiles), job.fileName)
						result, err = client.ExtractSync("", appID, job.filePath, "", job.fileName, nil)
					}
				}

				if err != nil {
					formatterOut.PrintError(fmt.Sprintf(i18n.T("failed_to_process"), job.fileName, err.Error()))
				} else {
					results[job.index] = result
				}
			}
		}()
	}

	// Send jobs
	for i, filePath := range validFiles {
		jobs <- fileJob{index: i, filePath: filePath, fileName: filepath.Base(filePath)}
	}
	close(jobs)

	wg.Wait()

	// Collect non-nil results in order
	var orderedResults []map[string]interface{}
	for _, r := range results {
		if r != nil {
			orderedResults = append(orderedResults, r)
		}
	}

	if exportPath != "" {
		if len(orderedResults) == 1 {
			fileHandler.WriteJSONOutput(orderedResults[0], exportPath)
		} else {
			fileHandler.WriteJSONOutput(map[string]interface{}{"results": orderedResults}, exportPath)
		}
		formatterOut.PrintSuccess(fmt.Sprintf(i18n.T("results_exported_to"), exportPath))
	} else {
		formatterOut.PrintJSON(orderedResults)
	}
}

func processURLs(urls []string, appID string, asyncMode bool, exportPath string, timeout, concurrency int, mode string) {
	cfg, err := config.Load()
	if err != nil {
		cliErr := errors.ClassifyException(err, mode+" url")
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
		cliErr := errors.ClassifyException(err, mode+" url")
		formatterOut.ExitWithError(cliErr)
	}

	concurrency = validateConcurrency(client, concurrency)

	fileHandler := file.New()
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

	formatterOut.PrintInfo(fmt.Sprintf(i18n.T("processing_urls"), len(validURLs)))

	if concurrency < 1 {
		concurrency = 1
	}

	type urlJob struct {
		index int
		url   string
	}
	jobs := make(chan urlJob, len(validURLs))
	results := make([]map[string]interface{}, len(validURLs))
	var wg sync.WaitGroup

	for w := 0; w < concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				var result map[string]interface{}
				var err error

				if mode == "parse" {
					if asyncMode {
						taskID, e := client.ParseAsync(job.url, appID, "", "", "")
						if e != nil {
							err = e
						} else {
							formatterOut.PrintProgress(job.index+1, len(validURLs), fmt.Sprintf("%s TaskID:%s", job.url, taskID))
							result, err = client.WaitForTask(taskID, client.QueryParseTask, timeout, 2)
						}
					} else {
						formatterOut.PrintProgress(job.index+1, len(validURLs), job.url)
						result, err = client.ParseSync(job.url, appID, "", "", "")
					}
				} else {
					if asyncMode {
						taskID, e := client.ExtractAsync(job.url, appID, "", "", "", nil)
						if e != nil {
							err = e
						} else {
							formatterOut.PrintProgress(job.index+1, len(validURLs), fmt.Sprintf("%s TaskID:%s", job.url, taskID))
							result, err = client.WaitForTask(taskID, client.QueryExtractTask, timeout, 2)
						}
					} else {
						formatterOut.PrintProgress(job.index+1, len(validURLs), job.url)
						result, err = client.ExtractSync(job.url, appID, "", "", "", nil)
					}
				}

				if err != nil {
					formatterOut.PrintError(fmt.Sprintf(i18n.T("failed_to_process_url"), job.url, err.Error()))
				} else {
					results[job.index] = result
				}
			}
		}()
	}

	for i, url := range validURLs {
		jobs <- urlJob{index: i, url: url}
	}
	close(jobs)

	wg.Wait()

	var orderedResults []map[string]interface{}
	for _, r := range results {
		if r != nil {
			orderedResults = append(orderedResults, r)
		}
	}

	if exportPath != "" {
		if len(orderedResults) == 1 {
			fileHandler.WriteJSONOutput(orderedResults[0], exportPath)
		} else {
			fileHandler.WriteJSONOutput(map[string]interface{}{"results": orderedResults}, exportPath)
		}
		formatterOut.PrintSuccess(fmt.Sprintf(i18n.T("results_exported_to"), exportPath))
	} else {
		formatterOut.PrintJSON(orderedResults)
	}
}

func processBase64(b64Strings []string, appID, fileName string, asyncMode bool, exportPath string, timeout, concurrency int, mode string) {
	cfg, err := config.Load()
	if err != nil {
		cliErr := errors.ClassifyException(err, mode+" base64")
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
		cliErr := errors.ClassifyException(err, mode+" base64")
		formatterOut.ExitWithError(cliErr)
	}

	concurrency = validateConcurrency(client, concurrency)

	fileHandler := file.New()

	formatterOut.PrintInfo(fmt.Sprintf(i18n.T("processing_files"), len(b64Strings)))

	if concurrency < 1 {
		concurrency = 1
	}

	type b64Job struct {
		index         int
		b64           string
		currentFileName string
	}
	jobs := make(chan b64Job, len(b64Strings))
	results := make([]map[string]interface{}, len(b64Strings))
	var wg sync.WaitGroup

	for w := 0; w < concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				var result map[string]interface{}
				var err error

				if mode == "parse" {
					if asyncMode {
						taskID, e := client.ParseAsync("", appID, "", job.b64, job.currentFileName)
						if e != nil {
							err = e
						} else {
							formatterOut.PrintProgress(job.index+1, len(b64Strings), fmt.Sprintf("%s TaskID:%s", job.currentFileName, taskID))
							result, err = client.WaitForTask(taskID, client.QueryParseTask, timeout, 2)
						}
					} else {
						formatterOut.PrintProgress(job.index+1, len(b64Strings), job.currentFileName)
						result, err = client.ParseSync("", appID, "", job.b64, job.currentFileName)
					}
				} else {
					if asyncMode {
						taskID, e := client.ExtractAsync("", appID, "", job.b64, job.currentFileName, nil)
						if e != nil {
							err = e
						} else {
							formatterOut.PrintProgress(job.index+1, len(b64Strings), fmt.Sprintf("%s TaskID:%s", job.currentFileName, taskID))
							result, err = client.WaitForTask(taskID, client.QueryExtractTask, timeout, 2)
						}
					} else {
						formatterOut.PrintProgress(job.index+1, len(b64Strings), job.currentFileName)
						result, err = client.ExtractSync("", appID, "", job.b64, job.currentFileName, nil)
					}
				}

				if err != nil {
					formatterOut.PrintError(fmt.Sprintf(i18n.T("failed_to_process"), job.currentFileName, err.Error()))
				} else {
					results[job.index] = result
				}
			}
		}()
	}

	for i, b64 := range b64Strings {
		currentFileName := fileName
		if len(b64Strings) > 1 {
			currentFileName = fmt.Sprintf("%s_%d", fileName, i+1)
		}
		jobs <- b64Job{index: i, b64: b64, currentFileName: currentFileName}
	}
	close(jobs)

	wg.Wait()

	var orderedResults []map[string]interface{}
	for _, r := range results {
		if r != nil {
			orderedResults = append(orderedResults, r)
		}
	}

	if exportPath != "" {
		if len(orderedResults) == 1 {
			fileHandler.WriteJSONOutput(orderedResults[0], exportPath)
		} else {
			fileHandler.WriteJSONOutput(map[string]interface{}{"results": orderedResults}, exportPath)
		}
		formatterOut.PrintSuccess(fmt.Sprintf(i18n.T("results_exported_to"), exportPath))
	} else {
		formatterOut.PrintJSON(orderedResults)
	}
}

func queryTask(taskID string, watch bool, exportPath string, timeout int, mode string) {
	cfg, err := config.Load()
	if err != nil {
		cliErr := errors.ClassifyException(err, mode+" query")
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
		cliErr := errors.ClassifyException(err, mode+" query")
		formatterOut.ExitWithError(cliErr)
	}

	var result map[string]interface{}
	var queryFunc func(string) (map[string]interface{}, error)

	if mode == "parse" {
		queryFunc = client.QueryParseTask
	} else {
		queryFunc = client.QueryExtractTask
	}

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
// If it's a file, reads and returns the URLs from the file; otherwise returns it as a single URL.
func resolveURLInput(input string) []string {
	info, err := os.Stat(input)
	if err == nil && !info.IsDir() {
		// Input is an existing file — read URLs from it
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
	// Not a file — treat as a single URL
	return []string{input}
}
