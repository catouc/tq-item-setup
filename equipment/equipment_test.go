package equipment

import (
	"fmt"
	"testing"

	"github.com/go-test/deep"
)

func TestSlotFromString(t *testing.T) {
	testData := []struct {
		Name string
		In   string
		Out  Slot
		OK   bool
	}{
		{
			Name: "ValidSlot",
			In:   "Amulet",
			Out:  Amulet,
			OK:   true,
		},
		{
			Name: "ValidSlot",
			In:   "Arm",
			Out:  Arm,
			OK:   true,
		},
		{
			Name: "ValidSlot",
			In:   "Head",
			Out:  Head,
			OK:   true,
		},
		{
			Name: "ValidSlot",
			In:   "Leg",
			Out:  Leg,
			OK:   true,
		},
		{
			Name: "ValidSlot",
			In:   "RingLeft",
			Out:  RingLeft,
			OK:   true,
		},
		{
			Name: "ValidSlot",
			In:   "RingRight",
			Out:  RingRight,
			OK:   true,
		},
		{
			Name: "ValidSlot",
			In:   "Torso",
			Out:  Torso,
			OK:   true,
		},
		{
			Name: "ValidSlot",
			In:   "WeaponLeft",
			Out:  WeaponLeft,
			OK:   true,
		},
		{
			Name: "ValidSlot",
			In:   "WeaponRight",
			Out:  WeaponRight,
			OK:   true,
		},
		{
			Name: "InvalidSlot",
			In:   "FooBar",
			Out:  UnknownSlot,
			OK:   false,
		},
	}
	for _, td := range testData {
		t.Run(td.Name, func(t *testing.T) {
			s, err := SlotFromString(td.In)
			if err != nil && td.OK {
				t.Errorf("unexpected error: %v", err)
			}
			if err == nil && !td.OK {
				t.Error("expected error but got nil")
			}
			if s != td.Out {
				t.Errorf("expected %s got %s instead", td.Out, s)
			}
		})
	}
}

func TestString(t *testing.T) {
	testData := []struct {
		Name string
		In   Slot
		Out  string
		OK   bool
	}{
		{
			Name: "ValidSlotAmulet",
			In:   Amulet,
			Out:  "Amulet",
			OK:   true,
		},
		{
			Name: "ValidSlotArm",
			In:   Arm,
			Out:  "Arm",
			OK:   true,
		},
		{
			Name: "ValidSlotHead",
			In:   Head,
			Out:  "Head",
			OK:   true,
		},
		{
			Name: "ValidSlotLeg",
			In:   Leg,
			Out:  "Leg",
			OK:   true,
		},
		{
			Name: "ValidSlotRingLeft",
			In:   RingLeft,
			Out:  "RingLeft",
			OK:   true,
		},
		{
			Name: "ValidSlotRingRight",
			In:   RingRight,
			Out:  "RingRight",
			OK:   true,
		},
		{
			Name: "ValidSlotTorso",
			In:   Torso,
			Out:  "Torso",
			OK:   true,
		},
		{
			Name: "ValidSlotWeaponLeft",
			In:   WeaponLeft,
			Out:  "WeaponLeft",
			OK:   true,
		},
		{
			Name: "ValidSlotWeaponRight",
			In:   WeaponRight,
			Out:  "WeaponRight",
			OK:   true,
		},
		{
			Name: "InvalidSlot",
			In:   UnknownSlot,
			Out:  "",
			OK:   true,
		},
	}
	for _, td := range testData {
		t.Run(td.Name, func(t *testing.T) {
			s := td.In.String()
			if s != td.Out {
				t.Errorf("expected %s got %s instead", td.Out, s)
			}
		})
	}
}

