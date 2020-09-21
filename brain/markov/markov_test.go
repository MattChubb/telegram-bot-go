package markov

import (
	"testing"
	"reflect"
	"github.com/mb-14/gomarkov"
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

func TestGenerateSentence(t *testing.T) {
	/*
	   Markov chains are inherently random, and so there's no real point in trying to test them deterministically.
	   Instead, we can feed it various prompts and check we get _something_ out the other side
	   TODO Check that what we get out is expected:
	   * Contains only the words "test" and "data" in that order
	   * Does not hit the tokensLengthLimit
	*/
	tables := []struct {
		testcase string
		input    []string
	}{
		{"Null", []string{}},
		//TODO Trim empty strings from input
		//{"Empty string", []string{""}},
		{"1 word", []string{"test"}},
		{"1 word 2", []string{"data"}},
		{"2 words", []string{"test", "data"}},
		{"3 words", []string{"test", "data", "test"}},
		{"Unknown word", []string{"testing"}},
	}

    brain := new(Brain)
    brain.Init(1, 32)

	brain.Train("test data test data test data")
	brain.Train("data test data test data")
	brain.Train("test data test data test data")

	for _, table := range tables {
		t.Logf("Testing: %s", table.testcase)
		got := brain.generateSentence(table.input)

		if len(got) < 1 {
			t.Errorf("prompt: %#v, got: %#v", table.input, got)
		} else if got[0] == gomarkov.StartToken {
			t.Errorf("Start token found, got: %#v", got)
		} else if got[len(got)-1] == gomarkov.EndToken {
			t.Errorf("End token found, got: %#v", got)
		} else {
			//t.Logf("Got: %#v", got)
			t.Logf("Passed (%d tokens returned)", len(got))
		}
	}
}

func TestTrimMessage(t *testing.T){
	tables := []struct {
		testcase string
		input    []string
		expected []string
	}{
		{"1 uncommon word", []string{"test"}, []string{"test"}},
		{"1 uncommon word, 1 common word", []string{"the", "test"}, []string{"test"}},
		{"Mention", []string{"test", "@self"}, []string{"test"}},
	}

	for _, table := range tables {
		t.Logf("Testing: %s", table.testcase)
		got := trimMessage(table.input)
		if !reflect.DeepEqual(got, table.expected) {
			t.Errorf("expected: %#v, got: %#v", table.expected, got)
		} else {
			t.Log("Passed")
		}
	}
}

func TestIsStopWord(t *testing.T){
	tables := []struct {
		testcase string
		input    string
		expected bool
	}{
		{"Non stopword", "test", false},
		{"Stopword", "the", true},
		{"Contains stopword", "theadore", false},
	}

	for _, table := range tables {
		t.Logf("Testing: %s", table.testcase)
		got := isStopWord(table.input)
		if !reflect.DeepEqual(got, table.expected) {
			t.Errorf("expected: %#v, got: %#v", table.expected, got)
		} else {
			t.Log("Passed")
		}
	}
}
