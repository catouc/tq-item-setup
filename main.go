package main

import (
	"log"

	"github.com/Deichindianer/tq-item-setup/equipment"
)

func main() {
	e, err := equipment.FromFile("equipment.yml")
	if err != nil {
		log.Fatal(err)
	}
	if err := e.Flush(); err != nil {
		log.Fatal(err)
	}
}