func TestFromFile(t *testing.T) {
	testData := []struct {
		Name string
		In   string
		Out  *Equipment
		OK   bool
	}{
		{
			Name: "ValidEquipment",
			In:   "../testData/validEquipment.yml",
			Out: &Equipment{
				Name: "TestEquipment",
				Path: "TestPath",
				Items: []Item{
					{
						SlotIdentifier: "Amulet",
						BaseName:       "TestBaseName",
						BaseRecord:     "Test/BaseRecord/record.dbr",
						PrefixName:     "TestPrefixName",
						PrefixRecord:   "Test/PrefixRecord/record.dbr",
						SuffixName:     "TestSuffixName",
						SuffixRecord:   "Test/SuffixRecord/record.dbr",
					},
				},
			},
			OK: true,
		},
		{
			Name: "ValidEquipmentMultipleItems",
			In:   "../testData/validEquipmentMultipleItems.yml",
			Out: &Equipment{
				Name: "TestEquipment",
				Path: "TestPath",
				Items: []Item{
					{
						SlotIdentifier: "Amulet",
						BaseName:       "TestBaseName",
						BaseRecord:     "Test/BaseRecord/record.dbr",
						PrefixName:     "TestPrefixName",
						PrefixRecord:   "Test/PrefixRecord/record.dbr",
						SuffixName:     "TestSuffixName",
						SuffixRecord:   "Test/SuffixRecord/record.dbr",
					},
					{
						SlotIdentifier: "Head",
						BaseName:       "TestBaseName",
						BaseRecord:     "Test/BaseRecord/record.dbr",
						PrefixName:     "TestPrefixName",
						PrefixRecord:   "Test/PrefixRecord/record.dbr",
						SuffixName:     "TestSuffixName",
						SuffixRecord:   "Test/SuffixRecord/record.dbr",
					},
				},
			},
			OK: true,
		},
		{
			Name: "ValidEquipmentNoItems",
			In:   "../testData/validEquipmentNoItems.yml",
			Out: &Equipment{
				Name:  "TestEquipment",
				Path:  "TestPath",
				Items: nil,
			},
			OK: true,
		},
		{
			Name: "InvalidEquipment",
			In:   "../testData/invalidEquipment.yml",
			Out:  nil,
			OK:   false,
		},
		{
			Name: "FileNotExists",
			In:   "../testData/IDoNotExist.notyml",
			Out:  nil,
			OK:   false,
		},
		{
			Name: "InvalidFile",
			In:   "../testData/invalidFile.yml",
			Out:  nil,
			OK:   false,
		},
	}
	for _, td := range testData {
		t.Run(td.Name, func(t *testing.T) {

			e, err := FromFile(td.In)
			if err != nil && td.OK {
				t.Errorf("unexpected error: %v", err)
			}
			if err == nil && !td.OK {
				t.Error("expected error but got nil")
			}
			if td.Out != nil {
				diff := deep.Equal(*e, *td.Out)
				if diff != nil {
					t.Errorf("result differs from expected file at filename %s: %+v", td.In, diff)
				}
			}
		})

	}
}

func TestNew(t *testing.T) {
	testData := []struct {
		Name         string
		InName       string
		InFolderPath string
		Out          *Equipment
		OK           bool
	}{
		{
			Name:         "ValidPath",
			InName:       "test_equip",
			InFolderPath: "test_out",
			Out: &Equipment{
				Name:  "test_equip",
				Path:  "test_out/test_equip",
				Items: nil,
			},
			OK: true,
		},
	}
	for _, td := range testData {
		t.Run(td.Name, func(t *testing.T) {
			e, err := New(td.InName, td.InFolderPath)
			if err != nil && td.OK {
				t.Errorf("unexpected error: %v", err)
			}
			if err == nil && !td.OK {
				t.Error("expected error but got nil")
			}
			diff := deep.Equal(*e, *td.Out)
			if diff != nil {
				t.Errorf("result differs from expected table: %+v", diff)
			}
		})
	}
}

func TestCreateAffixTable(t *testing.T) {
	testData := []struct {
		Name          string
		InPath        string
		InAffixName   string
		InAffixRecord string
		Out           *table
		OK            bool
	}{
		{
			Name:          "ValidTableUnixPath",
			InPath:        "test/Amulet/ItemPrefixTable.dbr",
			InAffixName:   "TestAffix",
			InAffixRecord: "records/item/LootMagicalAffixes/Prefix/Default/TestAffix.dbr",
			Out: &table{
				Path:    "test/Amulet/ItemPrefixTable.dbr",
				Headers: []byte("templateName,database\\Templates\\LootRandomizerTable.tpl,\nActorName,,\nClass,LootRandomizerTable.tpl,\nFileDescription,TestAffix,\n"),
				Body:    []byte("randomizerName1,records/item/LootMagicalAffixes/Prefix/Default/TestAffix.dbr,\nrandomizerWeight1,100,\n"),
			},
			OK: true,
		},
		{
			Name:          "ValidTableWindowsPath",
			InPath:        "test\\Amulet\\ItemPrefixTable.dbr",
			InAffixName:   "TestAffix",
			InAffixRecord: "records/item/LootMagicalAffixes/Prefix/Default/TestAffix.dbr",
			Out: &table{
				Path:    "test\\Amulet\\ItemPrefixTable.dbr",
				Headers: []byte("templateName,database\\Templates\\LootRandomizerTable.tpl,\nActorName,,\nClass,LootRandomizerTable.tpl,\nFileDescription,TestAffix,\n"),
				Body:    []byte("randomizerName1,records/item/LootMagicalAffixes/Prefix/Default/TestAffix.dbr,\nrandomizerWeight1,100,\n"),
			},
			OK: true,
		},
	}
	for _, td := range testData {
		t.Run(td.Name, func(t *testing.T) {
			table, err := createItemAffixTable(td.InPath, td.InAffixName, td.InAffixRecord)
			if err != nil && td.OK {
				t.Errorf("unexpected error: %v", err)
			}
			if err == nil && !td.OK {
				t.Error("expected error but got nil")
			}
			diff := deep.Equal(*table, *td.Out)
			if diff != nil {
				t.Errorf("result differs from expected table: %+v", diff)
			}
		})
	}
}

