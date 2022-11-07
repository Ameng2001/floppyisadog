package logmgr

import (
	"github.com/TarsCloud/TarsGo/tars"
)

var (
	RLOGGER = tars.GetLogger("migratetool")
	RERROR  = RLOGGER.Error
	RINFO   = RLOGGER.Info
	RWARN   = RLOGGER.Warn

	DLOGGER = tars.GetDayLogger("migratetool", 1)
)
