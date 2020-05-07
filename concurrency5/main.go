package main

import (
	"fmt"
	"github.com/guntenbein/goconcurrency/errworker"
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
		Figure{1, 2, 5, 0, 0},
		Figure{3, 2, 4, 0, 0},
		Figure{1, 10, 3, 0, 0},
		Figure{1, 10, -3, 0, 0},
		Figure{1, -10, 3, 0, 0},
		Figure{1, 10, 5, 0, 0}}

	squarec := make(chan Figure, n)

	volumec := make(chan Figure, n)

	go func() {
		if err := computeSquare(ff, squarec); err != nil {
			errc <- err
		}
		close(squarec)
	}()

	go func() {
		if err := computeVolume(squarec, volumec); err != nil {
			errc <- err
		}
		close(volumec)
	}()

	send(volumec, errc)

	close(errc)
	errGroup.Wait()
}

func computeSquare(ff []Figure, squarec chan<- Figure) error {
	ew := errworker.NewErrWorkgroup(2, true)
	for _, f := range ff {
		fClosure := f
		ew.Go(func() error {
			if fClosure.Length <= 0 || fClosure.Width <= 0 {
				return fmt.Errorf("invalid length or width value, should be positive non-zero, length: %d, width: %d", fClosure.Length, fClosure.Width)
			}
			fClosure.Square = fClosure.Length * fClosure.Width
			squarec <- fClosure
			return nil
		})
	}
	return ew.Wait()
}

func computeVolume(squarec <-chan Figure, volumec chan<- Figure) error {
	ew := errworker.NewErrWorkgroup(3, true)
	var err error
	for f := range squarec {
		fClosure := f
		ew.Go(func() error {
			if fClosure.Height <= 0 {
				err = fmt.Errorf("invalid height value, should be positive non-zero, height: %d", fClosure.Height)
				return err
			}
			fClosure.Volume = fClosure.Square * fClosure.Height
			volumec <- fClosure
			return nil
		})
	}
	return ew.Wait()
}

func send(sourcec <-chan Figure, errc chan<- error) {
	var err error
	count := 0
	batch := make([]Figure, 0, n)
	for f := range sourcec {
		if f.Volume > 40 {
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
