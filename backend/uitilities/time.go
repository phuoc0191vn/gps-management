package utilities

import (
	"net/url"
	"time"
)

func TimeInUTC(t time.Time) time.Time {
	return t.UTC()
}

func TimeInLocal(t time.Time) time.Time {
	return t.In(time.Local)
}

func GetQueryTime(query url.Values, key string) time.Time {
	values, ok := query[key]
	if !ok || len(values) < 1 {
		return time.Time{}
	}
	date := values[0]

	res, e := time.Parse(time.RFC3339, date)
	if e != nil {
		return time.Time{}
	}

	return res
}
