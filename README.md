
## Description

A small process that collects and outputs system metrics to designated log files.  Also runs executable scripts, returning simple data, which are placed in the appropriate location.

## Building

To build the binary, use the following example:

```bash
go build -ldflags "-X main.build_date=`date +%Y-%m-%d` -X main.VERSION=X.Y.Z -X main.COMMIT_SHA=`git rev-parse --verify HEAD`" -o stats-ag
```

Or use this method to build for another OS/Architecture
```bash
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.BUILD_DATE=`date +%Y-%m-%d` -X main.VERSION=X.Y.Z -X main.COMMIT_SHA=`git rev-parse --verify HEAD`" -o stats-ag
```


## Usage examples


## Authors



## License

Convered under the MIT License