package barista

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	CLEAN_LINE = "\x1b[1A\x1b[2K"
)

type Bar struct {
	length  uint32
	total   uint32
	prefix  string
	suffix  string
	current uint32
}

func NewBar(length, total uint32, prefix, suffix string) Bar {
	return Bar{
		length,
		total,
		prefix,
		suffix,
		0,
	}
}

func (b *Bar) step(increase uint32) {
	b.current += increase
	if b.current > b.total {
		b.current = b.total
	}
}

func (b Bar) printBar() {
	percent := float64(b.current) * 100.0 / float64(b.total)
	filledLength := int(b.length * b.current / b.total)
	end := ">"
	sfx := ""

	if b.current >= b.total {
		end = "="
		sfx = b.suffix
	}
	bar := strings.Repeat("=", filledLength) + end + strings.Repeat(" ", int(b.length-uint32(filledLength)))
	fmt.Printf("%s [%s] %.2f%% %s\n", b.prefix, bar, percent, sfx)
}

type Barista struct {
	bars         []Bar
	increaseChan chan barIncrease
	stopChan     chan struct{}
	wg           sync.WaitGroup
}

func New() Barista {
	return Barista{}
}

func (b *Barista) Add(bar Bar) {
	b.bars = append(b.bars, bar)
}

func (b *Barista) Step(index int, increase uint32) {
	if b.increaseChan != nil {
		b.increaseChan <- barIncrease{index, increase}
	}
}

type barIncrease struct {
	index    int
	increase uint32
}

func (b *Barista) Start(refreshTime time.Duration) {
	ticker := time.NewTicker(refreshTime)
	b.printBars()
	b.increaseChan = make(chan barIncrease, 128)
	b.stopChan = make(chan struct{})
	b.wg = sync.WaitGroup{}

	b.wg.Add(1)
	go func() {
		defer ticker.Stop()
		defer close(b.increaseChan)
		defer b.wg.Done()
		for {
			select {
			case <-ticker.C:
				b.clean()
				b.printBars()
			case bar, ok := <-b.increaseChan:
				if ok {
					b.bars[bar.index].step(bar.increase)
				}
			case <-b.stopChan:
				return
			}
		}
	}()
}

func (b *Barista) Stop() {
	close(b.stopChan)
	b.wg.Wait()
}

func (b Barista) clean() {
	for _ = range b.bars {
		fmt.Printf(CLEAN_LINE)
	}
}

func (b Barista) printBars() {
	for _, bar := range b.bars {
		bar.printBar()
	}
}
