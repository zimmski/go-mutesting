package example

func foo() int {
	n := 1

	for i := 0; i < 3; i++ {
		if i == 0 {
			n += 1
		} else if i == 1 {
			n += 2
		} else {
			n += 3
		}

		n++
	}

	n++

	n += bar()

	return n
}

func bar() int {
	return 4
}
