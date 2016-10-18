package main

import (
	"github.com/gizak/termui"
	"github.com/glaslos/etsz"
)

func makeSpark(data []int) *termui.Sparklines {
	spl3 := termui.NewSparkline()
	spl3.Data = data
	spl3.Title = "Enlarged Sparkline"
	spl3.Height = 8
	spl3.LineColor = termui.ColorYellow

	spls := termui.NewSparklines(spl3)
	spls.Height = 11
	spls.Width = 30
	spls.BorderFg = termui.ColorCyan
	spls.X = 21
	spls.BorderLabel = "Tweeked Sparkline"
	return spls
}

// TODO (glaslos): PoC, don't read this...
func main() {
	edb := etsz.New()

	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	// handle key i pressing, inserting entry to time series
	termui.Handle("/sys/kbd/i", func(termui.Event) {
		edb.Insert(1.0, "test")
	})

	// Handle q, closing the UI
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	data := edb.ReadInt("test")
	spls := makeSpark(data)
	termui.Render(spls)

	// handle a 1s timer
	termui.Handle("/timer/1s", func(e termui.Event) {

		t := e.Data.(termui.EvtTimer)
		// t is a EvtTimer
		if t.Count%2 == 0 {
			data := edb.ReadInt("test")
			spls := makeSpark(data)
			termui.Render(spls)
		}
	})

	termui.Loop()
}
