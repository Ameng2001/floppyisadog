package logmgr

import (
	"github.com/TarsCloud/TarsGo/tars"
)

var (
	RLOGGER = tars.GetLogger("webportalserver")
	RERROR  = RLOGGER.Errorf
	RINFO   = RLOGGER.Infof
	RWARN   = RLOGGER.Warnf

	DLOGGER = tars.GetDayLogger("webportalserver", 1)
)
