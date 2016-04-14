package main

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	DEBUG_LVL_NOTICE = iota + 1
	DEBUG_LVL_VERBOSE
	DEBUG_LVL_DEBUG
	DEBUG_ERROR
)

func metricsFileName(series_name string) string {
	return strings.Join([]string{strings.TrimRight(metrics_dir, "/"), (series_name + ".log")}, "/")
}

func YmdToString() string {
	t := time.Now()
	y, m, d := t.Date()
	return strconv.Itoa(y) + "-" + fmt.Sprintf("%02d", m) + "-" + fmt.Sprintf("%02d", d)
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

func DateStampAsString(with_milliseconds bool) string {
	t := time.Now()
	dt := YmdToString() + " " + fmt.Sprintf("%02d", t.Hour()) + ":" + fmt.Sprintf("%02d", t.Minute()) + ":" + fmt.Sprintf("%02d", t.Second())
	if with_milliseconds {
		dt = dt + "." + fmt.Sprintf("%03d", t.Nanosecond()/1000000)
	}
	return dt
}

func Log(lvl int, str string) {

	switch lvl {
	case 1:
		fmt.Printf("[%s] NOTICE: %s\n", DateStampAsString(true), str)
	case 2:
		fmt.Printf("[%s] DEBUG: %s\n", DateStampAsString(true), str)
	case 3:
		fmt.Printf("[%s] VERBOSE: %s\n", DateStampAsString(true), str)
	case 4:
		fmt.Printf("[%s] ERROR: %s\n", DateStampAsString(true), str)
	}
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
