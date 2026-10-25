package bod

import "fmt"

// abs value of an int
func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// simple manhattan distance formula
func Manhattan(x1, y1, x2, y2 int) int {
	return AbsInt(x2-x1) + AbsInt(y2-y1)
}

// go through the matrix and see how spaces around a cell have a given char in them
// helps with assigning chars to viz the matrix
func NCharHood(m [][]string, x, y int, char string) int {
	//n chars = x in neighboorhood
	count := 0
	if y > 0 && m[y-1][x] == char {
		count++
	}

	if y < len(m)-1 && m[y+1][x] == char {
		count++
	}

	if x > 0 && m[y][x-1] == char {
		count++
	}

	if x < len(m[y])-1 && m[y][x+1] == char {
		count++
	}

	return count
}

// loop over the map and print it
func ShowMap(m [][]string) {
	for _, val := range m {
		for _, v2 := range val {
			fmt.Print(v2)
		}
		fmt.Println()
	}
}

// find and replace for the map
func EditMapChars(m [][]string, find, insert string, show bool) {
	for idx1 := range m {
		for idx2 := range m[idx1] {
			if m[idx1][idx2] == find {
				m[idx1][idx2] = insert
			}
		}
	}
}
