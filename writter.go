package main

import (
	"fmt"
	"os"
	"sync"
)

type MetricsWriter struct {
	fileLock   *sync.Mutex
	GroupName  string
	timePrefix string
}

func NewMetricsWriter(series_name string, time_prefix string) *MetricsWriter {
	mw := &MetricsWriter{
		GroupName:  series_name,
		fileLock:   &sync.Mutex{},
		timePrefix: time_prefix,
	}
	// Create the dir if it doesn't exist
	if _, err := os.Stat(metrics_dir); os.IsNotExist(err) {
		err = os.Mkdir(metrics_dir, 0755)
	}
	return mw
}

func (mw *MetricsWriter) Save(data string) int {

	mw.fileLock.Lock()

	fh, err := os.OpenFile(metricsFileName(mw.GroupName), os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fh, err = os.Create(metricsFileName(mw.GroupName))
		if err != nil {
			Log(DEBUG_ERROR, fmt.Sprintf("Error creating file %s: %s", metricsFileName(mw.GroupName), err))
		}
	}
	defer fh.Close()

	nb, err := fh.WriteString(getDateStamp(mw.timePrefix) + " " + data + "\n")
	if err != nil {
		Log(DEBUG_ERROR, fmt.Sprintf("Error writting to %s: %s", metricsFileName(mw.GroupName), err))
	} else if debug {
		Log(DEBUG_LVL_NOTICE, fmt.Sprintf("Wrote %d bytes to %s", nb, metricsFileName(mw.GroupName)))
	}

	mw.fileLock.Unlock()
	return nb
}
