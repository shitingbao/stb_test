package main

import (
	"log"
	"time"
)

// 秒级时间戳 → 格式化时间（YYYY-MM-DD HH:MM:SS）
func TimestampToTime(ts int64) string {
	t := time.Unix(ts, 0)
	log.Println("TimestampToTime:", t.Format("2006-01-02 15:04:05"))
	return t.Format("2006-01-02 15:04:05")
}

// 毫秒时间戳 → 格式化时间（YYYY-MM-DD HH:MM:SS）
func MilliTimestampToTime(ms int64) string {
	t := time.Unix(0, ms*int64(time.Millisecond))
	log.Println("MilliTimestampToTime:", t.Format("2006-01-02 15:04:05"))
	return t.Format("2006-01-02 15:04:05")
}

// 格式化时间（YYYY-MM-DD HH:MM:SS）→ 时间戳（秒 / 毫秒）
func TimeToTimestamp(tstr string) (sec int64, ms int64, err error) {
	layout := "2006-01-02 15:04:05"

	t, err := time.ParseInLocation(layout, tstr, time.Local)
	if err != nil {
		return 0, 0, err
	}

	sec = t.Unix()
	ms = t.UnixNano() / int64(time.Millisecond)
	log.Println("sec,ms:", sec, ms)
	return sec, ms, nil
}
