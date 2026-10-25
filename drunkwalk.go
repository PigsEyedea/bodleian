package bod

import (
	"math/rand/v2"
	"slices"
	"strings"
)

var oldX int
var oldY int

var str strings.Builder
var totalTiles int
var filledTiles int
var targetPercent float64
var targetFilledTiles int

// this random walk, aka drunk walk. Drunk walk sounds more fun so we have a drunkard not a walker
type Drunkard struct {
	x     int
	y     int
	upDwn int
	up    int
	left  int
}

// stop drunkard from falling off the map
func clamp(drunky *Drunkard, width, height int) {
	if drunky.x < 0 {
		drunky.x = 0
	}
	if drunky.y < 0 {
		drunky.y = 0
	}
	if drunky.y > height-2 {
		drunky.y = height - 2
	}
	if drunky.x > width-1 {
		drunky.x = width - 1
	}
}

// walk one pace
func stumble(drunky *Drunkard, width, height, momentum int) {

	//momentum check
	//makes drunky repeat a step
	//if the rand value is greater than set momentum param then give drunky a new direction shot
	momentumRoll := rand.IntN(11)
	if momentumRoll > momentum {
		drunky.upDwn = rand.IntN(2)
		drunky.up = 0
		drunky.left = 0
	}

	//init helper vars
	upDwn := drunky.upDwn
	up := drunky.up
	left := drunky.left

	//rand int 2 means randomly pick 0 or 1
	//we use it like random bool
	//so, randomly pick if we are going left or right in case where are not going up or down
	//if we are going up or down rand int to say which one of the two
	//increment the coords to make that happen
	//y is kinda whacky if you are used to basic coords, our print out is like table not a chart
	//(y gets bigger as we go south not smaller, so up is subract and down is add)
	if upDwn == 0 {
		left = rand.IntN(2)
		if left == 1 {
			drunky.x--
		} else {
			drunky.x++
		}
	} else {
		up = rand.IntN(2)
		if up == 1 {
			drunky.y--
		} else {
			drunky.y++
		}
	}
	drunky.left = left
	drunky.up = up
	//if our  drunky just stumbled off the face of the planet quickly yank him back
	clamp(drunky, width, height)
}

// in the event the tile has already been stepped on, try for an empty tile, as many times as you like(the rolls param)
func reRoll(drunky *Drunkard, rolls, width, height, momentum int, matrix [][]string) {
	char := matrix[drunky.y][drunky.x]
	for i := 0; i < rolls; i++ {
		if char != " " {
			drunky.x = oldX
			drunky.y = oldY
			stumble(drunky, width, height, momentum)
			char = matrix[drunky.y][drunky.x]
		} else {
			break
		}
	}
}

// helper func to DRY up free space and boarder checker func
func checkNeighbors(data [][]string, x, y, width, height int, assumption bool) []bool {

	isMinX := x == 0
	isMinY := y == 0
	isMaxX := x == width-1
	isMaxY := y == height-2

	freeLeft := assumption
	freeRight := assumption
	freeUp := assumption
	freeDown := assumption

	if isMinX == false {
		freeLeft = data[y][x-1] == " "
	}

	if isMinY == false {
		freeUp = data[y-1][x] == " "
	}

	if isMaxX == false {
		freeRight = data[y][x+1] == " "
	}

	if isMaxY == false {
		freeDown = data[y+1][x] == " "
	}

	freeSlc := []bool{freeLeft, freeRight, freeUp, freeDown}
	return freeSlc

}

// flesh out map later by tracing around the map to make walls
func boarder(data [][]string, x, y, width, height int) bool {
	assum := true
	freeSlc := checkNeighbors(data, x, y, width, height, assum)
	boarderTile := false
	if slices.Contains(freeSlc, false) {
		boarderTile = true
	}
	return boarderTile
}

// see if there are any unplaced tiles around the drunkard
func free(drunky *Drunkard, width, height int, matrix [][]string) bool {
	freeSlc := checkNeighbors(matrix, drunky.x, drunky.y, width, height, false)
	return slices.Contains(freeSlc, true)
}

func DoDrunkWalk(width, height, momentum int, targetPercent float64) [][]string {

	var matrix [][]string

	totalTiles = height * width
	filledTiles = 1
	targetFilledTiles = int(float64(totalTiles) * targetPercent)

	//init drunkard, his name is Paul, hi Paul
	paul := Drunkard{
		x:     rand.IntN(width),
		y:     rand.IntN(height - 1),
		upDwn: rand.IntN(2),
		up:    0,
		left:  0}

	//matrix init
	for i := 0; i < height-1; i++ {
		slc := make([]string, width)
		for idx := range slc {
			slc[idx] = " "
		}
		matrix = append(matrix, slc)
	}

	nextTileChar := matrix[paul.y][paul.x]

	//stumbling phase
	for {
		//place
		matrix[paul.y][paul.x] = "x"
		if filledTiles >= targetFilledTiles {
			break
		}

		oldX = paul.x
		oldY = paul.y
		lastTileFree := free(&paul, width, height, matrix)

		stumble(&paul, width, height, momentum)

		nextTileChar = matrix[paul.y][paul.x]

		if nextTileChar != " " && lastTileFree == true {
			reRoll(&paul, 5, width, height, momentum, matrix)
		}
		// Paul's position may have changed, so check again
		nextTileChar = matrix[paul.y][paul.x]
		if nextTileChar == " " {
			filledTiles++
		}
	}

	return matrix

}
