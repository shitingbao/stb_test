package main

import (
	"log"
)

type pen struct {
	Name string
}

type pen2 struct {
	Name []string
}

func main() {
	n, err := loadHappy(1)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("done, fetched=%d\n", n)
}
