package testUtils

import (
	"time"

	"github.com/agiledragon/gomonkey/v2"
)

var defaultTime = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

func MockTime() *gomonkey.Patches {
	return gomonkey.ApplyFuncReturn(time.Now, defaultTime)
}
