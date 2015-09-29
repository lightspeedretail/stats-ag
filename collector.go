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

func main() {

	flag.StringVar(&metrics_dir, "m", "/var/log/stats_collector", "Location where metrics log files are written")
	flag.IntVar(&enable_scripts, "e", 0, "Enable custom scripts execution")
	flag.StringVar(&scripts_dir, "s", "/opt/stats_collector", "Location where custom metrics scripts are located")
	flag.StringVar(&time_prefix, "p", "SYSLOG", "Date prefix format for metric entries (RFC822Z, ISO8601, RFC3339, SYSLOG)")
	flag.Parse()

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
			res, _ := Call(core_stats, method)
			mw.Save(res[0].String())
			wg.Done()
		}(&wg, core_stats, k)
	}

	_, dir_exists_err := os.Stat(scripts_dir)
	if enable_scripts == 1 && dir_exists_err == nil {

		for _, src := range scripts {

			go func(wg *sync.WaitGroup, script_name string) {

				cm, err := NewCustomMetric(script_name)
				if err != nil {
					fmt.Println("Error:", err)
					wg.Done()
					return
				}

				mw := NewMetricsWriter(cm.Name, time_prefix)
				if stats, err := cm.GetStats(); err == nil {
					mw.Save(stats)
				} else {
					fmt.Println("GetStats() error: ", err)
				}
				wg.Done()
			}(&wg, scripts_dir+src.Name())
		}

	}

	wg.Wait()

}
