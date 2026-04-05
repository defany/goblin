package rt

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

func FnName(skip ...int) string {
	sk := 1
	if len(skip) > 0 {
		sk = skip[0] + 1
	}

	pc, _, _, ok := runtime.Caller(sk)
	if !ok {
		return "UNKNOWN"
	}

	funcName := runtime.FuncForPC(pc).Name()

	s := strings.Split(funcName, ".")
	if len(s) < 2 {
		return s[len(s)-1]
	}

	fnName := s[len(s)-1]
	pkg := strings.NewReplacer("(", "", ")", "", "*", "").Replace(s[len(s)-2])

	return fmt.Sprintf("%s.%s", pkg, fnName)
}

func Caller(skip int) (file string, line int, fn string) {
	pc, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "", 0, "UNKNOWN"
	}

	return file, line, runtime.FuncForPC(pc).Name()
}

func CallerShortLocation(skip int) string {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "unknown"
	}

	return shortPath(file) + ":" + strconv.Itoa(line)
}

func CallerUniqueKey(skip int) string {
	pc, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "unknown"
	}

	return file + ":" + strconv.Itoa(line) + ":" + runtime.FuncForPC(pc).Name()
}

func shortPath(file string) string {
	parts := strings.Split(file, "/")
	if len(parts) <= 2 {
		return file
	}

	return strings.Join(parts[len(parts)-2:], "/")
}
