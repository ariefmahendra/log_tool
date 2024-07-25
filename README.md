# LOG TOOL 
## Description
This Tools is used to read log file on server, download log file, filter log file based on the given keyword, read latest log file on server with buffer size, and read current ENV file use in this application

## Installation
### Prerequisite
- Go 1.16 or higher
- Git
- ENV must be set in the .env file
```text
### Step
```bash
git clone https://github.com/ariefmahendra/log_tool.git
cd log_tool
go install 
go run main.go
```

## Features 
- [x] Read List folder or file on server
- [x] Download Log From Server 
- [x] Filter log based on the given keyword
- [x] Read latest log file on server with buffer size
- [x] Read Current ENV file use in this application

## Rule Of Format Log File
- Log file must be in .log format
- When you want to filter log file, you must provide the keyword that you want to filter or if you not provide the keyword, it will show all log file by type request or response 
- Log will be print on the console only start from the Type and stop before the DEBUG or ERROR LOG

## EXAMPLE
- BEFORE
```text
DEBUG 2021-07-01 12:00:00.000 [main] - This is debug log
Type: Request (Or Response)
Your Request Or Response Log
DEBUG (OR ERROR) 2021-07-01 12:00:00.000 [main] - This is debug log
```
- AFTER
```text
🕒 TIME : Log Time 🕒
Type: Request
Your Request Or Response Log
```