func TestCreateItemTable(t *testing.T) {
	testData := []struct {
		Name          string
		InPath        string
		InLootPath    string
		InPrefixPath  string
		InSuffixPath  string
		InDescription string
		Out           *table
		OK            bool
	}{
		{
			Name:          "ValidTableUnixPath",
			InPath:        "test/Amulet/ItemTable.dbr",
			InLootPath:    "records/item/equipmenthelm/helm.dbr",
			InPrefixPath:  "test/Amulet/ItemPrefixTable.dbr",
			InSuffixPath:  "test/Amulet/ItemSuffixTable.dbr",
			InDescription: "TestItemTable",
			Out: &table{
				Path:    "test/Amulet/ItemTable.dbr",
				Headers: []byte("templateName,database\\Templates\\LootItemTable_FixedWeight.tpl,\nActorName,,\nClass,LootItemTable_FixedWeight.tpl,\nFileDescription,TestItemTable,\n"),
				Body: []byte("bothPrefixSuffix,100\n" +
					"lootName1,records/item/equipmenthelm/helm.dbr,\nlootWeight1,100,\n" +
					"prefixRandomizerChance,100,\nprefixRandomizerName1,test/Amulet/ItemPrefixTable.dbr,\nprefixRandomizerWeight1,,\n" +
					"suffixRandomizerChance,100,\nsuffixRandomizerName1,test/Amulet/ItemSuffixTable.dbr,\nsuffixRandomizerWeight1,,\n"),
			},
			OK: true,
		},
		{
			Name:          "ValidTableWindowsPath",
			InPath:        "test\\Amulet\\ItemTable.dbr",
			InLootPath:    "records\\item\\equipmenthelm\\helm.dbr",
			InPrefixPath:  "test\\Amulet\\ItemPrefixTable.dbr",
			InSuffixPath:  "test\\Amulet\\ItemSuffixTable.dbr",
			InDescription: "TestItemTable",
			Out: &table{
				Path:    "test\\Amulet\\ItemTable.dbr",
				Headers: []byte("templateName,database\\Templates\\LootItemTable_FixedWeight.tpl,\nActorName,,\nClass,LootItemTable_FixedWeight.tpl,\nFileDescription,TestItemTable,\n"),
				Body: []byte("bothPrefixSuffix,100\n" +
					"lootName1,records\\item\\equipmenthelm\\helm.dbr,\nlootWeight1,100,\n" +
					"prefixRandomizerChance,100,\nprefixRandomizerName1,test\\Amulet\\ItemPrefixTable.dbr,\nprefixRandomizerWeight1,,\n" +
					"suffixRandomizerChance,100,\nsuffixRandomizerName1,test\\Amulet\\ItemSuffixTable.dbr,\nsuffixRandomizerWeight1,,\n"),
			},
			OK: true,
		},
	}
	for _, td := range testData {
		t.Run(td.Name, func(t *testing.T) {
			table, err := createItemTable(td.InPath, td.InLootPath, td.InPrefixPath, td.InSuffixPath, td.InDescription)
			if err != nil && td.OK {
				t.Errorf("unexpected error: %v", err)
			}
			if err == nil && !td.OK {
				t.Error("expected error but got nil")
			}
			diff := deep.Equal(*table, *td.Out)
			if diff != nil {
				t.Errorf("result differs from expected table: %+v", diff)
			}
		})
	}
}

func TestCreateMerchantTable(t *testing.T) {
	testData := []struct {
		Name          string
		InPath        string
		InItemPath    string
		InDescription string
		Out           *table
		OK            bool
	}{
		{
			Name:          "ValidTableUnixPath",
			InPath:        "test/Amulet/MerchantTable.dbr",
			InItemPath:    "test/Amulet/ItemTable.dbr",
			InDescription: "TestMerchantTable",
			Out: &table{
				Path:    "test/Amulet/MerchantTable.dbr",
				Headers: []byte("templateName,database\\Templates\\LootMasterTable.tpl,\nActorName,,\nClass,LootMasterTable.tpl,\nFileDescription,TestMerchantTable,\n"),
				Body:    []byte("lootName1,test/Amulet/ItemTable.dbr,\nlootWeight1,100,\n"),
			},
			OK: true,
		},
	}
	for _, td := range testData {
		t.Run(td.Name, func(t *testing.T) {
			table, err := createMerchantTable(td.InPath, td.InItemPath, td.InDescription)
			if err != nil && td.OK {
				t.Errorf("unexpected error: %v", err)
			}
			if err == nil && !td.OK {
				t.Error("expected error but got nil")
			}
			diff := deep.Equal(*table, *td.Out)
			if diff != nil {
				fmt.Println(*table)
				t.Errorf("result differs from expected table: %+v", diff)
			}
		})
	}
}
