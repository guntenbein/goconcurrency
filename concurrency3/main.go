package main

import (
	"fmt"
	"sync"
)

type Figure struct {
	Length int
	Width  int
	Height int
	Square int
	Volume int
}

const (
	n           = 2
	statusOK    = 0
	statusError = 1
)

func main() {

	errc := make(chan error)
	status := statusOK

	errGroup := sync.WaitGroup{}
	errGroup.Add(1)

	go func() {
		for err := range errc {
			status = statusError
			fmt.Printf("error processing the code: %s\n", err)
		}
		errGroup.Done()
	}()

	ff := []Figure{
		Figure{1, 2, -5, 0, 0},
		Figure{3, 2, 4, 0, 0},
		Figure{1, 10, 3, 0, 0},
		Figure{1, 10, -3, 0, 0},
		Figure{-1, 10, 3, 0, 0},
		Figure{1, 10, 3, 0, 0}}

	squarec := make(chan Figure, n)

	volumec := make(chan Figure, n)

	go func() {
		computeSquare(ff, squarec, errc)
	}()

	go func() {
		computeVolume(squarec, volumec, errc)
	}()

	send(volumec, errc)

	close(errc)
	errGroup.Wait()
}

func computeSquare(ff []Figure, squarec chan<- Figure, errc chan<- error) {
	for _, f := range ff {
		if f.Length <= 0 || f.Width <= 0 {
			errc <- fmt.Errorf("invalid length or width value, should be positive non-zero, length: %d, width: %d", f.Length, f.Width)
		}
		f.Square = f.Length * f.Width
		squarec <- f
	}
	close(squarec)
}

func computeVolume(squarec <-chan Figure, volumec chan<- Figure, errc chan<- error) {
	var err error
	for f := range squarec {
		if f.Height <= 0 {
			err = fmt.Errorf("invalid height value, should be positive non-zero, height: %d", f.Height)
			errc <- err
		}
		// skip if error happens during previous figure calculation
		if err == nil {
			f.Volume = f.Square * f.Height
			volumec <- f
		}
	}
	close(volumec)
}

func send(sourcec <-chan Figure, errc chan<- error) {
	var err error
	count := 0
	batch := make([]Figure, 0, n)
	for f := range sourcec {
		if f.Volume > 25 {
			err = fmt.Errorf("cannot send figures with volume more than 25, volume: %d", f.Volume)
			errc <- err
		}
		// skip if error happens during sending
		if err == nil {
			batch = append(batch, f)
			count++
			if count == n {
				// imitate sending batch
				fmt.Println(batch)
				batch = make([]Figure, 0, n)
				count = 0
			}
		}
	}
	if err == nil && len(batch) != 0 {
		// imitate sending rest
		fmt.Println(batch)
	}
}
