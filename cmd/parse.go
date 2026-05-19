package cmd

import (
	"fmt"
	"os"
	"path/filepath"

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
		noWait, _ := cmd.Flags().GetBool("no-wait")
		export, _ := cmd.Flags().GetString("export")
		timeout, _ := cmd.Flags().GetInt("timeout")
		concurrency, _ := cmd.Flags().GetInt("concurrency")
		retry, _ := cmd.Flags().GetInt("retry")

		processLocalFiles(path, appID, asyncMode, noWait, export, timeout, concurrency, retry, "parse")
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
		noWait, _ := cmd.Flags().GetBool("no-wait")
		export, _ := cmd.Flags().GetString("export")
		timeout, _ := cmd.Flags().GetInt("timeout")
		concurrency, _ := cmd.Flags().GetInt("concurrency")
		retry, _ := cmd.Flags().GetInt("retry")

		urls := resolveURLInput(urlArg)
		processURLs(urls, appID, asyncMode, noWait, export, timeout, concurrency, retry, "parse")
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
		noWait, _ := cmd.Flags().GetBool("no-wait")
		export, _ := cmd.Flags().GetString("export")
		timeout, _ := cmd.Flags().GetInt("timeout")
		fileName, _ := cmd.Flags().GetString("file-name")
		concurrency, _ := cmd.Flags().GetInt("concurrency")
		retry, _ := cmd.Flags().GetInt("retry")

		processBase64(args, appID, fileName, asyncMode, noWait, export, timeout, concurrency, retry, "parse")
	},
}

// parseQueryCmd represents the parse query command
var parseQueryCmd = &cobra.Command{
	Use:   "query",
	Short: i18n.T("parse_query_title"),
	Long:  i18n.T("parse_query_title"),
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		watch, _ := cmd.Flags().GetBool("watch")
		taskFile, _ := cmd.Flags().GetString("file")
		export, _ := cmd.Flags().GetString("export")
		timeout, _ := cmd.Flags().GetInt("timeout")
		concurrency, _ := cmd.Flags().GetInt("concurrency")

		if taskFile == "" && len(args) == 0 {
			formatterOut.ExitWithError(errors.NewCLIError(
				"Either task IDs or --file must be provided",
				errors.ErrorTypeParam,
				errors.ExitParameterError,
				false,
				"Provide task IDs as arguments or use --file to load from a JSON file.",
				nil,
			))
		}

		queryTasks(args, watch, taskFile, export, timeout, concurrency, "parse")
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
	parseLocalCmd.Flags().Bool("no-wait", false, i18n.T("option_no_wait"))
	parseLocalCmd.Flags().String("export", "", i18n.T("option_export"))
	parseLocalCmd.Flags().Int("timeout", 900, i18n.T("option_timeout"))
	parseLocalCmd.Flags().Int("concurrency", 1, i18n.T("option_concurrency"))
	parseLocalCmd.Flags().Int("retry", 0, i18n.T("option_retry"))
	parseLocalCmd.MarkFlagRequired("app-id")

	// parse url flags
	parseURLCmd.Flags().String("app-id", "", i18n.T("option_app_id_parse"))
	parseURLCmd.Flags().Bool("async", false, i18n.T("option_async"))
	parseURLCmd.Flags().Bool("no-wait", false, i18n.T("option_no_wait"))
	parseURLCmd.Flags().String("export", "", i18n.T("option_export"))
	parseURLCmd.Flags().Int("timeout", 900, i18n.T("option_timeout"))
	parseURLCmd.Flags().Int("concurrency", 1, i18n.T("option_concurrency"))
	parseURLCmd.Flags().Int("retry", 0, i18n.T("option_retry"))
	parseURLCmd.MarkFlagRequired("app-id")

	// parse base64 flags
	parseBase64Cmd.Flags().String("app-id", "", i18n.T("option_app_id_parse"))
	parseBase64Cmd.Flags().Bool("async", false, i18n.T("option_async"))
	parseBase64Cmd.Flags().Bool("no-wait", false, i18n.T("option_no_wait"))
	parseBase64Cmd.Flags().String("export", "", i18n.T("option_export"))
	parseBase64Cmd.Flags().Int("timeout", 900, i18n.T("option_timeout"))
	parseBase64Cmd.Flags().String("file-name", "document", i18n.T("option_file_name"))
	parseBase64Cmd.Flags().Int("concurrency", 1, i18n.T("option_concurrency"))
	parseBase64Cmd.Flags().Int("retry", 0, i18n.T("option_retry"))
	parseBase64Cmd.MarkFlagRequired("app-id")

	// parse query flags
	parseQueryCmd.Flags().Bool("watch", false, i18n.T("option_watch"))
	parseQueryCmd.Flags().String("file", "", i18n.T("option_task_file"))
	parseQueryCmd.Flags().String("export", "", i18n.T("option_export"))
	parseQueryCmd.Flags().Int("timeout", 900, i18n.T("option_watch_timeout"))
	parseQueryCmd.Flags().Int("concurrency", 1, i18n.T("option_concurrency"))
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

func processLocalFiles(pathStr, appID string, asyncMode, noWait bool, exportPath string, timeout, concurrency, maxRetry int, mode string) {
	client, _ := initClientWithConfig(mode + " local")
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
		batchSubmit(jobs, concurrency, submitFn, exportPath)
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
		// extract
		if asyncMode {
			taskID, err := client.ExtractAsync("", appID, filePath, "", fileName, nil)
			if err != nil {
				return nil, err
			}
			formatterOut.PrintInfo(fmt.Sprintf("%s (task_id: %s)", fileName, taskID))
			return client.WaitForTask(taskID, client.QueryExtractTask, timeout, 2)
		}
		return client.ExtractSync("", appID, filePath, "", fileName, nil)
	}

	results := batchProcess(jobs, concurrency, maxRetry, processFn, exportPath)
	exitWithBatchResults(results)
}

