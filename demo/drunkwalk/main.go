package main

import (
	"fmt"
	"strings"
	"time"

	bod "github.com/PigsEyedea/bodleian"
)

func main() {
	sleeptime := time.Millisecond * 350
	fill := .2
	//myMap := bod.DoDrunkWalk(50, 22, mo, fill)
	for i := range 10 {
		for j := range 4 {
			fmt.Print("\033[H\033[2J\033[3J")
			dots := strings.Repeat(".", j%4)
			outMsgT := fmt.Sprintf("Giant Is Stumbling%s", dots)
			coloredMsgT := "\033[32m" + outMsgT + "\033[0m"
			fmt.Println(coloredMsgT)
			myMap := bod.DoDrunkWalk(50, 22, i, fill)
			bod.ShowMap(myMap)
			outMsg := fmt.Sprintf("W:50|H:22|Momentum:%d/9", i)
			coloredMsg := "\033[32m" + outMsg + "\033[0m"
			fmt.Println(coloredMsg)
			time.Sleep(sleeptime)

		}
	}

}
