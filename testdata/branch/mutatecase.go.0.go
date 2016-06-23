// +build example-main

package main

import (
	"fmt"
)

func main() {
	i := 1

	for i != 4 {
		switch {
		case i == 1:
			_, _ = fmt.Println, i
		case i == 2:
			fmt.Println(i * 2)
		default:
			fmt.Println(i * 3)
		}

		i++
	}
}
