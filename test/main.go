package main

import (
	"log"

	"github.com/google/uuid"
	u2 "github.com/pborman/uuid"
)

func main() {

	log.Println("base:", uuid.New().String())
	log.Println("base2:", u2.New())
}
