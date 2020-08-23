package main

import (
    "reflect"
    "testing"
)

func TestProcessString(t *testing.T) {
    tables := []struct{
        input string
        expected []string
    }{
        {"test", []string{"test"}},
        {"test data", []string{"test", "data"}},
        {"test data one", []string{"test", "data", "one"}},
        {"", []string{""}},
        {"test1data", []string{"test1data"}},
    }

    for _, table := range tables {
        got := processString(table.input)
        if ! reflect.DeepEqual(got, table.expected) {
            t.Errorf("expected: %#v, got: %#v", table.expected, got)
        }
    }
}
