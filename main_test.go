package main

import (
    "github.com/mb-14/gomarkov"
    "reflect"
    "testing"
)

func TestProcessString(t *testing.T) {
    tables := []struct{
        testcase string
        input string
        expected []string
    }{
        {"1 word", "test", []string{"test"}},
        {"2 words", "test data", []string{"test", "data"}},
        {"3 words", "test data one", []string{"test", "data", "one"}},
        {"0 words", "", []string{""}},
        {"alphanumeric", "test1data", []string{"test1data"}},
        {"punctuation", "test. data,", []string{"test.", "data,"}},
    }

    for _, table := range tables {
        t.Logf("Testing: %s", table.testcase)
        got := processString(table.input)
        if ! reflect.DeepEqual(got, table.expected) {
            t.Errorf("expected: %#v, got: %#v", table.expected, got)
        } else {
            t.Log("Passed")
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
    tables := []struct{
        testcase string
        input []string
    }{
        {"Null", []string{}},
        {"Empty string", []string{""}},
        {"1 word", []string{"test"}},
        {"1 word 2", []string{"data"}},
        {"2 words", []string{"test", "data"}},
        {"3 words", []string{"test", "data", "test"}},
        {"Unknown word", []string{"testing"}},
    }

    chain := gomarkov.NewChain(1)

    chain.Add(processString("test data test data test data"))
    chain.Add(processString("data test data test data"))
    chain.Add(processString("test data test data test data"))

    for _, table := range tables {
        t.Logf("Testing: %s", table.testcase)
        got := generateSentence(chain, table.input)

        if len(got) < 2 {
            t.Errorf("prompt: %#v, got: %#v", table.input, got)
        } else {
            //t.Logf("Got: %#v", got)
            t.Logf("Passed (%d tokens returned)", len(got))
        }
    }
}
