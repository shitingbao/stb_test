package main

import (
	"log"
)

func main() {
	defer func() {
		log.Println("11111")
	}()

	defer func() {
		log.Println("22222")
	}()

	defer func() {
		log.Println("33333")
	}()
	log.Println("base:")
}
