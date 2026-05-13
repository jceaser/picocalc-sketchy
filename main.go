//go:build tinygo

/************************************************************************************************100
A RPN calculator based on the gotools calculator I built some time ago

Compile with:
	 tinygo build -o tinygo_sketchy.uf2 -target=pico
	 tinygo build -o tinygo_sketchy2.uf2 -target=pico2

go.mod use to use: github.com/jceaser/picocalc-common v0.0.0-20260506162810-f084c2574cf0

***************************************************************************************************/

package main

import (
	//standard modules
	"fmt"
	"image/color"
	"time"

	//tiny go modules
	"machine"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"

	//the picocalc modules
	"github.com/jceaser/picocalc-common/i2ckbd"
	"github.com/jceaser/picocalc-common/ili948x"
)

// ***************************************************************************80
// MARK Globals

var (
	black = color.RGBA{0, 0, 0, 255}
	red   = color.RGBA{255, 0, 0, 255}
	green = color.RGBA{0, 255, 0, 255}
	blue  = color.RGBA{0, 0, 255, 255}
	white = color.RGBA{255, 255, 255, 255}
)

var font = &freemono.Regular12pt7b // Regular18pt7b
var lcd *ili948x.Ili948x
var adcPin machine.Pin
var FontHeight, FontHeightOuter, FontWidth, FontWidthOuter int16
var lastDebug string

func init() {
	i, o := tinyfont.LineWidth(font, "X")
	FontWidth = int16(i)
	FontWidthOuter = int16(o)

	FontHeight = int16(12)
	FontHeightOuter = FontHeight + 4
}

// ***************************************************************************80
// MARK - functions

func showError(err error) {
    tinyfont.WriteLine(lcd, font, 0, 64, err.Error(), red)
    debug("Error", err.Error())
}

func debug(label string, value any) {
	if len(lastDebug) > 0 {
		if lastDebug == label {
			return
		}
		tinyfont.WriteLineRotated(lcd, font, 298, 300, lastDebug, black, tinyfont.ROTATION_270)
	}
	display := fmt.Sprintf("%s = %v", label, value)
	tinyfont.WriteLineRotated(lcd, font, 298, 300, display, white, tinyfont.ROTATION_270)
	lastDebug = display
}

func SetPixel(lcd *ili948x.Ili948x, x, y int16, fg color.RGBA) {
	x = between(0, x, 299)
	y = between(0, y, 299)
	lcd.SetPixel(x, y, fg)
}

func between(low, value, high int16) int16 {
	return max(low, min(value, high))
}

func max(left, right int16) int16 {
	if left > right {
		return left
	}
	return right
}

func min(left, right int16) int16 {
	if left < right {
		return left
	}
	return right
}

func nextColor(currentColor color.RGBA) color.RGBA {

	switch currentColor {
		case black:
			return red
		case red:
			return green
		case green:
			return blue
		case blue:
			return white
		case white:
			return black
	}
	return white
}

// ***************************************************************************80
// MARK - Main

// remember 320x320 screen size
func main() {
	machine.Watchdog.Start()
	lcd = ili948x.InitDisplay()

	tinyfont.WriteLineRotated(lcd, font, 300, 1, "Sketchy", white, tinyfont.ROTATION_90)

	var keyboard i2ckbd.I2CKbd
	if err := keyboard.Init(); err != nil {
		showError(err)
	}

	tickSleep := time.Now().Unix()
	tick := 0
	var x, y, lastx, lasty int16
	maxx, maxy := int16(299), int16(299)
	delta := int16(2)
	activeColor := white
	for {
		tick++
		k, err := keyboard.GetChar()
		if err != nil {
			showError(err)
		}
		if k == 0 {
			//pass
			now := time.Now().Unix()
			//debug("ts", now-tickSleep)
			if (now - tickSleep) > 5 {
				//debug("k", now - tickSleep)
				machine.Watchdog.Update()
				time.Sleep(750 * time.Millisecond)
				machine.Watchdog.Update()
				continue
			}
		} else if k == i2ckbd.ESC_KEY || k == 0x09 {
			debug("escape/tab", k)
			activeColor = nextColor(activeColor)
		} else if k == i2ckbd.LEFT_KEY {
			x = max(0, x-delta)
		} else if k == i2ckbd.UP_KEY {
			y = max(0, y-delta)
		} else if k == i2ckbd.DOWN_KEY {
			y = min(y+delta, maxy)
		} else if k == i2ckbd.RIGHT_KEY {
			x = min(x+delta, maxx)
		} else if k == i2ckbd.HOME_KEY {
			x, y = 0, 0
			lastx, lasty = 0, 0
		} else if k == i2ckbd.DEL_KEY {
			lcd.FillRectangle(0, 0, 320, 320, ili948x.BLACK)
		} else if k == i2ckbd.END_KEY {
			x, y = 299, 299
			lastx, lasty = 299, 299
		} else {
			//pass
		}

		if x != lastx || y != lasty {
			//wider then tall
			for xp := x-1 ; xp < x+1 ; xp++ {
				for yp := y-1 ; yp <= y+1 ; yp++ {
					SetPixel(lcd, xp, yp, activeColor)
				}
			}
			lastx = x
			lasty = y
			tickSleep = time.Now().Unix()
		}
		machine.Watchdog.Update()
		time.Sleep(10 * time.Millisecond)
		machine.Watchdog.Update()
	}
}
