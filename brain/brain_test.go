package brain


import (
	"testing"
	"reflect"
)

func TestProcessString(t *testing.T) {
	tables := []struct {
		testcase string
		input    string
		expected []string
	}{
		{"1 word", "test", []string{"test"}},
		{"2 words", "test data", []string{"test", " ", "data"}},
		{"3 words", "test data one", []string{"test", " ", "data", " ", "one"}},
		{"0 words", "", []string{""}},
		{"alphanumeric", "test1data", []string{"test1data"}},
		{"punctuation", "test. data,", []string{"test", ". ", "data", ","}},
		{"Capitalisation", "Test. data,", []string{"test", ". ", "data", ","}},
		{"Mixed Case", "Test. Data,", []string{"test", ". ", "data", ","}},
		{"AlTeRnAtInG CaSe", "TeSt. DaTa,", []string{"test", ". ", "data", ","}},
	}

	for _, table := range tables {
		t.Logf("Testing: %s", table.testcase)
		got := ProcessString(table.input)
		if !reflect.DeepEqual(got, table.expected) {
			t.Errorf("expected: %#v, got: %#v", table.expected, got)
		} else {
			t.Log("Passed")
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
