package bod

import (
	"math/rand/v2"
)

// # Binary Space Partitioning (BSP) Overview! #
// Split the screen into smaller slices with random cuts, alternating between vertical and horizontal slices.
// Track each split as a leaf, when a leaf is split assign the current as the parent of the new leaves and
// assign the new leaves as the children of the current leaf.
// Once a leaf has reached a minimum size stop splitting.
// Once you have reached n depth of splitting, stop splitting.
// In DSA the structure you have just made is a tree.
// Find the bottom leaves, make randomly sized rectangles within the bounds of the levels.
// Find the sibling leaf of each leaf and connect them with a corridor. I do that by grid searching to find all
// the pairs of points between each sibling room's coordinates that tie for nearest distance, I then randomly
// select from the ties to make a corridor.
// Now track all the rooms that could not find a sibling, they are not of the same depth or "generation"
// of the rooms you found at the bottom of your tree structure, but once we join by grandparent instead of
// parent, you MAY have a match. So join by grandparent, now we have cousins - repeat the corridor creation routine
// and track the no-cousin leaves, keep on repeating until you reach the start of the tree/root node.
// Eventually everything will be connected this way.

// This BSP also has code to animate itself, basically I show changes to the dungeon as the for loops run.
// I have a function to print the matrix which is just a matrix of single char strings
// and then make the program sleep.

// Questions? Hit me up on Reddit.

// aforementioned leaves
type leaf struct {
	x1       int
	x2       int
	y1       int
	y2       int
	level    int
	children []leaf
	parent   *leaf
	letter   string
	room     room
}

// rooms inside leaves
type room struct {
	x1 int
	x2 int
	y1 int
	y2 int
}

// horizontal split, buffer regulates split size, if this the first split buffer is way bigger
func makeHSplit(l *leaf, screenHeight int, m [][]string) {

	buffer := 6
	if l.level == 0 {
		buffer = screenHeight / 3
	}

	top := (l.y2 - l.y1) - buffer
	hcut := (rand.IntN(top-buffer) + buffer)

	childa := leaf{
		x1: l.x1, y1: l.y1,
		x2: l.x2, y2: l.y1 + hcut,
		level:    l.level + 1,
		children: nil, parent: l,
		letter: "a",
	}

	childb := leaf{
		x1: l.x1, y1: l.y1 + hcut + 1,
		x2: l.x2, y2: l.y2,
		level:    l.level + 1,
		children: nil, parent: l,
		letter: "b",
	}

	l.children = append(l.children, childa)
	l.children = append(l.children, childb)

	for x := l.x1; x < l.x2; x++ {
		m[l.y1+hcut][x] = "#"
	}
}

// vertical slice
func makeVSplit(l *leaf, m [][]string) {

	buffer := 6

	top := (l.x2 - l.x1) - buffer
	vCut := (rand.IntN(top-buffer) + buffer)

	childa := leaf{
		x1: l.x1, y1: l.y1,
		x2: l.x1 + vCut, y2: l.y2,
		level:    l.level + 1,
		children: nil, parent: l,
		letter: "a",
	}

	childb := leaf{
		x1: l.x1 + vCut + 1, y1: l.y1,
		x2: l.x2, y2: l.y2,
		level:    l.level + 1,
		children: nil, parent: l,
		letter: "b",
	}

	l.children = append(l.children, childa)
	l.children = append(l.children, childb)

	for y := l.y1; y < l.y2; y++ {
		m[y][l.x1+vCut] = "#"
	}
}

// recursively build tree
func doSplits(l *leaf, screenHeight int, m [][]string) {

	level := l.level

	w := l.x2 - l.x1
	h := l.y2 - l.y1

	if w*h <= 20 {
		return
	}

	//even of h and v splits, check if there is space budget given buffer values
	if level%2 != 0 {
		if w >= 13 {
			makeVSplit(l, m)
		} else if h >= 13 {
			makeHSplit(l, screenHeight, m)
		} else {
			return
		}
	} else if h >= 13 {
		makeHSplit(l, screenHeight, m)
	} else if w >= 13 {
		makeVSplit(l, m)
	} else {
		return
	}

	if level == 5 {
		return
	}

	doSplits(&l.children[0], screenHeight, m)
	doSplits(&l.children[1], screenHeight, m)

}

