package mylog

import (
	"fmt"
	"runtime"
)
func LogError(err error) {
    if err != nil {
        // Retrieve the caller information
        _, file, line, ok := runtime.Caller(1)
        if ok {
            fmt.Printf("Error: %v (file: %s, line: %d)\n", err, file, line)
        } else {
            fmt.Printf("Error: %v\n", err)
        }
    }
}