package markov

import (
	"testing"
	"reflect"
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

func TestProcessString(t *testing.T) {
	tables := []struct {
		testcase string
		input    string
		expected []string
	}{
		{"1 word", "test", []string{"test"}},
		{"2 words", "test data", []string{"test", "data"}},
		{"3 words", "test data one", []string{"test", "data", "one"}},
		{"0 words", "", []string{""}},
		{"alphanumeric", "test1data", []string{"test1data"}},
		{"punctuation", "test. data,", []string{"test.", "data,"}},
		{"Capitalisation", "Test. data,", []string{"test.", "data,"}},
		{"Mixed Case", "Test. Data,", []string{"test.", "data,"}},
		{"AlTeRnAtInG CaSe", "TeSt. DaTa,", []string{"test.", "data,"}},
	}

	for _, table := range tables {
		t.Logf("Testing: %s", table.testcase)
		got := processString(table.input)
		if !reflect.DeepEqual(got, table.expected) {
			t.Errorf("expected: %#v, got: %#v", table.expected, got)
		} else {
			t.Log("Passed")
		}
	}
}

func TestTrain(t *testing.T) {
	tables := []struct {
		testcase string
		input    string
        errors   bool
	}{
        {"One word", "word", false},
        {"Two words", "two words", false},
        {"Two words with punctuation ", "two, words", false},
        {"Word and number", "1 one", false},
        {"Empty string", "", false},
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
