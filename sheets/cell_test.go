package sheets

import "testing"

func TestGetCellCoords(t *testing.T) {
	var tests = []struct {
		name   string
		input  string
		row    int32
		column int32
		err    bool
	}{
		{"A12A", "A12A", 0, 0, true},
		{"323", "323", 0, 0, true},
		{"BBB", "BBB", 0, 0, true},
		{"+-21", "+-21", 0, 0, true},
		{"A1", "A1", 0, 0, false},
		{"A2", "A2", 1, 0, false},
		{"B1", "B1", 0, 1, false},
		{"AB33", "AB33", 32, 27, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row, col, err := GetCellCoords(tt.input)
			if row != tt.row {
				t.Errorf("got %v, want %v", row, tt.row)
			}
			if col != tt.column {
				t.Errorf("got %v, want %v", col, tt.column)
			}
			if (err != nil) != tt.err {
				t.Errorf("got %v, want %v", err, tt.err)
			}
		})
	}
}

func TestGetCoordsCell(t *testing.T) {
	var tests = []struct {
		name   string
		row    int32
		column int32
		want   string
	}{
		{"A1", 0, 0, "A1"},
		{"A2", 1, 0, "A2"},
		{"B1", 0, 1, "B1"},
		{"AB33", 32, 27, "AB33"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetCoordsCell(tt.row, tt.column)
			if result != tt.want {
				t.Errorf("got %v, want %v", result, tt.want)
			}
		})
	}
}
