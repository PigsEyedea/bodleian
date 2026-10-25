package tests

//all tests ai generated

import (
	"testing"

	bod "github.com/PigsEyedea/bodleian"
)

// Check that a generated map has the expected dimensions.
func checkDimensions(t *testing.T, m [][]string, width, height int) {
	t.Helper()

	if len(m) != height {
		t.Fatalf("expected %d rows, got %d", height, len(m))
	}

	for y, row := range m {
		if len(row) != width {
			t.Fatalf("row %d: expected %d columns, got %d",
				y, width, len(row))
		}
	}
}

// Count occurrences of a particular tile.
func countTiles(m [][]string, tile string) int {
	count := 0

	for _, row := range m {
		for _, cell := range row {
			if cell == tile {
				count++
			}
		}
	}

	return count
}

func TestBSP(t *testing.T) {
	const width, height = 60, 40

	for i := 0; i < 5; i++ {
		m := bod.DoBSP(width, height)

		checkDimensions(t, m, width, height)

		for y, row := range m {
			for x, tile := range row {
				switch tile {
				case " ", ".", "#":
					// Valid BSP tiles.
				default:
					t.Fatalf("unexpected tile %q at (%d,%d)",
						tile, x, y)
				}
			}
		}

		if countTiles(m, ".") == 0 {
			t.Fatal("BSP generated no floor tiles")
		}
	}
}

func TestDrunkWalk(t *testing.T) {
	const (
		width    = 30
		height   = 20
		momentum = 5
		target   = 0.40
	)

	for i := 0; i < 3; i++ {
		m := bod.DoDrunkWalk(width, height, momentum, target)

		// Current implementation creates height-1 rows.
		checkDimensions(t, m, width, height-1)

		filled := countTiles(m, "x")
		total := width * (height - 1)

		expected := int(float64(total) * target)

		if filled < expected {
			t.Errorf("expected at least %d filled tiles, got %d",
				expected, filled)
		}

		if filled > total {
			t.Errorf("filled tiles %d exceeds total %d",
				filled, total)
		}

		for y, row := range m {
			for x, tile := range row {
				if tile != "x" && tile != " " {
					t.Fatalf("unexpected tile %q at (%d,%d)",
						tile, x, y)
				}
			}
		}
	}
}

func TestManhattan(t *testing.T) {
	tests := []struct {
		name     string
		x1, y1   int
		x2, y2   int
		expected int
	}{
		{"same point", 0, 0, 0, 0, 0},
		{"horizontal", 0, 0, 5, 0, 5},
		{"vertical", 0, 0, 0, 7, 7},
		{"diagonal", 1, 2, 4, 6, 7},
		{"negative coordinates", -2, -3, 2, 3, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bod.Manhattan(tt.x1, tt.y1, tt.x2, tt.y2)

			if got != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, got)
			}
		})
	}
}

func TestAbsInt(t *testing.T) {
	tests := []struct {
		input, expected int
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
	}

	for _, tt := range tests {
		got := bod.AbsInt(tt.input)

		if got != tt.expected {
			t.Errorf("AbsInt(%d): expected %d, got %d",
				tt.input, tt.expected, got)
		}
	}
}

func TestEditMapChars(t *testing.T) {
	m := [][]string{
		{".", "#", "."},
		{"#", ".", "#"},
	}

	bod.EditMapChars(m, "#", "X", false)

	expected := [][]string{
		{".", "X", "."},
		{"X", ".", "X"},
	}

	for y := range m {
		for x := range m[y] {
			if m[y][x] != expected[y][x] {
				t.Errorf("tile (%d,%d): expected %q, got %q",
					x, y, expected[y][x], m[y][x])
			}
		}
	}
}

func TestNCharHood(t *testing.T) {
	m := [][]string{
		{".", "#", "."},
		{"#", ".", "#"},
		{".", "#", "."},
	}

	got := bod.NCharHood(m, 1, 1, "#")

	if got != 4 {
		t.Errorf("expected 4 neighboring walls, got %d", got)
	}

	// Corner has only two possible cardinal neighbors.
	got = bod.NCharHood(m, 0, 0, "#")

	if got != 2 {
		t.Errorf("expected 2 neighboring walls, got %d", got)
	}
}