// really should be "get terminal leaves" ... whoops
// this project is almost done and we aren't changing it here.
// any how this grabs all the leaves with no kids using recursion
// note bucket = getLevel(&leaves[0],bucket) doesnt dump the bucket out
// the bucket keeps what's been appended as we recurse
func getLevel(l *leaf, bucket []leaf) []leaf {

	leaves := l.children

	if leaves[0].children == nil {
		bucket = append(bucket, leaves[0])
	} else {
		bucket = getLevel(&leaves[0], bucket)
	}
	if leaves[1].children == nil {
		bucket = append(bucket, leaves[1])
	} else {
		bucket = getLevel(&leaves[1], bucket)
	}

	return bucket

}

// make a roomly sized room by find out the budget based on the leaf size and set a floor on how low
// we can go in any dimension
func makeRoom(l *leaf) {

	couldAdd := (l.x2 - 1) - (l.x1 + 1)
	addX := rand.IntN(couldAdd)
	if addX < couldAdd/3 {
		addX = couldAdd / 3
	}
	if addX < 4 {
		addX = 4
	}

	leftOverX := couldAdd - addX
	jitterX := 0
	if leftOverX > 0 {
		jitterX = rand.IntN(leftOverX)
		if jitterX == 0 {
			jitterX = 1
		}
	}

	couldAdd = (l.y2 - 1) - (l.y1 + 1)
	addY := rand.IntN(couldAdd)
	if addY < couldAdd/2 {
		addY = couldAdd / 2
	}
	if addY < 4 {
		addY = 4
	}

	leftOverY := couldAdd - addY
	jitterY := 0
	if leftOverY > 0 {
		jitterY = rand.IntN(leftOverY)

	}

	roomY1 := l.y1 + 1 + jitterY
	roomX1 := l.x1 + 1 + jitterX

	roomY2 := roomY1 + addY - 1
	roomX2 := roomX1 + addX - 1

	rm := room{
		x1: roomX1,
		x2: roomX2,
		y1: roomY1,
		y2: roomY2,
	}

	l.room = rm

}

// loop ove the room coords and assign them to the main map
func mapRoom(l leaf, m [][]string) {
	rm := l.room
	for y := rm.y1; y <= rm.y2; y++ {
		for x := rm.x1; x <= rm.x2; x++ {
			m[y][x] = "0"
		}
	}
}

// given a slice of leaves, find the one with the same parent, without having dupe entries
// man I remember getting stumped on a coding interview question for a problem like this
// 350k job nbd... anyway the no dupes part is handled by the second level loop index [idx1+1:]
// we skip everything at and below the current index of loop one.
// once you have checked all the possible pairings for idx one, you dont need to ask about it again so dont
// look for it when you are on idx 2, also when on idx 3 dont look for pairings containing 1 or 2.
// maybe I just made you 350k... cheers.
func sibPairing(l []leaf) [][]leaf {

	bucket := [][]leaf{}

	for idx1, val1 := range l {
		for _, val2 := range l[idx1+1:] {
			if val1.parent == val2.parent && val1.parent != nil {
				bucket = append(bucket, []leaf{val1, val2})
			}
		}
	}

	return bucket
}

