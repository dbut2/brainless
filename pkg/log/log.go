package log

import (
	"fmt"
	"runtime"
)

func Dump(w ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	fmt.Printf("Dumped at %s:%d", file, line)
	for _, v := range w {
		fmt.Printf(" %+v", v)
	}
	fmt.Printf("\n")
}
