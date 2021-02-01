package doublemarkov

import (
	log "github.com/sirupsen/logrus"
	"testing"
	"reflect"
    "regexp"
	"github.com/mb-14/gomarkov"
)

func TestMain(m *testing.M) {
    //log.SetLevel(log.DebugLevel)
    log.SetLevel(log.InfoLevel)
    m.Run()
}

func newBrain(order int, length int) *Brain {
    brain := new(Brain)
    brain.Init(order, length)
    brain.Train("test data test data")
    brain.Train("data test data")
    brain.Train("test data")
    return brain
}

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

	for _, table := range tables {
		t.Logf("Testing: %s", table.testcase)
        brain := new(Brain)
		brain.Init(table.order, table.length)
        t.Log("Initialised without crashing")
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

	for _, table := range tables {
		t.Logf("Testing: %s", table.testcase)
        brain := new(Brain)
        brain.Init(1, 32)

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
        order    int
        expected string
	}{
		{"Empty string, order 1", "", 1, `^[(Test)|(Data)][( test)|( data)]*$`},
		{"Empty string, order 2", "", 2, `^[(Test)|(Data)][( test)|( data)]*$`},
		{"1 word, order 1", "test", 1, `^[(Test)|(Data)][( test)|( data)]*$`},
		{"1 word, order 2", "test", 2, `^[(Test)|(Data)][( test)|( data)]*$`},
		{"1 word, order 2, no double space", "test", 2, `[^\s{2}]`},
		{"1 word 2", "data", 1, `^[(Test)|(Data)][( test)|( data)]*$`},
		{"2 words", "test data", 1, `^[(Test)|(Data)][( test)|( data)]*$`},
		{"3 words", "test data test", 1, `^[(Test)|(Data)][( test)|( data)]*$`},
		{"Unknown word", "testing", 1, `^Testing$`},
		{"Unknown word, order 2", "testing", 2, `^Testing$`},
	}

    const length = 32
	for _, table := range tables {
        t.Logf("Testing: %s", table.testcase)
        brain := newBrain(table.order, length)
        brain.Train("test data test data test data")
        brain.Train("data test data test data")
        brain.Train("test data test data test data")

        //TODO Test error handling
        got, _ := brain.Generate(table.input)

        if len(got) < 1 {
            t.Errorf("prompt: %#v, got: %#v", table.input, got)
        } else if len(got) > length * 5 {
            t.Errorf("Response largr than lengthlimit, got: %#v", got)
        } else {
            //t.Logf("Got: %#v", got)
            t.Logf("Passed (%d characters returned)", len(got))
        }

        if got[0] == 'T' || got[0] == 'D' {
            t.Logf("Passed (First letter %q capitalised)", got[0])
        } else {
            t.Errorf("First letter %q not capitalised", got[0])
        }

        if match, _ := regexp.Match(table.expected, []byte(got)); ! match {
            t.Errorf("Output not as expected, got: %#v", got)
        }
	}

    brain := new(Brain)
    brain.Init(1, length)
    brain.Train("test data test data test data")
    brain.Train("data test data test data")
    brain.Train("test data test data test data")
    brain.Train("test subject data")
    got, _ := brain.Generate("subject")
    if got[0:7] == "Subject" {
        t.Errorf("Nothing generated before subject, got: %#v", got)
    } else if got[len(got)-7:len(got)] == "subject" {
        t.Errorf("Nothing generated after subject, got: %#v", got)
    } else if match, _ := regexp.Match(`ubject subject`, []byte(got)); match {
        t.Errorf("Subject generated twice, got: %#v", got)
    } else if match, _ := regexp.Match(`ubjectsubject`, []byte(got)); match {
        t.Errorf("Subject generated twice, got: %#v", got)
    } else if match, _ := regexp.Match(`\s{2}`, []byte(got)); match {
        t.Errorf("Double spacing, got: %#v", got)
    } else {
        t.Logf("Passed (text generated either side of subject): %#v", got[len(got)-7:len(got)])
    }
}

