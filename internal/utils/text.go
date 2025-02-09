package utils

import rl "github.com/gen2brain/raylib-go/raylib"

// CenterText draws the given text centered around the passed-in position
func CenterText(text string, position rl.Vector2, fontSize int32) {
	textSize := rl.MeasureText(text, fontSize)
	rl.DrawText(text, int32(position.X)-textSize/2, int32(position.Y)-fontSize/2, fontSize, rl.Black)
}
