package gameobjects

import "testing"

type MockGameObject struct {
	alive bool
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

func TestGameObjectCollectionUpdate(t *testing.T) {
	collection := NewGameObjectCollection()

	// Add some mock objects
	collection.Add(&MockGameObject{alive: true})
	collection.Add(&MockGameObject{alive: false})
	collection.Add(&MockGameObject{alive: true})
	collection.Add(&MockGameObject{alive: false})

	// Update the collection
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
