package main

import "log"

func logVerbose(msg string) {
	if displayVerboseLogs {
		log.Println(msg)
	}
}
