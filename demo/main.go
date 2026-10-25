package main

import (
	"fmt"

	bod "github.com/PigsEyedea/bodleian"
)

func main() {
	myMap := bod.DoBSP(60, 25)
	bod.ShowMap(myMap)
	fmt.Println()
	myMap = bod.DoDrunkWalk(60, 25, 6, .2)
	bod.ShowMap(myMap)
}
