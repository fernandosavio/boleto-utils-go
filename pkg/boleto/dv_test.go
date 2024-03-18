package boleto

import "testing"

func TestMod11(t *testing.T) {
	testValues := []struct {
		input    string
		expected string
		skip     int
	}{
		{"1119_444455555555556666666666666666666666666", "1", 4},
		{"1049_898100000214032006561000100040099726390", "9", 4},
		{"7569_903800002500001434301033723400014933001", "6", 4},
		{"0019_667900002434790000002656973019362470618", "1", 4},
		{"0019_586200000773520000002464206011816073018", "5", 4},
		{"7559_896700003787000003389850761252543475984", "2", 4},
		{"2379_672000003097902028060007024617500249000", "1", 4},
		{"2379_672000003249052028269705944177105205220", ".", 4},
		{"1119_100255555555556666666666666666666666666", ".", 4},
		{"8220000215048200974123220154098290108605940", ".", -1},
		{"01230067896", ".", -1},
		{"31230067896", ".", -1},
		{"01231068896", "9", -1},
		{"01230267896", "8", -1},
		{"01231167896", "7", -1},
		{"01232067896", "6", -1},
		{"01241067896", "5", -1},
		{"01250067896", "4", -1},
		{"02232067896", "3", -1},
		{"01250067897", "2", -1},
		{"01230367896", "1", -1},
	}

	for _, tt := range testValues {
		got := mod11([]byte(tt.input), bdot, tt.skip)

		if result := string(got); result != tt.expected {
			t.Fatalf("%q: expected %q got %q", tt.input, tt.expected, result)
		}
	}
}

func TestMod10(t *testing.T) {
	testValues := []struct {
		input    string
		expected string
	}{
		{"01230067896", "3"},
		{"01230167896", "2"},
		{"01230267896", "1"},
		{"01230367896", "0"},
		{"01230467896", "9"},
		{"01230567896", "8"},
		{"01230667896", "7"},
		{"01230767896", "6"},
		{"01230867896", "5"},
		{"01230967896", "4"},
	}

	for _, tt := range testValues {
		got := mod10([]byte(tt.input))

		if result := string(got); result != tt.expected {
			t.Fatalf("%q: expected %q got %q", tt.input, tt.expected, result)
		}
	}
}
