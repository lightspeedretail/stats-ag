stats_ag CHANGELOG
===========================


0.1.2
-----
- Updated `gopsutil` related fuction calls to reflect the 2.0.0 [release](https://github.com/shirou/gopsutil/releases/tag/v2.0.0)
- [bugfix] If the scripts dir was specified, but scripts were not enabled, a deadlock situation would occur as the application was expecting more stats to execute then actually existed. 
- Removed the `-e` flag to enable custome scripts as the `-s` flag with a default value of empty, makes it unecessary.  

0.1.1
-----
- Fixed syslog date format and example
- Added debug option `-d`
- Updated how flag default options are printed
- Added `-v` option to print version info
- Updated build script so that build version is automatically included as well as the build date and the commit hash

0.1.0
-----
- Initial version. 
