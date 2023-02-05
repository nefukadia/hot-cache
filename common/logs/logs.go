package logs

import (
	"log"
	"runtime/debug"
)

func init() {

}

func Info(v ...any) {
	values := []any{"[INFO]"}
	values = append(values, v...)
	log.Println(values...)
}

func Warning(v ...any) {
	values := []any{"[WARNING]"}
	v = append(v, string(debug.Stack()))
	values = append(values, v...)
	log.Println(values...)
}

func WarningWithoutStack(v ...any) {
	values := []any{"[WARNING]"}
	values = append(values, v...)
	log.Println(values...)
}

func Error(v ...any) {
	values := []any{"[ERROR]"}
	v = append(v, string(debug.Stack()))
	values = append(values, v...)
	log.Println(values...)
}
