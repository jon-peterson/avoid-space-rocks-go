package playfield

import (
	"avoid_the_space_rocks/internal/gameobjects"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

type Spaceship struct {
	gameobjects.Transform
	SpriteSheet rl.Texture2D // The texture with the packed sprites
	frameWidth  float32      // Width of each frame
	frameHeight float32      // Height of each frame
	frameCount  int          // The number of frames in the sheet
}

func MakeSpaceship() Spaceship {
	ship := Spaceship{
		SpriteSheet: rl.LoadTexture("assets/sprites/spaceship.png"),
	}
	ship.frameWidth = float32(ship.SpriteSheet.Width)
	ship.frameHeight = float32(ship.SpriteSheet.Height / 3)
	ship.frameCount = 3
	// Traditionally starts pointing straight up
	ship.Rotation = rl.Vector2{
		X: 0.0,
		Y: -1.0,
	}
	return ship
}

// Draw the spaceship at its current position and rotation
func (s *Spaceship) Draw() {
	frame := s.frame(1)
	destination := rl.Rectangle{
		X:      s.Transform.Position.X,
		Y:      s.Transform.Position.Y,
		Width:  s.frameWidth,
		Height: s.frameHeight,
	}
	rotationDegrees := math.Atan2(float64(s.Rotation.Y), float64(s.Rotation.X)) * 180 / math.Pi
	rl.DrawTexturePro(s.SpriteSheet, frame, destination, rl.Vector2{}, float32(rotationDegrees), rl.Black)
}

// Returns the rectangle coordinates of the specified frame in the spreadsheet
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
