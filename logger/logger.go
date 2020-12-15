package logger

import (
	"log"
	"os"
	"strings"
)

var logger *log.Logger

func SetDefault(l *log.Logger) {
	logger = l
}

func Info(args ...interface{}) {
	logger.Println(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Printf(format, args...)
}

func Init() {
	path := "logs"
	mode := os.ModePerm
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, mode)
	}
	file, _ := os.Create(strings.Join([]string{path, "log.txt"}, "/"))
	defer file.Close()
	loger := log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
	SetDefault(loger)
}
