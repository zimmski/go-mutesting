// +build example-main

package example

func foo() int {
	n := 1

	for i := 0; i < 3; i++ {
		if i == 0 {
			n++
		} else if i == 1 {
			n += 2
		} else {
			n += 3
		}

		n++
	}

	if n < 0 {
		n = 0
	}

	n++

	n += bar()

	bar()
	bar()

	switch {
	case n < 20:
		n++
	case n > 20:
		n--
	default:

	}

	return n
}

func bar() int {
	return 4
}
