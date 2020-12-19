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
	fmt.Printf("%s %s %s", e.Slots[equipment.Head].Item.Prefix, e.Slots[equipment.Head].Item.Base, e.Slots[equipment.Head].Item.Suffix)
	if err := e.Slots[equipment.Head].SetItem(
		"Alcyoneus' Mask",
		"records\\xpack\\item\\equipmentarmor\\helm\\mi_n_gigantesmelee02.dbr",
		"Robust",
		"Records\\Item\\LootMagicalAffixes\\Prefix\\Default\\Rare_StrLife_02.dbr",
		"of Glory",
		"Records\\XPack\\Item\\LootMagicalAffixes\\Suffix\\Default\\Rare_%xp%elementalResist_02.dbr",
	); err != nil {
		log.Fatal(err)
	}
}
