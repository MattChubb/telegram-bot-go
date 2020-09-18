package markov

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

func TestTrain(t *testing.T) {
	tables := []struct {
		testcase string
		input    string
        errors   bool
	}{
        {"Train on a word", "word", false},
    }

    brain := new(Brain)
    brain.Init(1)

	for _, table := range tables {
		t.Logf("Testing: %s", table.testcase)
		err := brain.Train(table.input)

        if !table.errors && err != nil {
            t.Errorf("Expected no errors, got %#v", err)
        } else if table.errors && err == nil {
            t.Errorf("Expected errors, but got none")
        } else {
            t.Log("Initialised without crashing")
        }
    }
}
