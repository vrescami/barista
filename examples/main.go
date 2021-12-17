package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"gitlab.devtools.intel.com/vrescami/barista/pkg"
)

const (
	RED    = "\033[1;31m%s\033[0m"
	GREEN  = "\033[1;32m%s\033[0m"
	YELLOW = "\033[1;33m%s\033[0m"
	BLUE   = "\033[1;34m%s\033[0m"
)

func colorize(msg, color string) string {
	return fmt.Sprintf(color, msg)
}

func main() {
	bar := barista.New()

	for i := 0; i < 20; i++ {
		prefix := "c001n000" + strconv.Itoa(i) + ":" + colorize(" uploading", BLUE)
		suffix := colorize("done", GREEN)
		bar.Add(barista.NewBar(30, 30, prefix, suffix))
	}

	bar.Start(200 * time.Millisecond)
	defer bar.Stop()

	wg := sync.WaitGroup{}
	wg.Add(30 * 20)
	for i := 0; i < 30; i++ {
		for b := 0; b < 20; b++ {
			go func(b int) {
				defer wg.Done()
				time.Sleep(time.Duration(rand.Intn(10000)) * time.Millisecond)
				bar.Step(b, 1)
			}(b)
		}
	}

	wg.Wait()

	// give time to print the final bars
	time.Sleep(200 * time.Millisecond)
}
