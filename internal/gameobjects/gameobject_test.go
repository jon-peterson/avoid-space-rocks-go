package gameobjects

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"testing"
)

type MockGameObject struct {
	name   string
	alive  bool
	enemy  bool
	hitbox rl.Rectangle
}

func (m *MockGameObject) Update(delta float32) error {
	return nil
}

func (m *MockGameObject) Draw() error {
	return nil
}

func (m *MockGameObject) IsAlive() bool {
	return m.alive
}

func (m *MockGameObject) IsEnemy() bool {
	return m.enemy
}

func (m *MockGameObject) OnCollision(_ Collidable) error {
	return nil
}

func (m *MockGameObject) GetHitbox() rl.Rectangle {
	return m.hitbox
}

func TestGameObjectCollectionUpdate(t *testing.T) {
	collection := NewGameObjectCollection()

	// Add some mock objects
	collection.Add(&MockGameObject{alive: true})
	collection.Add(&MockGameObject{alive: false})
	collection.Add(&MockGameObject{alive: true})
	collection.Add(&MockGameObject{alive: false})

	// Update the collection twice (so it can remove the dead objects after adding)
	collection.Update(0.1)
	collection.Update(0.1)

	// Verify that only alive objects remain
	if len(collection.objects) != 2 {
		t.Errorf("Expected 2 alive objects, got %d", len(collection.objects))
	}

	for _, obj := range collection.objects {
		if !obj.IsAlive() {
			t.Errorf("Found a dead object in the collection")
		}
	}
}

func TestGameObjectCollectionCollisionCheck(t *testing.T) {
	collection := NewGameObjectCollection()

	// Add some mock collidable objects
	collection.Add(&MockGameObject{alive: true, hitbox: rl.NewRectangle(0, 0, 10, 10)})
	collection.Add(&MockGameObject{alive: true, hitbox: rl.NewRectangle(5, 5, 10, 10)})
	collection.Add(&MockGameObject{alive: true, hitbox: rl.NewRectangle(20, 20, 10, 10)})
	collection.Update(0.1)

	// Verify that collisions are detected correctly
	// In this case, the first two objects should collide
	if !rl.CheckCollisionRecs(collection.objects[0].(Collidable).GetHitbox(), collection.objects[1].(Collidable).GetHitbox()) {
		t.Errorf("Expected collision between first two objects")
	}

	// The third object should not collide with the first two
	if rl.CheckCollisionRecs(collection.objects[0].(Collidable).GetHitbox(), collection.objects[2].(Collidable).GetHitbox()) {
		t.Errorf("Did not expect collision between first and third objects")
	}
	if rl.CheckCollisionRecs(collection.objects[1].(Collidable).GetHitbox(), collection.objects[2].(Collidable).GetHitbox()) {
		t.Errorf("Did not expect collision between second and third objects")
	}
}

func TestGameObjectCollectionAny(t *testing.T) {
	collection := NewGameObjectCollection()

	// Add some mock objects
	collection.Add(&MockGameObject{alive: true, name: "John Bonham"})
	collection.Add(&MockGameObject{alive: false, name: "Robert Plant"})
	collection.Add(&MockGameObject{alive: true, name: "Jimmy Page"})
	collection.Add(&MockGameObject{alive: false, name: "John Paul Jones"})
	collection.Update(0.1)

	// Test case: Check if any object is alive
	if !collection.Any(func(obj GameObject) bool {
		return obj.IsAlive()
	}) {
		t.Errorf("Expected at least one alive object")
	}

	// Test case: Check if any object is Jimmy Page
	if !collection.Any(func(obj GameObject) bool {
		// Cast obj to MockGameObject to access the name field
		if mockObj, ok := obj.(*MockGameObject); ok {
			return mockObj.name == "Jimmy Page"
		}
		return false
	}) {
		t.Errorf("Expected to find Jimmy Page")
	}

	if collection.Any(func(obj GameObject) bool {
		// Cast obj to MockGameObject to access the name field
		if mockObj, ok := obj.(*MockGameObject); ok {
			return mockObj.name == "Perry Como"
		}
		return false
	}) {
		t.Errorf("Should not have found Perry Como")
	}
}

func TestGameObjectCollectionForEach(t *testing.T) {
	collection := NewGameObjectCollection()

	// Add some mock objects
	collection.Add(&MockGameObject{alive: true})
	collection.Add(&MockGameObject{alive: false})
	collection.Add(&MockGameObject{alive: true})
	collection.Add(&MockGameObject{alive: false})
	collection.Update(0.1)

	collection.ForEach(func(obj GameObject) {
		if mockObj, ok := obj.(*MockGameObject); ok {
			mockObj.alive = false
		}
	})

	// Test case: Check if any object is alive
	if collection.Any(func(obj GameObject) bool {
		return obj.IsAlive()
	}) {
		t.Errorf("Expected all dead objects")
	}
}

func TestGameObjectCollectionHasRemainingEnemies(t *testing.T) {
	collection := NewGameObjectCollection()

	// Add some mock objects
	collection.Add(&MockGameObject{alive: true, enemy: true})
	collection.Add(&MockGameObject{alive: false, enemy: true})
	collection.Add(&MockGameObject{alive: true, enemy: false})
	collection.Add(&MockGameObject{alive: false, enemy: false})
	collection.Update(0.1)

	// Test case: Check if there are any remaining enemies
	if !collection.HasRemainingEnemies() {
		t.Errorf("Expected to find remaining enemies")
	}

	// Remove all enemies
	collection.ForEach(func(obj GameObject) {
		if mockObj, ok := obj.(*MockGameObject); ok && mockObj.enemy {
			mockObj.alive = false
		}
	})
	collection.Update(0.1)
	if collection.HasRemainingEnemies() {
		t.Errorf("Did not expect to find remaining enemies after marking all dead")
	}

	// Add a new non-enemy object, still should have no enemies
	collection.Add(&MockGameObject{alive: true, enemy: false})
	if collection.HasRemainingEnemies() {
		t.Errorf("Newly added non-enemy should leave collection no remaining enemies")
	}

	// Add a new enemy object, now it should say true again
	collection.Add(&MockGameObject{alive: true, enemy: true})
	if !collection.HasRemainingEnemies() {
		t.Errorf("Newly added enemy should leave collection with remaining enemies")
	}
}
