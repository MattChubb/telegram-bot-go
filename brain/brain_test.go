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

func TestExtractSubject(t *testing.T) {
	tables := []struct {
		testcase string
		input    []string
        length   int
		expected []string
	}{
		{"1 uncommon word", []string{"test"}, 1, []string{"test"}},
		{"Empty string", []string{""}, 1, []string{}},
		{"Empty string, order 2", []string{""}, 2, []string{}},
		{"Empty sentence", []string{}, 1, []string{}},
		{"Empty sentence, order 2", []string{}, 2, []string{}},
		{"1 uncommon word, 1 common word", []string{"the", "test"}, 1, []string{"test"}},
		{"Punctuation", []string{".", ".", "."}, 1, []string{}},
		{"Sentence", []string{"the", " ", "test"}, 1, []string{"test"}},
		{"Sentence, all stopwords", []string{"the", "the", "the"}, 1, []string{}},
		{"Mention", []string{"test", "@self"}, 1, []string{"test"}},
		{"2 word subject", []string{"the", " ", "test", " ", "and"}, 2, []string{" ", "test"}},
		{"2 word subject, end of sentence", []string{"the", " ", "test"}, 2, []string{" ", "test"}},
		{"2 word subject, beginning of sentence", []string{"test", " ", "the", " "}, 2, []string{"test", " "}},
	}

	for _, table := range tables {
		t.Logf("Testing: %s", table.testcase)
		got := ExtractSubject(table.input, table.length)
		if !reflect.DeepEqual(got, table.expected) {
			t.Errorf("FAIL, expected: %#v, got: %#v", table.expected, got)
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
			t.Errorf("FAIL, expected: %#v, got: %#v", table.expected, got)
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
			t.Errorf("FAIL, expected: %#v, got: %#v", table.expected, got)
		} else {
			t.Log("Passed")
		}
	}
}
