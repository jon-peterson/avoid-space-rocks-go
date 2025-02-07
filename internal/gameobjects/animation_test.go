package gameobjects

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"testing"
)

func TestSpriteSheet_frame(t *testing.T) {
	sheet, err := NewSpriteSheet("test.png", 2, 2)
	if err != nil {
		t.Fatalf("Failed to create SpriteSheet: %v", err)
	}

	tests := []struct {
		row, col int
		expected rl.Rectangle
	}{
		{0, 0, rl.Rectangle{X: 0, Y: 0, Width: float32(sheet.frameWidth), Height: float32(sheet.frameHeight)}},
		{0, 1, rl.Rectangle{X: float32(sheet.frameWidth), Y: 0, Width: float32(sheet.frameWidth), Height: float32(sheet.frameHeight)}},
		{1, 0, rl.Rectangle{X: 0, Y: float32(sheet.frameHeight), Width: float32(sheet.frameWidth), Height: float32(sheet.frameHeight)}},
		{1, 1, rl.Rectangle{X: float32(sheet.frameWidth), Y: float32(sheet.frameHeight), Width: float32(sheet.frameWidth), Height: float32(sheet.frameHeight)}},
	}

	for _, tt := range tests {
		frame, err := sheet.frame(tt.row, tt.col)
		if err != nil {
			t.Errorf("Unexpected error for frame (%d, %d): %v", tt.row, tt.col, err)
		}
		if frame != tt.expected {
			t.Errorf("Expected frame %v, got %v", tt.expected, frame)
		}
	}

	// Test out of bounds
	_, err = sheet.frame(-1, 0)
	if err == nil {
		t.Error("Expected error for out of bounds frame (-1, 0), got nil")
	}

	_, err = sheet.frame(0, -1)
	if err == nil {
		t.Error("Expected error for out of bounds frame (0, -1), got nil")
	}

	_, err = sheet.frame(2, 0)
	if err == nil {
		t.Error("Expected error for out of bounds frame (2, 0), got nil")
	}

	_, err = sheet.frame(0, 2)
	if err == nil {
		t.Error("Expected error for out of bounds frame (0, 2), got nil")
	}
}

func TestSpriteSheet_GetRectangle(t *testing.T) {
	sheet, err := NewSpriteSheet("test.png", 2, 2)
	if err != nil {
		t.Fatalf("Failed to create SpriteSheet: %v", err)
	}

	center := rl.NewVector2(50, 50)
	expected := rl.Rectangle{
		X:      center.X - float32(sheet.frameWidth)/2,
		Y:      center.Y - float32(sheet.frameHeight)/2,
		Width:  float32(sheet.frameWidth),
		Height: float32(sheet.frameHeight),
	}

	rect := sheet.GetRectangle(center)
	if rect != expected {
		t.Errorf("Expected rectangle %v, got %v", expected, rect)
	}
}
