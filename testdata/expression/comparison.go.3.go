// +build example-main

package main

import "fmt"

func main() {
	if 1 > 2 {
		fmt.Printf("1 is greater than 2!")
	}

	if 1 < 2 {
		fmt.Printf("1 is less than 2!")
	}

	if 1 <= 2 {
		fmt.Printf("1 is less than or equal to 2!")
	}

	if 1 > 2 {
		fmt.Printf("1 is greater than or equal to 2!")
	}

	if 1 == 2 {
		fmt.Print("1 is equal to 2!")
	}
}
