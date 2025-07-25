package utils

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"sync"
)

const spacing = 2.0

var font rl.Font
var once sync.Once

func getFont() *rl.Font {
	once.Do(func() {
		font = rl.LoadFontEx("assets/fonts/Orbitron-Regular.ttf", 32, nil, 250)
	})
	return &font
}

// CenterText draws the given text centered around the passed-in position
func CenterText(text string, position rl.Vector2, fontSize int) {
	font := getFont()
	textSize := rl.MeasureTextEx(*font, text, float32(fontSize), spacing)
	pos := rl.Vector2{X: position.X - textSize.X/2, Y: position.Y - textSize.Y/2}
	WriteText(text, pos, fontSize)
}

func WriteText(text string, position rl.Vector2, fontSize int) {
	font := getFont()
	rl.DrawTextEx(*font, text, position, float32(fontSize), spacing, rl.Black)
}
