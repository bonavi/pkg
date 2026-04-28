package stackTrace

import (
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"sync/atomic"
)

type stackTracer struct {
	serviceName  string
	isEnabled    atomic.Bool
	samplingRate int
}

var stackTracerInstance = &stackTracer{
	serviceName:  "",
	isEnabled:    atomic.Bool{},
	samplingRate: 1,
}

func SetIsEnabled(enabled bool) {
	stackTracerInstance.isEnabled.Store(enabled)
}

func Init(serviceName string, isEnabled bool, samplingRate int) {
	stackTracerInstance = &stackTracer{
		serviceName: serviceName,
		isEnabled:   atomic.Bool{},
		samplingRate: samplingRate,
	}

	stackTracerInstance.isEnabled.Store(isEnabled)
}

func GetStackTrace(skip int) []string {

	if !stackTracerInstance.isEnabled.Load() || stackTracerInstance.samplingRate < 1 {
		return nil
	}

	if rand.Intn(stackTracerInstance.samplingRate) != 0 {
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
