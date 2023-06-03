package main

import (
	"log"

	"github.com/fatih/color"
)

func SendLog(format string, v ...any) {
	log.Printf(color.CyanString(format), v...)
}

func RecvLog(format string, v ...any) {
	log.Printf(color.GreenString(format), v...)
}

func ErrLog(format string, v ...any) {
	log.Printf(color.RedString(format), v...)
}
