package main

import (
	"flag"
	"log"

	"github.com/Deichindianer/tq-item-setup/equipment"
)

func main() {
	var pathFlag = flag.String("path", "str_lvl_45.yml", "set the path to the equipment yml, defaults to equipment.yml")
	flag.Parse()
	e, err := equipment.FromFile(*pathFlag)
	if err != nil {
		log.Fatal(err)
	}
	if err := e.Flush(); err != nil {
		log.Fatal(err)
	}
}
