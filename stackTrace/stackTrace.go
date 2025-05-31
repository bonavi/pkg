package stackTrace

import (
	"fmt"
	"runtime"
	"strings"
	"sync/atomic"
)

type stackTracer struct {
	serviceName string
	isEnabled   atomic.Bool
}

var stackTracerInstance = &stackTracer{
	serviceName: "",
}

func SetIsEnabled(enabled bool) {
	stackTracerInstance.isEnabled.Store(enabled)
}

func Init(serviceName string, isEnabled bool) {
	stackTracerInstance = &stackTracer{
		serviceName: serviceName,
		isEnabled:   atomic.Bool{},
	}

	stackTracerInstance.isEnabled.Store(isEnabled)
}

func GetStackTrace(skip int) []string {

	if !stackTracerInstance.isEnabled.Load() {
		return nil
	}

	var pcs [32]uintptr
	n := runtime.Callers(0, pcs[:])
	var path []string
	for i := skip; i < n; i++ {
		_, file, line, _ := runtime.Caller(i)
		if strings.Contains(file, stackTracerInstance.serviceName) {
			path = append(path, fmt.Sprintf("%s:%d", file, line))
		}
	}
	return path
}
