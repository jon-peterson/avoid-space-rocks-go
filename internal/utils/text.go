package utils

import rl "github.com/gen2brain/raylib-go/raylib"

// CenterText draws the given text centered around the passed-in position
func CenterText(text string, position rl.Vector2, fontSize int32) {
	textSize := rl.MeasureText(text, fontSize)
	pos := rl.Vector2{X: position.X - float32(textSize)/2, Y: position.Y - float32(fontSize)/2}
	WriteText(text, pos, fontSize)
}

func WriteText(text string, position rl.Vector2, fontSize int32) {
	rl.DrawText(text, int32(position.X), int32(position.Y), fontSize, rl.Black)
}
