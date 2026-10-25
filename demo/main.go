package main

import (
	bod "github.com/PigsEyedea/bodleian"
)

// use all the functions to do bsp!
func main() {
	myMap := bod.DoBSP(60, 25)
	bod.ShowMap(myMap)
}
