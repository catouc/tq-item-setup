package equipment

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-yaml/yaml"
)

const (
	affixTableClass       = "LootRandomizerTable.tpl"
	affixTableTemplate    = "database\\Templates\\LootRandomizerTable.tpl"
	itemTableClass        = "LootItemTable_FixedWeight.tpl"
	itemTableTemplate     = "database\\Templates\\LootItemTable_FixedWeight.tpl"
	merchantTableClass    = "LootMasterTable.tpl"
	merchantTableTemplate = "database\\Templates\\LootMasterTable.tpl"

	// I know this could be iota but this is way more readable :)

	UnknownSlot = -1
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

// Equipment is an entire equipment of a Titan Quest char plus all metadata for filesystem storage
type Equipment struct {
	Name  string `yaml:"Name"`
	Path  string `yaml:"Path"`
	Items []Item `yaml:"Items"`
}

type Item struct {
	SlotIdentifier string `yaml:"SlotIdentifier"`
	BaseName       string `yaml:"BaseName"`
	BaseRecord     string `yaml:"BaseRecord"`
	PrefixName     string `yaml:"PrefixName"`
	PrefixRecord   string `yaml:"PrefixRecord"`
	SuffixName     string `yaml:"SuffixName"`
	SuffixRecord   string `yaml:"SuffixRecord"`
}

// Slot is the slot where the equipment goes, lol
type Slot int

func (i Slot) String() string {
	switch i {
	case Amulet:
		return "Amulet"
	case Arm:
		return "Arm"
	case Head:
		return "Head"
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

func SlotFromString(s string) (Slot, error) {
	switch s {
	case "Amulet":
		return Amulet, nil
	case "Arm":
		return Arm, nil
	case "Head":
		return Head, nil
	case "Leg":
		return Leg, nil
	case "RingLeft":
		return RingLeft, nil
	case "RingRight":
		return RingRight, nil
	case "Torso":
		return Torso, nil
	case "WeaponLeft":
		return WeaponLeft, nil
	case "WeaponRight":
		return WeaponRight, nil
	}
	return UnknownSlot, fmt.Errorf("unexpectected Slot %s", s)
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
		Path: filepath.Join(folderPath, name),
	}
	return &e, nil
}

func (e *Equipment) Flush() error {
	if err := os.MkdirAll(e.Path, 0644); err != nil {
		return fmt.Errorf("failed to create %s: %v", e.Path, err)
	}
	if err := e.createItems(); err != nil {
		return err
	}
	return nil
}

func (i *Item) Validate() error {
	_, err := SlotFromString(i.SlotIdentifier)
	if err != nil {
		return err
	}
	return nil
}

func FromFile(path string) (*Equipment, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open equipment file: %v", err)
	}
	var e Equipment
	if err := yaml.Unmarshal(f, &e); err != nil {
		return nil, fmt.Errorf("failed to unmarshal equip file: %v", err)
	}
	for _, i := range e.Items {
		if err := i.Validate(); err != nil {
			return nil, fmt.Errorf("item %s is not valid: %v", i.BaseName, err)
		}
	}
	return &e, nil
}

func createItemAffixTable(path, affixName, affixRecord, description string) (*table, error) {
	headers, err := createTableHeader("itemAffixTable", affixName)
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

func createItemTable(path string, prefixTable, suffixTable *table, description string) (*table, error) {
	// TODO: find a way to make descriptions pretty, maybe this needs to be implemented on top of the item after all just for that?
	headers, err := createTableHeader("itemTable", "")
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

func createMerchantTable(path string, itemTable *table, description string) (*table, error) {
	headers, err := createTableHeader("merchantTable", "")
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

func createTableHeader(tableType, tableDescription string) ([]byte, error) {
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

func (e *Equipment) createItems() error {
	for _, item := range e.Items {
		if err := e.createItem(item); err != nil {
			return err
		}
	}
	return nil
}

func (e *Equipment) createItem(item Item) error {
	if err := item.Validate(); err != nil {
		return fmt.Errorf("item is invalid: %v", err)
	}
	basePath := filepath.Join(e.Path, item.SlotIdentifier)
	if err := os.MkdirAll(basePath, 0644); err != nil {
		return fmt.Errorf("failed to create %s: %v", basePath, err)
	}
	prefixTable, err := createItemAffixTable(filepath.Join(basePath, "itemPrefixTable.dbr"), item.PrefixName, item.PrefixRecord, item.PrefixName)
	if err != nil {
		return fmt.Errorf("failed to initialise %s: %v", prefixTable.Path, err)
	}
	if err := prefixTable.write(); err != nil {
		return fmt.Errorf("failed to write table to %s: %v", prefixTable.Path, err)
	}

	suffixTable, err := createItemAffixTable(filepath.Join(basePath, "itemSuffixTable.dbr"), item.SuffixName, item.SuffixRecord, item.SuffixName)
	if err != nil {
		return fmt.Errorf("failed to initialise %s: %v", suffixTable.Path, err)
	}
	if err := suffixTable.write(); err != nil {
		return fmt.Errorf("failed to write table to %s: %v", suffixTable.Path, err)
	}

	itemTable, err := createItemTable(
		filepath.Join(basePath, "itemTable.dbr"),
		prefixTable,
		suffixTable,
		fmt.Sprintf("%s %s %s", item.PrefixName, item.BaseName, item.SuffixName),
	)
	if err != nil {
		return fmt.Errorf("failed to initialise %s: %v", itemTable.Path, err)
	}
	if err := itemTable.write(); err != nil {
		return fmt.Errorf("failed to write table to %s: %v", itemTable.Path, err)
	}

	merchantTable, err := createMerchantTable(filepath.Join(basePath, "merchantTable.dbr"), itemTable, item.BaseName)
	if err != nil {
		return fmt.Errorf("failed to initialise %s: %v", merchantTable.Path, err)
	}
	if err := merchantTable.write(); err != nil {
		return fmt.Errorf("failed to write table to %s: %v", merchantTable.Path, err)
	}
	return nil
}

func (t *table) write() error {
	err := ioutil.WriteFile(t.Path, append(t.Headers, t.Body...), 0644)
	if err != nil {
		return fmt.Errorf("failed writing table file %s: %v", t.Path, err)
	}
	return nil
}
