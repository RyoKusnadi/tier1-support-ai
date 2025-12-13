package logger

import (
	"fmt"
	"log"
)

func Info(msg string, fields map[string]interface{}) {
	log.Println(format("INFO", msg, fields))
}

func Error(msg string, fields map[string]interface{}) {
	log.Println(format("ERROR", msg, fields))
}

func format(level, msg string, fields map[string]interface{}) string {
	out := "level=" + level + " msg=\"" + msg + "\""
	for k, v := range fields {
		out += " " + k + "=" + toString(v)
	}
	return out
}

func toString(v interface{}) string {
	return "\"" + logValue(v) + "\""
}

func logValue(v interface{}) string {
	return fmt.Sprintf("%v", v)
}
