package helper

import (
	"runtime"
	"strings"
)

func GetCallerInfo(needFuncName bool) (string, string, int) {
	pc, file, line, _ := runtime.Caller(2)

	temp := strings.Split(file, "/")
	file = temp[len(temp)-1]

	var funcName string
	if needFuncName {
		temp = strings.Split(runtime.FuncForPC(pc).Name(), ".")
		funcName = temp[len(temp)-1]
	}

	return funcName, file, line
}
