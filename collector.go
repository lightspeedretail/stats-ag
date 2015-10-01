package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

var metrics_dir, scripts_dir, time_prefix string
var enable_scripts int
var core_stats map[string]interface{}
var debug bool

func main() {

	if os.Args[1] == "-v" {
		fmt.Printf("Stats-ag Version %s\n", VERSION)
		os.Exit(0)
	}

	flag.IntVar(&enable_scripts, "e", 0, "Enable custom scripts execution")
	flag.StringVar(&metrics_dir, "m", "/var/log/stats_collector", "Location where metrics log files are written")
	flag.StringVar(&scripts_dir, "s", "/opt/stats_collector", "Location where custom metrics scripts are located")
	flag.StringVar(&time_prefix, "p", "SYSLOG", "Date prefix format for metric entries (RFC822Z, ISO8601, RFC3339, SYSLOG)")
	flag.BoolVar(&debug, "d", false, "Enable verbose debug mode")
	flag.Parse()

	if flag.NArg() == 0 {
		usage := `
Usage: stats-ag [OPTIONS]
Options:
	-e [ENABLE_CUSTOM_SCRIPTS] (default = 0)
	-m [METRICS_DIR] (default = /var/log/stats_collector)
	-s [CUSTOM_SCRIPTS_DIR] (default = /opt/stats_collector)
	-p [TIME_PREFIX_FORMAT] (default = SYSLOG)
	-d [DEBUG] (default = false)
	`
		fmt.Printf("%s\n", usage)
		os.Exit(0)
	}

	if debug {
		fmt.Printf(
			"\nstats-ag config values:\n---------------------------\nenable_scripts = %d\nmetrics_dir = %s\nscripts_dir = %s\ntime_prefix = %s\ndebug = %t\n\n",
			enable_scripts,
			metrics_dir,
			scripts_dir,
			time_prefix,
			debug,
		)
	}

	var wg sync.WaitGroup
	scripts, _ := ioutil.ReadDir(scripts_dir)

	core_stats = map[string]interface{}{
		"memory": GetMemStats,
		"load":   GetLoadStats,
		"disk":   GetDiskStats,
		"cpu":    GetCpuStats,
		"host":   GetHostStats,
		"net":    GetNetStats,
	}

	wg.Add(len(core_stats) + len(scripts))

	for k, _ := range core_stats {
		go func(wg *sync.WaitGroup, core_stats map[string]interface{}, method string) {
			mw := NewMetricsWriter(method, time_prefix)
			if debug {
				fmt.Printf("%s [DEBUG] fetching core stat: %s\n", getDateStamp(time_prefix), method)
			}
			res, _ := Call(core_stats, method)
			mw.Save(res[0].String())
			wg.Done()
		}(&wg, core_stats, k)
	}

	_, dir_exists_err := os.Stat(scripts_dir)
	if enable_scripts == 1 && dir_exists_err == nil {

		for _, src := range scripts {

			go func(wg *sync.WaitGroup, script_name string) {

				if debug {
					fmt.Printf("%s [DEBUG] calling custom stat: %s\n", getDateStamp(time_prefix), script_name)
				}
				cm, err := NewCustomMetric(script_name)
				if err != nil {
					fmt.Printf("%s [ERROR] NewCustomMetric() error: %s\n", getDateStamp(time_prefix), err)

					wg.Done()
					return
				}

				mw := NewMetricsWriter(cm.Name, time_prefix)
				if stats, err := cm.GetStats(); err == nil {
					mw.Save(stats)
				} else {
					fmt.Printf("%s [ERROR] GetStats() error: %s\n", getDateStamp(time_prefix), err)
				}
				wg.Done()
			}(&wg, scripts_dir+src.Name())
		}

	}

	wg.Wait()

}
