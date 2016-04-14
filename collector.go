package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

var metrics_dir, scripts_dir, time_prefix string
var enable_scripts int
var debug bool
var core_stats map[string]interface{}

func main() {

	v := flag.Bool("v", false, "prints current version and exits")
	flag.StringVar(&metrics_dir, "m", "/var/log/stats-ag", "Location where metrics log files are written")
	flag.StringVar(&scripts_dir, "s", "/opt/stats-ag/scripts", "Location where custom metrics scripts are located")
	flag.StringVar(&time_prefix, "p", "SYSLOG", "Date prefix format for metric entries (RFC822Z, ISO8601, RFC3339, SYSLOG)")
	flag.BoolVar(&debug, "d", false, "Enable verbose debug mode")
	flag.Parse()

	if *v {
		fmt.Printf("Stats-ag Version %s (Build date: %s)\nCommit SHA: %s\n", VERSION, BUILD_DATE, COMMIT_SHA)
		os.Exit(0)
	}

	if flag.NFlag() == 0 {
		flag.PrintDefaults()
	}

	if debug {
		fmt.Printf(
			"\nstats-ag config values:\n---------------------------\nmetrics_dir = %s\nscripts_dir = %s\ntime_prefix = %s\ndebug = %v\n\n",
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

	wg.Add(len(core_stats))

	_, dir_exists_err := os.Stat(scripts_dir)
	if dir_exists_err == nil {
		enable_scripts = 1
		wg.Add(len(scripts))
	} else {
		enable_scripts = 0
	}

	for k, _ := range core_stats {
		go func(wg *sync.WaitGroup, core_stats map[string]interface{}, method string) {
			mw := NewMetricsWriter(method, time_prefix)
			if debug {
				Log(DEBUG_LVL_NOTICE, fmt.Sprintf("Fetching core stat: %s", method))
			}
			res, _ := Call(core_stats, method)
			mw.Save(res[0].String())
			wg.Done()
		}(&wg, core_stats, k)
	}

	if enable_scripts == 1 {

		for _, src := range scripts {

			go func(wg *sync.WaitGroup, script_name string) {

				if debug {
					Log(DEBUG_LVL_NOTICE, fmt.Sprintf("Calling custom stat: %s", script_name))
				}
				cm, err := NewCustomMetric(script_name)
				if err != nil {
					Log(DEBUG_ERROR, fmt.Sprintf("NewCustomMetric() error: %s", err))
					wg.Done()
					return
				}

				mw := NewMetricsWriter(cm.Name, time_prefix)
				if stats, err := cm.GetStats(); err == nil {
					mw.Save(stats)
				} else {
					Log(DEBUG_ERROR, fmt.Sprintf("GetStats() error: %s", err))
				}
				wg.Done()
			}(&wg, strings.Join([]string{strings.TrimRight(scripts_dir, "/"), src.Name()}, "/"))
		}

	}

	wg.Wait()

}
