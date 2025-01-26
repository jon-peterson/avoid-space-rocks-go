package playfield

import (
	"avoid_the_space_rocks/internal/gameobjects"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

type Spaceship struct {
	gameobjects.Rigidbody
	SpriteSheet rl.Texture2D // The texture with the packed sprites
	frameWidth  float32      // Width of each frame
	frameHeight float32      // Height of each frame
	frameCount  int          // The number of frames in the sheet
	FuelBurning bool         // Is the user burning fuel to accelerate?
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

// Update the status of the spaceship given the current state of the game
func (s *Spaceship) Update() {
	if s.FuelBurning {
		s.Velocity = rl.Vector2Add(s.Velocity, rl.Vector2Scale(s.Rotation, 0.01))
	} else {
		s.Velocity = rl.Vector2{}
	}
	s.Position = rl.Vector2Add(s.Position, s.Velocity)
}

// Draw the spaceship at its current position and rotation
func (s *Spaceship) Draw() {
	frame := s.frame(s.frameIndex())
	destination := rl.Rectangle{
		X:      s.Transform.Position.X,
		Y:      s.Transform.Position.Y,
		Width:  s.frameWidth,
		Height: s.frameHeight,
	}
	// The origin is the center of the sprite, so rotation works around that center
	origin := rl.Vector2{
		X: s.frameWidth / 2,
		Y: s.frameHeight / 2,
	}
	rotationDegrees := math.Atan2(float64(s.Rotation.Y), float64(s.Rotation.X)) * 180 / math.Pi
	rl.DrawTexturePro(s.SpriteSheet, frame, destination, origin, float32(rotationDegrees), rl.Black)
}

// frameIndex returns the index of the correct frame to use in the sprite sheet. There are two
// fuel burning frames, so the index is either 0, 1, or 2.
func (s *Spaceship) frameIndex() int {
	if s.FuelBurning {
		t := rl.GetTime()
		if t-math.Floor(t) < 0.5 {
			return 1
		} else {
			return 2
		}
	}
	return 0
}

// Returns the rectangle coordinates of the specified frame in the spreadsheet
func (s *Spaceship) frame(n int) rl.Rectangle {
	if n < 0 || n > s.frameCount-1 {
		// Is this the right thing to do here?
		panic("invalid frame number")
	}
	return rl.Rectangle{
		X:      0,
		Y:      float32(n) * s.frameHeight,
		Width:  s.frameWidth,
		Height: s.frameHeight,
	}
}
