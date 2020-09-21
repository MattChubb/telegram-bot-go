package markov

import (
	"testing"
	"reflect"
)

func TestInit(t *testing.T) {
	tables := []struct {
		testcase string
		order    int
        length   int
	}{
        {"Chain of order 1", 1, 32},
        {"Chain of order 2", 2, 32},
        {"Chain of order 100", 100, 32}, //Don't try this at home!
        {"Chain of order 0", 0, 32},
        {"Chain of order -1", -1, 32},
    }

    brain := new(Brain)

	for _, table := range tables {
		t.Logf("Testing: %s", table.testcase)
		brain.Init(table.order, table.length)
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
    brain.Init(1, 32)

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

func TestGenerate(t *testing.T) {
	tables := []struct {
		testcase string
		input    string
	}{
		//TODO Trim empty strings from input
		//{"Empty string", []string{""}},
		{"1 word", "test"},
		{"1 word 2", "data"},
		{"2 words", "test data"},
		{"3 words", "test data test"},
		{"Unknown word", "testing"},
	}

    brain := new(Brain)
    brain.Init(1, 32)

	brain.Train("test data test data test data")
	brain.Train("data test data test data")
	brain.Train("test data test data test data")

	for _, table := range tables {
		t.Logf("Testing: %s", table.testcase)
        //TODO Test error handling
		got, _ := brain.Generate(table.input)

		if len(got) < 1 {
			t.Errorf("prompt: %#v, got: %#v", table.input, got)
		} else {
			//t.Logf("Got: %#v", got)
			t.Logf("Passed (%d characters returned)", len(got))
		}

        if got[0] == 'T' || got[0] == 'D' {
            t.Logf("Passed (First letter %q capitalised)", got[0])
        } else {
            t.Errorf("First letter %q not capitalised", got[0])
        }
	}
}
