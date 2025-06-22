package log

import (
	"os"
	"pkg/errors"

	"pkg/log/model"
)

func InitDefaultLogger(
	systemInfo model.SystemInfo,
	loggerCfg LoggerSettingsEnv,
) error {

	var logHandlers []Handler
	switch loggerCfg.LogFormat {
	case TextFormat:
		logHandlers = append(logHandlers, NewTextHandler(os.Stdout, loggerCfg.LogLevel))
	case JSONFormat:
		logHandlers = append(logHandlers, NewJSONHandler(os.Stdout, loggerCfg.LogLevel))
	}

	if err := Init(systemInfo, logHandlers...); err != nil {
		return errors.Default.Wrap(err).
			SkipThisCall()
	}

	return nil
}
