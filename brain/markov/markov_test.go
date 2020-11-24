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

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		order   int
		data    []string
		want    string
		wantErr bool
	}{
		{"Empty chain", 2, []string{}, `{"Chain":{"int":2,"spool_map":{},"freq_mat":{}},"LengthLimit":31}`, false},
		{"Empty chain, order 1", 1, []string{}, `{"Chain":{"int":1,"spool_map":{},"freq_mat":{}},"LengthLimit":31}`, false},
		{"Trained once", 1, []string{"test"}, `{"Chain":{"int":1,"spool_map":{"$":0,"^":2,"test":1},"freq_mat":{"0":{"1":1},"1":{"2":1}}},"LengthLimit":31}`, false},
		{"Trained on more data", 1, []string{"test data", "test data", "test node"}, `{"Chain":{"int":1,"spool_map":{"$":0,"^":3,"data":2,"node":4,"test":1},"freq_mat":{"0":{"1":3},"1":{"2":2,"4":1},"2":{"3":2},"4":{"3":1}}},"LengthLimit":31}`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            brain := new(Brain)
            brain.Init(tt.order, 31)
			for _, data := range tt.data {
				brain.Train(data)
			}

			got, err := brain.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Brain.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("Brain.MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		args    []byte
		wantErr bool
	}{
		{"Empty chain", []byte(`{"Chain":{"int":1,"spool_map":{},"freq_mat":{}},"LengthLimit":31}`), false},
		{"More complex chain", []byte(`{"Chain":{"int":1,"spool_map":{"$":0,"^":3,"data":2,"node":4,"test":1},"freq_mat":{"0":{"1":3},"1":{"2":2,"4":1},"2":{"3":2},"4":{"3":1}}},"LengthLimit":31}`), false},
		{"Invalid json", []byte(`{{"int":2,"spool_map":{},"freq_mat":{}}`), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            brain := new(Brain)

			if err := brain.UnmarshalJSON(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("Brain.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
            } else {
                t.Log("Successfully unmarshalled json")
			}

            if !tt.wantErr {
                //An error unmarshalling means we don't have a brain to train or generate from
                if err := brain.Train("test"); (err != nil) != tt.wantErr {
                    t.Errorf("Brain.Train() error = %v, wantErr %v", err, tt.wantErr)
                } else {
                    t.Log("Successfully trained unmarshalled brain")
                }

                if _, err := brain.Generate("test"); (err != nil) != tt.wantErr {
                    t.Errorf("Brain.Generate() error = %v, wantErr %v", err, tt.wantErr)
                } else {
                    t.Log("Successfully generated using trained unmarshalled brain")
                }
            }
		})
	}
}
