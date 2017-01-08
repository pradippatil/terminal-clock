package main

import (
	"flag"
	"strings"
	"time"

	runewidth "github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

func draw(midW, midH, w int, frame string) {
	lines := strings.Split(th.digits[frame], "\n")
	for x, line := range lines {
		y := 0
		for _, cell := range line {
			termbox.SetCell(midW+y+w, midH+x, cell, th.fg, th.bg)
			y++
		}
	}

}

func drawBox(x, y, w, h int, fg, bg termbox.Attribute) {
	// unicode box drawing chars around clock
	termbox.SetCell(x, y, '┌', fg, bg)
	termbox.SetCell(x, y+h, '└', fg, bg)
	termbox.SetCell(x+w, y, '┐', fg, bg)
	termbox.SetCell(x+w, y+h, '┘', fg, bg)
	fill(x+1, y, w-1, 1, termbox.Cell{Ch: '─', Fg: fg, Bg: bg})
	fill(x+1, y+h, w-1, 1, termbox.Cell{Ch: '─', Fg: fg, Bg: bg})
	fill(x, y+1, 1, h-1, termbox.Cell{Ch: '│', Fg: fg, Bg: bg})
	fill(x+w, y+1, 1, h-1, termbox.Cell{Ch: '│', Fg: fg, Bg: bg})
}

// fill portion of width w and height h with given Cell starting from co-ordinates x, y
func fill(x, y, w, h int, cell termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, cell.Ch, cell.Fg, cell.Bg)
		}
	}
}

func printString(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

var th Theme

func main() {
	//defer profile.Start(profile.CPUProfile).Stop()
	var delay int
	var T string
	var center bool
	var hourType int
	flag.IntVar(&delay, "d", 1, "refresh delay in seconds")
	flag.StringVar(&T, "t", "ansi", "choose ansi or electronics")
	flag.BoolVar(&center, "c", false, "show clock at center if true")
	flag.IntVar(&hourType, "h", 12, "24/12 hour format")
	flag.Parse()
	if T == "electronics" {
		th = electronics
	} else {
		th = ansi

	}

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	termbox.SetOutputMode(termbox.Output256)
	//draw(0, 0)

loop:
	for {
		select {
		case ev := <-eventQueue:
			if (ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc) || ev.Type == termbox.EventInterrupt {
				break loop
			}
		default:
			//midW, midH = 1, 1
			//h, m, s := time.Now().Clock()
			//const layout = "3:04:05PM"
			layout := "15:04:05"
			if hourType == 12 {
				layout = "3:04:05PM"
			}
			//const layout = "Monday, Jan 2, 2006 3:04:05PM"
			currentTime := time.Now().Format(layout)
			day := time.Now().Format("Monday, Jan 2, 2006")
			midW, midH := 1, 1
			if center {
				w, h := termbox.Size()
				midW, midH = (w-th.chWidth*len(currentTime))/2, (h-th.chHeight)/2
			}
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			for i, v := range currentTime {
				draw(midW, midH, th.chWidth*i, string(v))

			}
			drawBox(midW-1, midH, th.chWidth*len(currentTime), th.chHeight, th.fg, th.bg)
			printString(midW, midH, termbox.ColorBlack, termbox.ColorMagenta, day)

			termbox.Flush()
			time.Sleep(time.Duration(delay) * time.Second)
		}
	}
}
