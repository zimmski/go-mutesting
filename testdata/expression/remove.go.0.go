// +build example-main

package main

import (
	"fmt"
)

func main() {
	i := 1

	for i != 4 {
		if true && i <= 1 {
			fmt.Println(i)
		} else if (i >= 2 && i <= 2) || i*1 == 1+1 {
			fmt.Println(i * 2)
		} else {

		}

		i++
	}
}
