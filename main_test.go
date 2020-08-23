package main

import (
    "reflect"
    "testing"
)

func TestProcessString(t *testing.T) {
    input := "test data"
    expected := []string{"test", "data"}

    got := processString(input)
    if ! reflect.DeepEqual(got, expected) {
        t.Errorf("expected: %#v, got: %#v", expected, got)
    }
}
