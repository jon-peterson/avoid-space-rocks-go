package gameobjects

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Body struct {
	Transform
	acceleration rl.Vector2
}

func (b *Body) String() string {
	return fmt.Sprintf("%v acc (%f,%f)", b.Transform, b.acceleration.X, b.acceleration.Y)
}
