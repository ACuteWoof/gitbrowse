package template

import (
	"fmt"
	"log"
	"math"
	"git.lewoof.xyz/clone/gitbrowse/config"
)

func getFormattedSize(size float64) (formatted string) {
	var precisionDivisor = math.Pow(10, config.SizePrecision)

	if size < 1024 {
		formatted = fmt.Sprintf("%g B", size)
	} else {
		kb := size / 1024
		if kb < 1024 {
			formatted = fmt.Sprintf("%g KiB", math.Round(kb*precisionDivisor)/precisionDivisor)
		} else {
			mb := kb / 1024
			if mb < 1024 {
				formatted = fmt.Sprintf("%g MiB", math.Round(mb*precisionDivisor)/precisionDivisor)
			} else {
				gb := mb / 1024
				formatted = fmt.Sprintf("%g GiB", math.Round(gb*precisionDivisor)/precisionDivisor)
			}
		}
	}
	return
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
