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

// Equipment is an entire equipment of a Titan Quest char plus all metadata for filesystem storage
type Equipment struct {
	Name  string
	Path  string
	Slots map[string]slot
}

// slot represents one quipment slot
type slot struct {
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
		Path:  fmt.Sprintf("%s/%s", folderPath, name),
		Slots: make(map[string]slot),
	}
	for _, sName := range allSlots() {
		s := slot{
			Name: sName,
			Item: item{Base: "", Prefix: "", Suffix: "", Record: ""},
			Path: fmt.Sprintf("%s/%s", e.Path, sName),
		}
		os.MkdirAll(s.Path, 0644)
		err := s.init()
		if err != nil {
			return nil, fmt.Errorf("failed to initialise %s: %v", s.Name, err)
		}

		e.Slots[sName] = s
	}
	return &e, nil
}

func (s slot) init() error {
	err := s.createItem("", "", "", "", "", "")
	if err != nil {
		return fmt.Errorf("failed to create item: %v", err)
	}
	return nil
}

func (s slot) createItemAffixTable(path, affixName, affixRecord string) (*table, error) {
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

func (s slot) createItemTable(path string, prefixTable, suffixTable *table) (*table, error) {
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

func (s slot) createMerchantTable(path string, itemTable *table) (*table, error) {
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

func (s slot) createTableHeader(tableType, tableDescription string) ([]byte, error) {
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

func (s slot) createItem(baseName, baseRecord, prefixName, prefixRecord, suffixName, suffixRecord string) error {
	prefixTable, err := s.createItemAffixTable(fmt.Sprintf("%s/%s.dbr", s.Path, "itemPrefixTable"), prefixName, prefixRecord)
	if err != nil {
		return fmt.Errorf("failed to initialise %s: %v", prefixTable.Path, err)
	}
	err = prefixTable.write()
	if err != nil {
		return fmt.Errorf("failed to write table to %s: %v", prefixTable.Path, err)
	}
	s.Item.PrefixTable = prefixTable

	suffixTable, err := s.createItemAffixTable(fmt.Sprintf("%s/%s.dbr", s.Path, "itemSuffixTable"), suffixName, suffixRecord)
	if err != nil {
		return fmt.Errorf("failed to initialise %s: %v", suffixTable.Path, err)
	}
	err = suffixTable.write()
	if err != nil {
		return fmt.Errorf("failed to write table to %s: %v", suffixTable.Path, err)
	}
	s.Item.SuffixTable = suffixTable

	itemTable, err := s.createItemTable(fmt.Sprintf("%s/%s.dbr", s.Path, "itemTable"), prefixTable, suffixTable)
	if err != nil {
		return fmt.Errorf("failed to initialise %s: %v", itemTable.Path, err)
	}
	err = itemTable.write()
	if err != nil {
		return fmt.Errorf("failed to write table to %s: %v", itemTable.Path, err)
	}
	s.Item.LootTable = itemTable

	merchantTable, err := s.createMerchantTable(fmt.Sprintf("%s/%s.dbr", s.Path, "merchantTable"), itemTable)
	if err != nil {
		return fmt.Errorf("failed to initialise %s: %v", merchantTable.Path, err)
	}
	err = merchantTable.write()
	if err != nil {
		return fmt.Errorf("failed to write table to %s: %v", merchantTable.Path, err)
	}
	s.MerchantTable = merchantTable
	return nil
}

func (s slot) setItem(baseName, baseRecord, prefixName, prefixRecord, suffixName, suffixRecord string) error {
	s.Item.Base = baseName
	s.Item.Prefix = prefixName
	s.Item.Suffix = suffixName

	s.createItem(baseName, baseRecord, prefixName, prefixRecord, suffixName, suffixRecord)
	return nil
}

func (t table) write() error {
	err := ioutil.WriteFile(t.Path, append(t.Headers, t.Body...), 0644)
	if err != nil {
		return fmt.Errorf("failed writing table file %s: %v", t.Path, err)
	}
	return nil
}

func (t table) update(path string, header, body []byte) error {
	return nil
}