func TestGenerateSentence(t *testing.T) {
	tables := []struct {
		testcase string
		input    []string
        order    int
        expected string
	}{
		{"Null, order 1", []string{}, 1, `((test)|(data)|\W)`},
		{"Null, order 2", []string{}, 2, `((test)|(data)|\W)`},
		{"Empty string, order 1", []string{""}, 1, `^$`},
		{"Empty string, order 2", []string{""}, 2, `^$`},
		{"1 word", []string{"test"}, 1, `((test)|(data)|\W)`},
		{"1 word 2", []string{"data"}, 2, `((test)|(data)|\W)`},
		{"2 words, order 1", []string{"test", "data"}, 1, `((test)|(data)|\W)`},
		{"2 words, order 2", []string{"test", "data"}, 2, `((test)|(data)|\W)`},
		{"3 words", []string{"test", "data", "test"}, 2, `((test)|(data)|\W)`},
		{"Unknown word", []string{"testing"}, 2, `testing`},
	}

    const length = 32
	for _, table := range tables {
		t.Logf("Testing: %s", table.testcase)
        brain := newBrain(table.order, length)

		got := brain.generateSentence(brain.fwdChain, table.input)

		if len(got) < 1 {
			t.Errorf("prompt: %#v, got: %#v", table.input, got)
		} else if len(got) > length/2 {
			t.Errorf("Response largr than lengthlimit, got: %#v", got)
		} else if got[0] == gomarkov.StartToken {
			t.Errorf("Start token found, got: %#v", got)
		} else if got[len(got)-1] == gomarkov.EndToken {
			t.Errorf("End token found, got: %#v", got)
		}

        for _, word := range got {
            if match, _ := regexp.Match(table.expected, []byte(word)); ! match {
                t.Errorf("Output not as expected, got: %#v", got)
            }
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
		{"Empty chain", 2, []string{}, `{"BckChain":{"int":2,"spool_map":{},"freq_mat":{}},"FwdChain":{"int":2,"spool_map":{},"freq_mat":{}},"LengthLimit":31}`, false},
		{"Empty chain, order 1", 1, []string{}, `{"BckChain":{"int":1,"spool_map":{},"freq_mat":{}},"FwdChain":{"int":1,"spool_map":{},"freq_mat":{}},"LengthLimit":31}`, false},
		{"Trained once", 1, []string{"test"}, `{"BckChain":{"int":1,"spool_map":{"$":0,"^":2,"test":1},"freq_mat":{"0":{"1":1},"1":{"2":1}}},"FwdChain":{"int":1,"spool_map":{"$":0,"^":2,"test":1},"freq_mat":{"0":{"1":1},"1":{"2":1}}},"LengthLimit":31}`, false},
		{"Trained on more data", 1, []string{"test data", "test data", "test node"}, `{"BckChain":{"int":1,"spool_map":{" ":2,"$":0,"^":4,"data":1,"node":5,"test":3},"freq_mat":{"0":{"1":2,"5":1},"1":{"2":2},"2":{"3":3},"3":{"4":3},"5":{"2":1}}},"FwdChain":{"int":1,"spool_map":{" ":2,"$":0,"^":4,"data":3,"node":5,"test":1},"freq_mat":{"0":{"1":3},"1":{"2":3},"2":{"3":2,"5":1},"3":{"4":2},"5":{"4":1}}},"LengthLimit":31}`, false},
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
		{"Empty chain", []byte(`{"BckChain":{"int":1,"spool_map":{},"freq_mat":{}},"FwdChain":{"int":1,"spool_map":{},"freq_mat":{}},"LengthLimit":31}`), false},
		{"More complex chain", []byte(`{"BckChain":{"int":1,"spool_map":{"$":0,"^":3,"data":2,"node":4,"test":1},"freq_mat":{"0":{"1":3},"1":{"2":2,"4":1},"2":{"3":2},"4":{"3":1}}},"FwdChain":{"int":1,"spool_map":{"$":0,"^":3,"data":2,"node":4,"test":1},"freq_mat":{"0":{"1":3},"1":{"2":2,"4":1},"2":{"3":2},"4":{"3":1}}},"LengthLimit":31}`), false},
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