// this is a grid search described at the top
// grid search to find all the pairs of points between
// each sibling room's coordinates that tie for nearest distance, then randomly
// select from the ties to make a corridor.
func nearestPair(l1, l2 leaf, m [][]string) (int, int, int, int) {

	rm1 := l1.room
	rm2 := l2.room

	bestX0 := rm1.x1
	bestY0 := rm1.y1
	bestX1 := rm2.x1
	bestY1 := rm2.y1
	bestDist := Manhattan(bestX0, bestY0, bestX1, bestY1)
	//find actual best
	for x0 := rm1.x1; x0 <= rm1.x2; x0++ {
		for y0 := rm1.y1; y0 <= rm1.y2; y0++ {
			if m[y0][x0] == "0" {
				for x1 := rm2.x1; x1 <= rm2.x2; x1++ {
					for y1 := rm2.y1; y1 <= rm2.y2; y1++ {
						if m[y1][x1] == "0" {
							dist := Manhattan(x0, y0, x1, y1)
							if dist < bestDist {
								bestX0 = x0
								bestY0 = y0
								bestX1 = x1
								bestY1 = y1
								bestDist = dist
							}
						}
					}
				}
			}
		}
	}
	//find ties
	bestPairs := [][]int{}
	for x0 := rm1.x1; x0 <= rm1.x2; x0++ {
		for y0 := rm1.y1; y0 <= rm1.y2; y0++ {
			if m[y0][x0] == "0" {
				for x1 := rm2.x1; x1 <= rm2.x2; x1++ {
					for y1 := rm2.y1; y1 <= rm2.y2; y1++ {
						if m[y1][x1] == "0" {
							dist := Manhattan(x0, y0, x1, y1)
							if dist == bestDist {
								bestPairs = append(bestPairs, []int{x0, y0, x1, y1})
							}
						}
					}
				}
			}
		}
	}
	//randomly select the pair to use
	winner := bestPairs[rand.IntN(len(bestPairs))]
	return winner[0], winner[1], winner[2], winner[3]
}

// should be tunnel merge collect
// we take a slice of sibs and make corridors, then promote the siblings to parent... or make a fraud parent
// then we add the fraud parent to a slice so we have the next level to pair up
func mergeCollect(sibs [][]leaf, m [][]string) []leaf {

	bucket := []leaf{}

	//tunnel
	for _, pair := range sibs {
		//x0,y0 := pair[0].center()
		//x0,y0 = mostCenter(x0,y0,pair[0],m)
		//x1,y1 := pair[1].center()
		//x1,y1 = mostCenter(x1,y1,pair[1],m)
		x0, y0, x1, y1 := nearestPair(pair[0], pair[1], m)

		xMin := x0
		if x0 > x1 {
			xMin = x1
		}

		xMax := x0
		if x0 < x1 {
			xMax = x1
		}

		yMin := y0
		if y0 > y1 {
			yMin = y1
		}

		yMax := y0
		if y0 < y1 {
			yMax = y1
		}

		for x := xMin; x < xMax; x++ {
			m[yMin][x] = "0"
		}
		for y := yMin; y < yMax; y++ {
			m[y][xMax] = "0"
		}

		//merge
		//placeholder
		parent := &pair[0]
		if pair[0].parent != nil {
			parent = pair[0].parent.parent
		} else {
			parent = nil
		}
		mother := leaf{
			room: room{
				x1: xMin,
				y1: yMin,
				x2: xMax,
				y2: yMax,
			},
			children: []leaf{pair[0], pair[1]},
			parent:   parent,
		}
		//collect
		bucket = append(bucket, mother)

	}
	return bucket

}

// get the leaves without siblings
func getOnlyChildren(l []leaf) []leaf {
	//singles get promoted to mega rooms... is what I used to call this
	singles := []leaf{}
	for idx1, leaf1 := range l {
		sibCount := 0
		for idx2, leaf2 := range l {
			if leaf1.parent == leaf2.parent && idx1 != idx2 {
				sibCount++
			}
		}
		if sibCount == 0 {
			if leaf1.parent == nil {
				singles = append(singles, leaf1)
				continue
			}
			mr := leaf{
				x1:       leaf1.x1,
				y1:       leaf1.y1,
				x2:       leaf1.x2,
				y2:       leaf1.y2,
				room:     leaf1.room,
				children: []leaf{leaf1},
				parent:   leaf1.parent,
			}
			singles = append(singles, mr)
		}
	}
	return singles
}

