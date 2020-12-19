package main

import (
	"log"

	"github.com/Deichindianer/tq-item-setup/equipment"
)

func main() {
	/*
		templateName,database\Templates\LootItemTable_FixedWeight.tpl,
		ActorName,,
		Class,LootItemTable_FixedWeight,
		FileDescription,,
	*/
	_, err := equipment.New("lvl_30", "out")
	if err != nil {
		log.Fatal(err)
	}
}
