package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

var metrics_dir, scripts_dir, time_prefix string
var enable_scripts, debug int
var core_stats map[string]interface{}

func main() {

	flag.IntVar(&enable_scripts, "e", 0, "Enable custom scripts execution")
	flag.StringVar(&metrics_dir, "m", "/var/log/stats-ag", "Location where metrics log files are written")
	flag.StringVar(&scripts_dir, "s", "/opt/stats-ag/scripts", "Location where custom metrics scripts are located")
	flag.StringVar(&time_prefix, "p", "SYSLOG", "Date prefix format for metric entries (RFC822Z, ISO8601, RFC3339, SYSLOG)")
	flag.IntVar(&debug, "d", 0, "Enable verbose debug mode")

	if len(os.Args) >= 2 && os.Args[1] == "-v" {
		fmt.Printf("Stats-ag Version %s (Build date: %s)\nCommit SHA: %s\n", VERSION, BUILD_DATE, COMMIT_SHA)
		os.Exit(0)
	}

	flag.Parse()

	if flag.NFlag() == 0 {
		flag.PrintDefaults()
	}

	if debug == 1 {
		fmt.Printf(
			"\nstats-ag config values:\n---------------------------\nenable_scripts = %d\nmetrics_dir = %s\nscripts_dir = %s\ntime_prefix = %s\ndebug = %d\n\n",
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
			if debug == 1 {
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

				if debug == 1 {
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
