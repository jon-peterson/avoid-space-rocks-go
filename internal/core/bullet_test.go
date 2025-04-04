package core

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"testing"
)

func TestBullet_IsAlive(t *testing.T) {
	bullet := NewBullet(rl.NewVector2(0, 0), rl.NewVector2(0, 0), true)

	// Test that the bullet is alive immediately after creation
	if !bullet.IsAlive() {
		t.Errorf("Expected bullet to be alive immediately after creation")
	}

	// Test that the bullet is not alive after its lifetime has passed
	bullet.ageMs = bulletLifetimeMs + 1
	if bullet.IsAlive() {
		t.Errorf("Expected bullet to be dead after its lifetime has passed")
	}
}

func TestBullet_GetHitbox(t *testing.T) {
	position := rl.NewVector2(10, 20)
	bullet := NewBullet(position, rl.NewVector2(0, 0), true)

	expectedHitbox := rl.Rectangle{
		X:      position.X,
		Y:      position.Y,
		Width:  1,
		Height: 1,
	}

	hitbox := bullet.GetHitbox()
	if hitbox != expectedHitbox {
		t.Errorf("Expected hitbox %v, got %v", expectedHitbox, hitbox)
	}
}

func TestBullet_OnCollisionWithRock(t *testing.T) {
	bullet := NewBullet(rl.NewVector2(0, 0), rl.NewVector2(0, 0), true)
	rock := NewRock(RockBig, rl.NewVector2(0, 0))
	err := bullet.OnCollision(&rock)
	if err != nil {
		t.Errorf("Unexpected error during collision: %v", err)
	}

	// Check if the rock is destroyed after collision
	if rock.IsAlive() {
		t.Errorf("Expected rock to be destroyed after collision")
	}

	// Check if the bullet is still alive after collision
	if bullet.IsAlive() {
		t.Errorf("Expected bullet to be destroyed after collision")
	}
}

func TestBullet_OnCollisionWithAlien(t *testing.T) {
	bullet := NewBullet(rl.NewVector2(0, 0), rl.NewVector2(0, 0), true)
	alien := NewAlien(AlienBig, rl.NewVector2(0, 0))
	err := bullet.OnCollision(&alien)
	if err != nil {
		t.Errorf("Unexpected error during collision: %v", err)
	}

	// Check if the alien is destroyed after collision
	if alien.IsAlive() {
		t.Errorf("Expected alien to be destroyed after collision")
	}

	// Check if the bullet is still alive after collision
	if bullet.IsAlive() {
		t.Errorf("Expected bullet to be destroyed after collision")
	}
}

func TestBullet_OnCollisionWithSpaceshipFiredBySpaceship(t *testing.T) {
	bullet := NewBullet(rl.NewVector2(0, 0), rl.NewVector2(0, 0), true)
	ship := NewSpaceship()
	ship.Alive = true
	err := bullet.OnCollision(&ship)
	if err != nil {
		t.Errorf("Unexpected error during collision: %v", err)
	}

	if !ship.IsAlive() {
		t.Errorf("Expected ship to be alive after collision with own bullet")
	}

	if !bullet.IsAlive() {
		t.Errorf("Expected bullet to be alive after collision with ship")
	}
}

func TestBullet_OnCollisionWithSpaceshipFiredByAlien(t *testing.T) {
	bullet := NewBullet(rl.NewVector2(0, 0), rl.NewVector2(0, 0), false)
	ship := NewSpaceship()
	err := bullet.OnCollision(&ship)
	if err != nil {
		t.Errorf("Unexpected error during collision: %v", err)
	}

	if ship.IsAlive() {
		t.Errorf("Expected ship to be dead after collision with alien bullet")
	}

	if bullet.IsAlive() {
		t.Errorf("Expected alien bullet to be dead after collision with ship")
	}
}
