package assert

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"runtime/debug"
)

func runAssert(msg string, args ...interface{}) {
	slogValues := []interface{}{
		"msg",
		msg,
		"area",
		"Assert",
	}
	slogValues = append(slogValues, args...)
	fmt.Fprintf(os.Stderr, "ARGS: %+v\n", args)
	fmt.Fprintf(os.Stderr, "ASSERT\n")
	for i := 0; i < len(slogValues); i += 2 {
		fmt.Fprintf(os.Stderr, "   %s=%v\n", slogValues[i], slogValues[i+1])
	}
	fmt.Fprintln(os.Stderr, string(debug.Stack()))
	os.Exit(1)
}

func Assert(truth bool, msg string, data ...any) {
	if !truth {
		runAssert(msg, data...)
	}
}

func Nil(item any, msg string, data ...any) {
	slog.Info("Nil Check", "item", item)
	if item == nil {
		return
	}
	slog.Error("Nil#not nil encountered")
	runAssert(msg, data...)
}

func NotNil(item any, msg string, data ...any) {
	if item == nil || reflect.ValueOf(item).Kind() == reflect.Ptr && reflect.ValueOf(item).IsNil() {
		slog.Error("NotNil#nil encountered")
		runAssert(msg, data...)
	}
}

func Never(msg string, data ...any) {
	runAssert(msg, data...)
}

func NoError(err error, msg string, data ...any) {
	if err != nil {
		data = append(data, "error", err)
		runAssert(msg, data...)
	}
}
