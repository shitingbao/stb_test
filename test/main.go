package main

import (
	"log"

	"github.com/gogf/gf/v2/util/gconv"
)

type pen struct {
	Name string
}

type pen2 struct {
	Name []string
}

func main() {
	log.Println(gconv.Int(false))
	// TimestampToTime(1764926287)
}
