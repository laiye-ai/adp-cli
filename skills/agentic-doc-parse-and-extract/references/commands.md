# Complete Command List of ADP CLI 

```bash
# Configuration Management
adp config get # View current configuration
adp config set --api-key YOUR_API_KEY # Set API Key
adp config set --api-base-url API_BASE_URL # Set API Base URL
adp config clear # Clear configuration 

# Application Management
adp app-id list # Query the list of available applications 

# Document Parsing
adp parse url <file_url> --app-id <app_id> # Synchronous processing for URL document parsing
adp parse local <file_path> --app-id <app_id> # Synchronous processing for local document parsing
adp parse url <file_url> --app-id <app_id> --async # Asynchronous processing for URL document parsing
adp parse local <file_path> --app-id <app_id> --async # Asynchronous processing for local document parsing 

# Document Extraction
adp extract url <file_url> --app-id <app_id> # Synchronous processing for extracting URL documents
adp extract local <file_path> --app-id <app_id> # Synchronous processing for extracting local documents
adp extract url <file_url> --app-id <app_id> --async # Asynchronous processing for extracting URL documents
adp extract local <file_path> --app-id <app_id> --async # Asynchronous processing for extracting local documents 

# Asynchronous Task Query
adp parse query <task_id> # Query to check the status of the parsing task
adp extract query <task_id> # Query to check the status of the extraction task 

# Create a custom extraction application
adp custom-app create --app-name <app-name> --extract-fields <json> --parse-mode <mode> --enable-long-doc <bool> --long-doc-config <json> # Create a custom extraction application
adp custom-app get-config --api-key "your_api_key" --app-id "app_id" --config-version "v1" # View the configuration of the custom application
adp custom-app delete --api-key "your_api_key" --app-id "app_id" # Delete the custom application
adp custom-app delete-version --api-key "your_api_key" --app-id "app_id" --config-version "v2" # Delete the specified version of the custom application
adp custom-app ai-generate --app-id "app_id" --file-url "file_url" # AI generates field recommendations 

# Batch Processing
adp parse local <folder path> --app-id <app_id> --export <folder path> --concurrency <concurrency number> # Batch parsing of documents in the local folder
adp extract local <folder path> --app-id <app_id> --export <folder path> --concurrency <concurrency number> # Batch extraction of documents in the local folder
adp parse url <URL list file path> --export <folder path> --app-id <app_id> --concurrency <concurrency number> # Batch parsing of documents in the URL list file
adp extract url <URL list file path> --export <folder path> --app-id <app_id> --concurrency <concurrency number> # Batch extraction of documents in the URL list file 

# Assistance
adp --help # View the complete list of commands and usage instructions
adp credit # Check the current account's point balance
adp app-id cache # View the cached application id
```

# ADP Error Code
