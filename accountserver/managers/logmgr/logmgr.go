package logmgr

import (
	"github.com/TarsCloud/TarsGo/tars"
)

var (
	RLOGGER = tars.GetLogger("accountserver")
	RERROR  = RLOGGER.Error
	RINFO   = RLOGGER.Info
	RWARN   = RLOGGER.Warn

	DLOGGER = tars.GetDayLogger("accountserver", 1)
)