// Gitbrowse: a simple web server for git.
// Copyright (C) 2026 Vithushan
// 
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
// 
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package template

import (
	"fmt"
	"log"
	"math"
	"git.lewoof.xyz/gitbrowse/config"
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
