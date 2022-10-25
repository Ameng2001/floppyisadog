package log

import (
	"github.com/TarsCloud/TarsGo/tars"
)

var (
	RLOGGER = tars.GetLogger("foauthserver")
	RERROR 	= RLOGGER.Error
	RINFO	= RLOGGER.Info
	RWARN	= RLOGGER.Warn


	DLOGGER = tars.GetDayLogger("foauthserver", 1)
)
