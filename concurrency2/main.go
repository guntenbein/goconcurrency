package main

import (
	"fmt"
)

type Figure struct {
	Length int
	Width  int
	Height int
	Square int
	Volume int
}

const n = 2

func main() {

	ff := []Figure{
		Figure{1, 2, 5, 0, 0},
		Figure{3, 2, 4, 0, 0},
		Figure{1, 10, 3, 0, 0}}

	squarec := make(chan Figure, n)

	volumec := make(chan Figure, n)

	go func() {
		computeSquare(ff, squarec)
	}()

	go func() {
		computeVolume(squarec, volumec)
	}()

	send(volumec)
}

func computeSquare(ff []Figure, squarec chan<- Figure) {
	for _, f := range ff {
		f.Square = f.Length * f.Width
		squarec <- f
	}
	close(squarec)
}

func computeVolume(squarec <-chan Figure, volumec chan<- Figure) {
	for f := range squarec {
		f.Volume = f.Square * f.Height
		volumec <- f
	}
	close(volumec)
}

func send(sourcec <-chan Figure) error {
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
	return nil
}
