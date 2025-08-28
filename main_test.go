package main

import "testing"

func TestHello(t *testing.T) {
	want := "Hello, World!"
	if got := Hello(); got != want {
		t.Fatalf("Hello() = %q; want %q", got, want)
	}
}
