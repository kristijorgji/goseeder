package goseeder

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

const (
	infoColor    = "\033[1;34m%s\033[0m"
	noticeColor  = "\033[1;36m%s\033[0m"
	warningColor = "\033[1;33m%s\033[0m"
	errorColor   = "\033[1;31m%s\033[0m"
	debugColor   = "\033[0;36m%s\033[0m"
)

var printError = color(errorColor)

func color(colorString string) func(...interface{}) string {
	return func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
}

func findString(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func prepareStatement(table string, row map[string]string) (strings.Builder, []interface{}) {
	var left strings.Builder
	var right strings.Builder

	var args []interface{}

	left.WriteString(fmt.Sprintf("insert into %s (", table))
	right.WriteString("values (")

	i := 0

	for k, v := range row {
		if i == 0 {
			left.WriteString(k)
			right.WriteString("?")
		} else {
			left.WriteString(fmt.Sprintf(", %s", k))
			right.WriteString(", ?")
		}

		args = append(args, parseValue(v))
		i++
	}

	left.WriteString(") ")
	right.WriteString(")")
	left.WriteString(right.String())

	return left, args
}

func parseValue(value string) interface{} {
	if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
		return parsed
	}
	if parsed, err := strconv.ParseFloat(value, 32); err == nil {
		return parsed
	}
	if parsed, err := strconv.ParseBool(value); err == nil {
		return parsed
	}

	return value
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
