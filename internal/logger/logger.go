package logger

import (
	"log"
	"os"
)

var err_log *log.Logger
var debug_log *log.Logger
var info_log *log.Logger

func init() {
	err_log = log.New(os.Stdout, "[Error]", log.LstdFlags)
	debug_log = log.New(os.Stdout, "[Debug]", log.LstdFlags)
	info_log = log.New(os.Stdout, "[Info]", log.LstdFlags)
}

func Info(v ...any) {
	info_log.Println(v...)
}

func Debug(v ...any) {
	debug_log.Println(v...)
}

func Error(v ...any) {
	err_log.Println(v...)
}
