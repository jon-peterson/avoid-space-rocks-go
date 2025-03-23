package core

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"testing"
)

func TestWraparound(t *testing.T) {
	world := NewWorld(800, 600)

	tests := []struct {
		input    rl.Vector2
		expected rl.Vector2
	}{
		{input: rl.Vector2{X: -10, Y: 300}, expected: rl.Vector2{X: 800, Y: 300}},
		{input: rl.Vector2{X: 810, Y: 300}, expected: rl.Vector2{X: 0, Y: 300}},
		{input: rl.Vector2{X: 400, Y: -10}, expected: rl.Vector2{X: 400, Y: 600}},
		{input: rl.Vector2{X: 400, Y: 610}, expected: rl.Vector2{X: 400, Y: 0}},
	}

	for _, test := range tests {
		result := world.Wraparound(test.input)
		if result != test.expected {
			t.Errorf("Wraparound(%v) = %v, want %v", test.input, result, test.expected)
		}
	}
}

func TestIsOutsideEdges(t *testing.T) {
	world := NewWorld(800, 600)

	tests := []struct {
		input    rl.Vector2
		expected bool
	}{
		{input: rl.Vector2{X: -10, Y: 300}, expected: true},
		{input: rl.Vector2{X: 810, Y: 300}, expected: true},
		{input: rl.Vector2{X: 400, Y: -10}, expected: true},
		{input: rl.Vector2{X: 400, Y: 610}, expected: true},
		{input: rl.Vector2{X: 400, Y: 300}, expected: false},
	}

	for _, test := range tests {
		result := world.IsOutsideEdges(test.input)
		if result != test.expected {
			t.Errorf("IsOutsideEdges(%v) = %v, want %v", test.input, result, test.expected)
		}
	}
}
