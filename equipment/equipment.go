package equipment

import (
	"fmt"
	"io/ioutil"
	"os"
)

const (
	affixTableClass       = "LootRandomizerTable.tpl"
	affixTableTemplate    = "database\\Templates\\LootRandomizerTable.tpl"
	itemTableClass        = "LootItemTable_FixedWeight.tpl"
	itemTableTemplate     = "database\\Templates\\LootItemTable_FixedWeight.tpl"
	merchantTableClass    = "LootMasterTable.tpl"
	merchantTableTemplate = "database\\Templates\\LootMasterTable.tpl"
)

func allSlots() []string {
	return []string{"Amulet", "Arm", "Head", "Leg", "RingLeft", "RingRight", "Torso", "WeaponLeft", "WeaponRight"}
}

func itemTables() []string {
	return []string{"merchantTable", "itemTable", "itemPrefixTable", "itemSuffixTable"}
}

// Equipment is an entire equipment of a Titan Quest char plus all metadata for filesystem storage
type Equipment struct {
	Name  string
	Path  string
	Slots []Slot
}

// Slot represents one quipment slot
type Slot struct {
	Item Item
	Name string
	Path string
}

/*
type Slots struct {
	Head        Item
	Torso       Item
	Arm         Item
	Leg         Item
	RingLeft    Item
	RingRight   Item
	WeaponLeft  Item
	WeaponRight Item
	Amulet      Item
}
*/

// Item is all data any given item needs to be constructed
type Item struct {
	Base   string
	Prefix string
	Suffix string
	Record string
}

type table struct {
	Path    string
	Headers []byte
	Body    []byte
}

// New sets up a new equipment directory and initialises empty equipment tables
func New(name, folderPath string) (*Equipment, error) {
	e := Equipment{
		Name: name,
		Path: fmt.Sprintf("%s/%s", folderPath, name),
	}
	/*err := os.MkdirAll(fmt.Sprintf("%s/%s", e.Path, e.Name), 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %v", e.Path, err)
	}*/
	for _, s := range allSlots() {
		slot := Slot{
			Name: s,
			Item: Item{Base: "", Prefix: "", Suffix: "", Record: ""},
			Path: fmt.Sprintf("%s/%s", e.Path, s),
		}
		os.MkdirAll(fmt.Sprintf("%s", slot.Path), 0644)
		err := slot.init()
		if err != nil {
			return nil, fmt.Errorf("failed to initialise %s: %v", slot.Name, err)
		}

		e.Slots = append(e.Slots, slot)
	}
	// create item table
	// create prefix table
	// create suffix table
	return &e, nil
}

func (s Slot) init() error {
	prefixTable, err := s.createItemAffixTable(fmt.Sprintf("%s/%s.dbr", s.Path, "itemPrefixTable"), true)
	if err != nil {
		return fmt.Errorf("failed to initialise %s: %v", prefixTable.Path, err)
	}
	err = prefixTable.write()
	if err != nil {
		return fmt.Errorf("failed to write table to %s: %v", prefixTable.Path, err)
	}

	suffixTable, err := s.createItemAffixTable(fmt.Sprintf("%s/%s.dbr", s.Path, "itemSuffixTable"), false)
	if err != nil {
		return fmt.Errorf("failed to initialise %s: %v", suffixTable.Path, err)
	}
	err = suffixTable.write()
	if err != nil {
		return fmt.Errorf("failed to write table to %s: %v", suffixTable.Path, err)
	}

	itemTable, err := s.createItemTable(fmt.Sprintf("%s/%s.dbr", s.Path, "itemTable"), prefixTable, suffixTable)
	if err != nil {
		return fmt.Errorf("failed to initialise %s: %v", itemTable.Path, err)
	}
	err = itemTable.write()
	if err != nil {
		return fmt.Errorf("failed to write table to %s: %v", itemTable.Path, err)
	}

	merchantTable, err := s.createMerchantTable(fmt.Sprintf("%s/%s.dbr", s.Path, "merchantTable"), itemTable)
	if err != nil {
		return fmt.Errorf("failed to initialise %s: %v", merchantTable.Path, err)
	}
	err = merchantTable.write()
	if err != nil {
		return fmt.Errorf("failed to write table to %s: %v", merchantTable.Path, err)
	}
	return nil
}

func (s Slot) createItemAffixTable(path string, prefix bool) (*table, error) {
	var description string
	if prefix {
		description = s.Item.Prefix
	} else {
		description = s.Item.Suffix
	}
	headers, err := s.createTableHeader("itemAffixTable", description)
	if err != nil {
		return nil, fmt.Errorf("failed to create %s header: %v", path, err)
	}
	body := []byte("randomizerName1,,\nrandomizerWeight1,,\n")
	t := table{
		Path:    path,
		Headers: headers,
		Body:    body,
	}
	return &t, nil
}

func (s Slot) createItemTable(path string, prefixTable, suffixTable *table) (*table, error) {
	headers, err := s.createTableHeader("itemTable", fmt.Sprintf("%s %s %s", s.Item.Prefix, s.Item.Base, s.Item.Suffix))
	if err != nil {
		return nil, fmt.Errorf("failed to create %s header: %v", path, err)
	}
	baseConfig := []byte("bothPrefixSuffix,100\n")
	lootConfig := []byte("lootName1,,\nlootWeight1,100,\n")
	prefixConfig := []byte(fmt.Sprintf("prefixRandomizerChance,100,\nprefixRandomizerName1,%s,\nprefixRandomizerWeight1,,\n", prefixTable.Path))
	suffixConfig := []byte(fmt.Sprintf("suffixRandomizerChance,100,\nsuffixRandomizerName1,%s,\nsuffixRandomizerWeight1,,\n", suffixTable.Path))
	configs := [][]byte{baseConfig, lootConfig, prefixConfig, suffixConfig}
	var body []byte
	for _, c := range configs {
		body = append(body, c...)
	}
	t := table{
		Path:    path,
		Headers: headers,
		Body:    body,
	}
	return &t, nil
}

func (s Slot) createMerchantTable(path string, itemTable *table) (*table, error) {
	headers, err := s.createTableHeader("merchantTable", s.Item.Base)
	if err != nil {
		return nil, fmt.Errorf("failed to create %s header: %v", path, err)
	}
	body := []byte(fmt.Sprintf("lootName1,%s,\nlootWeight1,100,\n", itemTable.Path))
	t := table{
		Path:    path,
		Headers: headers,
		Body:    body,
	}
	return &t, nil
}

func (s Slot) createTableHeader(tableType, tableDescription string) ([]byte, error) {
	var template string
	var class string
	switch tableType {
	case "merchantTable":
		template = merchantTableTemplate
		class = merchantTableClass
	case "itemTable":
		template = itemTableTemplate
		class = itemTableClass
	case "itemAffixTable":
		template = affixTableTemplate
		class = affixTableClass
	default:
		return nil, fmt.Errorf("wrong tableType: %s", tableType)
	}
	return []byte(fmt.Sprintf("templateName,%s,\nActorName,,\nClass,%s,\nFileDescription,%s,\n", template, class, tableDescription)), nil
}

func (t table) write() error {
	err := ioutil.WriteFile(t.Path, append(t.Headers, t.Body...), 0644)
	if err != nil {
		return fmt.Errorf("failed writing table file %s: %v", t.Path, err)
	}
	return nil
}
