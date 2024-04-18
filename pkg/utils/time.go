package utils

import (
	"strconv"
	"time"
)

func ConvertTimestampStringToTime(timestampString string) time.Time {
	i, err := strconv.ParseFloat(timestampString, 64)
	if err != nil {
		panic(err)
	}
	return time.Unix(int64(i), 0)
}
