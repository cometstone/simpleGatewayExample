package common

import (
	"strconv"
	"time"
)

func RequestID() string {
	base :=strconv.FormatInt(time.Now().UnixNano(),10)
	return strconv.Itoa(Conf.Api.ServerID) + base
}