func DoBSP(screenWidth, screenHeight int) [][]string {
	root := leaf{
		x1: 0, y1: 0,
		x2:       screenWidth,
		y2:       screenHeight,
		level:    0,
		children: nil,
		parent:   nil,
	}
	//"the main map"
	rows := make([][]string, screenHeight)
	//fill map with "."
	for i := 0; i < screenHeight; i++ {
		row := make([]string, screenWidth)
		for i := range row {
			row[i] = "."
		}
		rows[i] = row
	}

	doSplits(&root, screenHeight, rows)

	//make leaf bucket fill leaf bucket wih terminal leaves
	leafBucket := []leaf{}
	leafBucket = getLevel(&root, leafBucket)
	for idx := range leafBucket {
		makeRoom(&leafBucket[idx])
		mapRoom(leafBucket[idx], rows)

	}

	EditMapChars(rows, "#", ".", true)

	current := leafBucket
	//while len of current is more than 1
	// pair up the sibling leaves, find the leaves with no siblings
	// run mergecollect to get the next gen to pair up and throw those leaves in a slice with the singles
	// ... after a slices means make these into elements... once len(current)==1 we are at the root node.
	for len(current) > 1 {
		sibs := sibPairing(current)
		singles := getOnlyChildren(current)

		current = append(
			mergeCollect(sibs, rows),
			singles...,
		)
	}

	////section for cleaning and showing map

	//turn the dead space to blank spaces
	EditMapChars(rows, ".", " ", true)
	//turn the rooms into more traditional floor chars
	EditMapChars(rows, "0", ".", true)

	//make walls, check if a cell is blank and is surronded by any floor chars
	for idx1 := range rows {
		for idx2 := range rows[idx1] {
			if rows[idx1][idx2] != " " {
				continue
			}
			n := NCharHood(rows, idx2, idx1, ".")
			if n > 0 {
				rows[idx1][idx2] = "#"
			}
		}

	}

	//make a deep copy of the map.
	//slices have funny pointer stuff going on under the hood so when you set something equal to a slice
	//editing the new slice impacts the old one.
	twinMap := make([][]string, len(rows))
	for i := range twinMap {
		twinMap[i] = append([]string(nil), rows[i]...)
	}

	//cleaning the map, sometimes making a corridor looks awkward then tunneling over other corridors and rooms
	//specifically the corridor walls cover over things in a visually unpleasing way
	//To tidy up we start by looking for the pattern below, and changing the middle char to "."
	//    .
	//   ###
	//    .
	//to start we ask if the current cell is "#" if not then skip this iteration
	//if so then check if we can check both above AND below, if not then the check will be false for both above AND below
	//if so then we actually the values for each, looking for "."
	//same routine for left and right, looking for "#"
	//finally if all the checks are true we go to our twin map and change the "#" to "."
	//after the routine we set rows(our map) equal to our twinMap
	//This way we can make edits while iterating without changing the map, which would change the outcome of our checks
	//as a by product of this procces we get island walls so we turn those into floor as well

	for idx1 := range rows {
		for idx2 := range rows[idx1] {
			//do even need to check?
			if rows[idx1][idx2] != "#" {
				continue
			}
			//default to false
			above := false
			below := false
			//if we can check then we can check
			if idx1 > 0 && idx1 < len(rows)-1 {
				above = rows[idx1-1][idx2] == "."
				below = rows[idx1+1][idx2] == "."
			}
			//same pattern as above
			left := false
			right := false
			if idx2 > 0 && idx2 < len(rows[idx1])-1 {
				left = rows[idx1][idx2-1] == "#"
				right = rows[idx1][idx2+1] == "#"
			}

			//store check if all true
			allTrue := above && below && left && right
			//edit twin map
			if allTrue {
				twinMap[idx1][idx2] = "."
			}
		}

	}

	rows = twinMap

	//check for those wall islands
	for idx1 := range rows {
		for idx2 := range rows[idx1] {
			//do we even need to check?
			if rows[idx1][idx2] != "#" {
				continue
			}
			//island = "." char neighbors = 4
			n := NCharHood(rows, idx2, idx1, ".")
			if n == 4 {
				rows[idx1][idx2] = "."
			}
		}
	}

	return rows

}
