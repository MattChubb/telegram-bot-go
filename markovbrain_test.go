package markovbrain

import (
	"testing"
)

func TestInit(t *testing.T) {
	tables := []struct {
		testcase string
		input    int
	}{
        {"Chain of order 1", 1},
        {"Chain of order 2", 2},
        {"Chain of order 100", 100}, //Don't try this at home!
        {"Chain of order 0", 0},
        {"Chain of order -1", -1},
    }

    brain := new(Brain)

	for _, table := range tables {
		t.Logf("Testing: %s", table.testcase)
		brain.Init(table.input)
        t.Log("Initialised without crashing")
    }
}
