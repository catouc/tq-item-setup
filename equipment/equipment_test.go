package equipment

import (
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
