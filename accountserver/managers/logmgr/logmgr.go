package logmgr

import (
	"github.com/TarsCloud/TarsGo/tars"
)

var (
	RLOGGER = tars.GetLogger("accountserver")
	RERROR  = RLOGGER.Errorf
	RINFO   = RLOGGER.Infof
	RWARN   = RLOGGER.Warnf

	DLOGGER = tars.GetDayLogger("accountserver", 1)
	DERRORF = DLOGGER.Errorf
	DINFOF  = DLOGGER.Infof
	DWARNF  = DLOGGER.Warnf
)
