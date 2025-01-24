package playfield

import (
	"avoid_the_space_rocks/internal/gameobjects"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Spaceship struct {
	gameobjects.Transform
	SpriteSheet rl.Texture2D
	frameWidth  float32
	frameHeight float32
	frameCount  int
}

func MakeSpaceship() Spaceship {
	ship := Spaceship{
		SpriteSheet: rl.LoadTexture("assets/sprites/spaceship.png"),
	}
	ship.frameWidth = float32(ship.SpriteSheet.Width)
	ship.frameHeight = float32(ship.SpriteSheet.Height / 3)
	ship.frameCount = 3
	return ship
}

func (s *Spaceship) String() string {
	return fmt.Sprintf("Yo: %z", s.Transform)
}

func (s *Spaceship) Draw() {
	frame := s.frame(1)
	destination := rl.Rectangle{
		X:      s.Transform.Position.X,
		Y:      s.Transform.Position.Y,
		Width:  s.frameWidth,
		Height: s.frameHeight,
	}
	rl.DrawTexturePro(s.SpriteSheet, frame, destination, rl.Vector2{}, 0, rl.Black)
}

func (s *Spaceship) frame(n int) rl.Rectangle {
	if n < 1 || n > s.frameCount {
		// Is this the right thing to do here?
		panic("invalid frame number")
	}
	return rl.Rectangle{
		X:      0,
		Y:      0,
		Width:  s.frameWidth,
		Height: float32(n) * s.frameHeight,
	}
}
