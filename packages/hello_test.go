package main

import (
    "strings"
    "testing"
)

const HelloWorld = "Hello World"

func TestHello(t *testing.T) {
    if (!strings.Contains(HelloWorld, "Hello")) {
        t.Error()
    }
}

func TestWorld(t *testing.T) {
    if (!strings.Contains(HelloWorld, "World")) {
        t.Error()
    }
}
