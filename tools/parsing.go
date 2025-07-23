package tools

import (
	"github.com/erikpa1/TurtleIntelligenceBackend/lg"
	"strconv"
)

func StringToInt32(value string, defval int32) int32 {
	status, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		lg.LogE(err)
		return defval
	} else {
		return int32(status)
	}

}

func StringToInt64(value string, defval int64) int64 {
	status, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		lg.LogE(err)
		return defval
	} else {
		return status
	}

}
