package main

import (
	"fmt"
	"log"

	"github.com/Deichindianer/tq-item-setup/equipment"
)

func main() {
	e, err := equipment.New("lvl_30", "out")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s %s %s", e.Slots["Head"].Item.Prefix, e.Slots["Head"].Item.Base, e.Slots["Head"].Item.Suffix)
}
