package main

import (
	"log"

	"github.com/gogf/gf/util/gconv"
)

type pen struct {
	Name string
}

type pen2 struct {
	Name []string
}

func main() {
	p := pen{Name: "aa,bb"}

	p2 := pen2{}

	err := gconv.Struct(p, &p2)
	log.Println(err, gconv.String(p2))
}
