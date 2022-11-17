package logmgr

import (
	"github.com/TarsCloud/TarsGo/tars"
)

var (
	RLOGGER = tars.GetLogger("webportalserver")
	RERROR  = RLOGGER.Error
	RINFO   = RLOGGER.Info
	RWARN   = RLOGGER.Warn

	DLOGGER = tars.GetDayLogger("webportalserver", 1)
)
