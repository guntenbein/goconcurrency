package main

import "fmt"

type Figure struct {
	Length int
	Width  int
	Square int
}

const n = 2

func main() {
	ff := []Figure{Figure{1, 2, 0}, Figure{3, 2, 0}, Figure{1, 10, 0}}

	squarec := make(chan Figure, n)

	go func() {
		computeSquare(ff, squarec)
	}()

	send(squarec)
}

func computeSquare(ff []Figure, squarec chan<- Figure) {
	for _, f := range ff {
		f.Square = f.Length * f.Width
		squarec <- f
	}
	close(squarec)
}

func send(sourcec <-chan Figure) {
	count := 0
	batch := make([]Figure, 0, n)
	for f := range sourcec {
		batch = append(batch, f)
		count++
		if count == n {
			// imitate sending batch
			fmt.Println(batch)
			batch = make([]Figure, 0, n)
			count = 0
		}
	}
	// imitate sending rest
	fmt.Println(batch)
}
