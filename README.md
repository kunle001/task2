# File Modification Tracker

## Overview

The **File Modification Tracker** is a tool for monitoring and tracking changes in a specified directory. It features a user-friendly interface for managing the service and provides HTTP endpoints for health checks, log retrieval, and command execution. The application is compatible with macOS and Windows.

## Features

- Monitors file modifications in a specified directory.
- Runs as a Windows service or macOS application.
- Provides a user interface for controlling the service.
- Exposes HTTP endpoints for health checks, log retrieval, and command execution.

## Setup

### Prerequisites

- **Go 1.18 or later**
- **`pkgbuild`** (for macOS packaging)
- **`go mod tidy`** (for managing Go dependencies)

### Configuration

The application reads configuration from a `config.yaml` file in the project directory. Include the following fields:

- `directory`: Directory to monitor for file modifications.
- `check_frequency`: Interval (in seconds) to check for changes.
- `api_endpoint`: URL for reporting or retrieving data.

**Note:** Ensure the `directory` is a valid path on your machine!!!!!!.

#### Example Configuration

```yaml
directory: /path/to/monitor
check_frequency: 60
api_endpoint: http://example.com/api
```

### Building the Application

To build the macOS binary, run:

```bash
make build
```

This compiles the Go code and places the binary in the `bin` directory.

### Running the Application

To run the application directly from the source, use:

```bash
make run
```

This starts the application.

### User Interface

The UI includes three buttons:

- **Start**: Begins monitoring by starting the timer and worker threads.
- **Stop**: Stops the monitoring service by halting the worker and timer threads.
- **Fetch Logs**: Retrieves all logs.

### HTTP Endpoints

- **Health Check**: `localhost:8080/health`
  Checks the health status of the service.

- **Logs**: `localhost:8080/logs`
  Retrieves all modification logs for the specified directory.
  Sample Response

  ```json
  {
    "status": true,
    "data": [
      {
        "timestamp": "2024-09-16T07:25:27Z", //Time at which the stats was fetched.
        "fileStats": [
          //array of files in the directory and their respective logs.
          {
            "path": "/Users/apple/Desktop/Work/Me/Blockchain/First/simpleStorage.sol",
            "size": "556",
            "last_accessed": "2024-04-30T11:11:44+01:00",
            "last_modified": "2024-04-30T11:09:34+01:00",
            "last_changed": "2024-04-30T11:09:34+01:00"
          },
          {
            "path": "/Users/apple/Desktop/Work/Me/Blockchain/First/package.json",
            "size": "201",
            "last_accessed": "2024-05-12T15:08:10+01:00",
            "last_modified": "2024-05-12T15:08:06+01:00",
            "last_changed": "2024-05-12T15:08:06+01:00"
          }
        ]
      }
    ]
  }
  ```

- **Commands**: `localhost:8080/command`
  A POST endpoint where you can send commands for the worker thread to execute. Example payload:

  ```json
  {
    "command": [
      "ls -l", // List files in the current directory with detailed info
      "cat /path/to/file.txt", // Display contents of a specific file
      "grep 'search_term' /path/to/file", // Search for a term within a file
      "mkdir /path/to/new_directory", // Create a new directory
      "cp /path/to/source /path/to/destination" // Copy a file or directory
    ]
  }
  ```

  NB: Dangerous commands such as "rm", "del", "unlink", "rmdir", "erase", "destroy", are not allowed

### Running Tests

To execute unit tests, use:

```bash
make test
```

This runs all defined tests and displays the results.

### Creating a macOS Package

To package the application into a `.pkg` file for macOS, use:

```bash
make package
```

The `.pkg` file will be located in the root directory.

### Cleaning Up

To remove generated files and clean the project directory, run:

```bash
make clean
```

This deletes the binary, packaging directories, and `.pkg` file if present.

### Installing the Application

To install the application on macOS:

1. Execute the `.pkg` file created by `make package`.
2. Follow the prompts to complete the installation.
3. The application will be installed in `/usr/local/bin` by default.

## Contact

For questions or support, please contact [adekunle.olanipekun.ko@gmail.com].
