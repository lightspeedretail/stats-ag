package main

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

func metricsFileName(series_name string) string {
	return metrics_dir + "/" + series_name + ".log"
}

func getDateStamp(format string) string {

	t := time.Now().UTC()

	if format == "RFC822Z" {

		return t.Format(time.RFC822Z)

	} else if format == "ISO8601" { // 2011-04-19T03:44:01.103Z

		return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.%03dZ",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second(), (t.Nanosecond() / 10000000))

	} else if format == "RFC3339" {

		return t.Format(time.RFC3339)

	} else if format == "SYSLOG" { // Oct  1 13:22:25

		return fmt.Sprintf("%s %d %02d:%02d:%02d",
			t.Month().String()[0:3], t.Day(),
			t.Hour(), t.Minute(), t.Second())
	}
	return t.Format(time.RFC822Z)
}

func Call(m map[string]interface{}, name string, params ...interface{}) (result []reflect.Value, err error) {
	f := reflect.ValueOf(m[name])
	if len(params) != f.Type().NumIn() {
		err = errors.New("Wrong number of parameters!")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return
}
