package gameobjects

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"testing"
)

type MockGameObject struct {
	name   string
	alive  bool
	hitbox rl.Rectangle
}

func (m *MockGameObject) Update() error {
	return nil
}

func (m *MockGameObject) Draw() error {
	return nil
}

func (m *MockGameObject) IsAlive() bool {
	return m.alive
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
	collection.Update()
	collection.Update()

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
	collection.Update()

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
	collection.Update()

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
