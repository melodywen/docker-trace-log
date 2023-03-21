package helper

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// 在函数中获取调用者
func getCallerInfo(needFuncName bool) (string, string, int) {
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

// EnterExitFunc 打印函数进出日志:
//     使用方法
// 	defer EnterExitFunc()()
func EnterExitFunc() func() {
	funcName, file, line := getCallerInfo(true)
	start := time.Now()

	fmt.Printf("enter %s func (%s:%d)", funcName, file, line)
	return func() {
		_, file, line = getCallerInfo(false)
		fmt.Printf("exit %s (%s) func (%s:%d)", funcName, time.Since(start), file, line)
	}
}
