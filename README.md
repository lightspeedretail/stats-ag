
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

## Command options

- `-v`: Show application version and exit
- `-d`: Enable verbose debug mode (default: false)
- `-m`: Location where metrics log files are written (default: /var/log/stats-ag)
- `-p`: Date prefix format for metric entries, either RFC822Z, ISO8601, RFC3339 or SYSLOG (default: SYSLOG)
- `-s`: Location where custom metrics scripts are located (default: /opt/stats-ag/scripts)


## Usage examples

Calling the script manually:

```bash
stats-ag -m /var/log/stats-ag/metrics/ -s /opt/stats-ag/scripts/ -p ISO8601
```

Or having a cron run it every minute:

```bash
* * * * * /opt/stats-ag/stats-ag -m /var/log/stats-ag/metrics/ -s /opt/stats-ag/scripts/ -p ISO8601
```

## Authors

Alain Lefebvre  (alain.lefebvre '**at**' lightspeedhq.com)


## License

Convered under the MIT License