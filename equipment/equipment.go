package equipment

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	affixTableClass       = "LootRandomizerTable.tpl"
	affixTableTemplate    = "database\\Templates\\LootRandomizerTable.tpl"
	itemTableClass        = "LootItemTable_FixedWeight.tpl"
	itemTableTemplate     = "database\\Templates\\LootItemTable_FixedWeight.tpl"
	merchantTableClass    = "LootMasterTable.tpl"
	merchantTableTemplate = "database\\Templates\\LootMasterTable.tpl"

	// I know this could be iota but this is way more readable :)

	Amulet      = 0
	Arm         = 1
	Head        = 2
	Leg         = 3
	RingLeft    = 4
	RingRight   = 5
	Torso       = 6
	WeaponLeft  = 7
	WeaponRight = 8
)

//var allSlots = []string{"Amulet", "Arm", "Head", "Leg", "RingLeft", "RingRight", "Torso", "WeaponLeft", "WeaponRight"}
var allSlots = []Slot{Amulet, Arm, Head, Leg, RingLeft, RingRight, Torso, WeaponLeft, WeaponRight}

// Equipment is an entire equipment of a Titan Quest char plus all metadata for filesystem storage
type Equipment struct {
	Name  string
	Path  string
	Slots map[Slot]equipmentSlot
}

// Slot is the slot where the equipment goes, lol
type Slot int

func (i Slot) String() string {
	switch i {
	case Amulet:
		return "Amulet"
	case Arm:
		return "Arm"
	case Leg:
		return "Leg"
	case RingLeft:
		return "RingLeft"
	case RingRight:
		return "RingRight"
	case Torso:
		return "Torso"
	case WeaponLeft:
		return "WeaponLeft"
	case WeaponRight:
		return "WeaponRight"
	}
	return ""
}

// equipmentSlot represents one quipment equipmentSlot
type equipmentSlot struct {
	Item          item
	Name          string
	Path          string
	MerchantTable *table
}

// item is all data any given item needs to be constructed
type item struct {
	Base        string
	LootTable   *table
	Prefix      string
	PrefixTable *table
	SuffixTable *table
	Suffix      string
	Record      string
}

type table struct {
	Path    string
	Headers []byte
	Body    []byte
}

// New sets up a new equipment directory and initialises empty equipment tables
func New(name, folderPath string) (*Equipment, error) {
	e := Equipment{
		Name:  name,
		Path:  filepath.Join(folderPath, name),
		Slots: make(map[Slot]equipmentSlot),
	}
	for _, is := range allSlots {
		s := equipmentSlot{
			Name: is.String(),
			Item: item{Base: "", Prefix: "", Suffix: "", Record: ""},
			Path: filepath.Join(e.Path, is.String()),
		}
		if err := os.MkdirAll(s.Path, 0644); err != nil {
			return nil, fmt.Errorf("failed to create %s: %v", s.Path, err)
		}
		if err := s.init(); err != nil {
			return nil, fmt.Errorf("failed to initialise %s: %v", s.Name, err)
		}
		e.Slots[is] = s
	}
	return &e, nil
}

func (s equipmentSlot) init() error {
	err := s.createItem("", "", "", "", "", "")
	if err != nil {
		return fmt.Errorf("failed to create item: %v", err)
	}
	return nil
}

func (s equipmentSlot) createItemAffixTable(path, affixName, affixRecord string) (*table, error) {
	headers, err := s.createTableHeader("itemAffixTable", affixName)
	if err != nil {
		return nil, fmt.Errorf("failed to create %s header: %v", path, err)
	}
	body := []byte(fmt.Sprintf("randomizerName1,%s,\nrandomizerWeight1,100,\n", affixRecord))
	t := table{
		Path:    path,
		Headers: headers,
		Body:    body,
	}
	return &t, nil
}

func (s equipmentSlot) createItemTable(path string, prefixTable, suffixTable *table) (*table, error) {
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

func (s equipmentSlot) createMerchantTable(path string, itemTable *table) (*table, error) {
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

func (s equipmentSlot) createTableHeader(tableType, tableDescription string) ([]byte, error) {
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

func (s equipmentSlot) createItem(baseName, baseRecord, prefixName, prefixRecord, suffixName, suffixRecord string) error {
	prefixTable, err := s.createItemAffixTable(filepath.Join(s.Path, "itemPrefixTable.dbr"), prefixName, prefixRecord)
	if err != nil {
		return fmt.Errorf("failed to initialise %s: %v", prefixTable.Path, err)
	}
	if err := prefixTable.write(); err != nil {
		return fmt.Errorf("failed to write table to %s: %v", prefixTable.Path, err)
	}
	s.Item.PrefixTable = prefixTable

	suffixTable, err := s.createItemAffixTable(filepath.Join(s.Path, "itemSuffixTable.dbr"), suffixName, suffixRecord)
	if err != nil {
		return fmt.Errorf("failed to initialise %s: %v", suffixTable.Path, err)
	}
	if err := suffixTable.write(); err != nil {
		return fmt.Errorf("failed to write table to %s: %v", suffixTable.Path, err)
	}
	s.Item.SuffixTable = suffixTable

	itemTable, err := s.createItemTable(filepath.Join(s.Path, "itemTable.dbr"), prefixTable, suffixTable)
	if err != nil {
		return fmt.Errorf("failed to initialise %s: %v", itemTable.Path, err)
	}
	if err := itemTable.write(); err != nil {
		return fmt.Errorf("failed to write table to %s: %v", itemTable.Path, err)
	}
	s.Item.LootTable = itemTable

	merchantTable, err := s.createMerchantTable(filepath.Join(s.Path, "merchantTable.dbr"), itemTable)
	if err != nil {
		return fmt.Errorf("failed to initialise %s: %v", merchantTable.Path, err)
	}
	if err := merchantTable.write(); err != nil {
		return fmt.Errorf("failed to write table to %s: %v", merchantTable.Path, err)
	}
	s.MerchantTable = merchantTable
	return nil
}

func (s equipmentSlot) SetItem(baseName, baseRecord, prefixName, prefixRecord, suffixName, suffixRecord string) error {
	s.Item.Base = baseName
	s.Item.Prefix = prefixName
	s.Item.Suffix = suffixName

	if err := s.createItem(baseName, baseRecord, prefixName, prefixRecord, suffixName, suffixRecord); err != nil {
		return fmt.Errorf("failed to set item: %v", err)
	}
	return nil
}

func (t table) write() error {
	err := ioutil.WriteFile(t.Path, append(t.Headers, t.Body...), 0644)
	if err != nil {
		return fmt.Errorf("failed writing table file %s: %v", t.Path, err)
	}
	return nil
}