func processURLs(urls []string, appID string, asyncMode, noWait bool, exportPath string, timeout, concurrency, maxRetry int, mode string) {
	client, _ := initClientWithConfig(mode + " url")
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

	// noWait mode: submit only, output task IDs
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
		if asyncMode {
			taskID, err := client.ExtractAsync(url, appID, "", "", "", nil)
			if err != nil {
				return nil, err
			}
			formatterOut.PrintInfo(fmt.Sprintf("%s (task_id: %s)", job.Name, taskID))
			return client.WaitForTask(taskID, client.QueryExtractTask, timeout, 2)
		}
		return client.ExtractSync(url, appID, "", "", "", nil)
	}

	results := batchProcess(jobs, concurrency, maxRetry, processFn, exportPath)
	exitWithBatchResults(results)
}

func processBase64(b64Strings []string, appID, fileName string, asyncMode, noWait bool, exportPath string, timeout, concurrency, maxRetry int, mode string) {
	client, _ := initClientWithConfig(mode + " base64")
		// Build batch jobs
	jobs := make([]BatchJob, len(b64Strings))
	for i := range b64Strings {
		name := fileName
		if len(b64Strings) > 1 {
			name = fmt.Sprintf("%s_%d", fileName, i+1)
		}
		jobs[i] = BatchJob{Index: i, Name: name}
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
		batchSubmit(jobs, concurrency, submitFn, exportPath)
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
		if asyncMode {
			taskID, err := client.ExtractAsync("", appID, "", b64, currentFileName, nil)
			if err != nil {
				return nil, err
			}
			formatterOut.PrintInfo(fmt.Sprintf("%s (task_id: %s)", currentFileName, taskID))
			return client.WaitForTask(taskID, client.QueryExtractTask, timeout, 2)
		}
		return client.ExtractSync("", appID, "", b64, currentFileName, nil)
	}

	results := batchProcess(jobs, concurrency, maxRetry, processFn, exportPath)
	exitWithBatchResults(results)
}

func queryTasks(taskIDs []string, watch bool, taskFile, exportPath string, timeout, concurrency int, mode string) {
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

	if len(taskIDs) == 1 {
		// Single task query: preserve existing behavior
		querySingleTask(client, taskIDs[0], watch, exportPath, timeout, mode)
		return
	}

	// Multi task query: use batch processing
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
