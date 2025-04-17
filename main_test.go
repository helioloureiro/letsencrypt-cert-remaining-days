package main

import (
	"testing"
	"crypto/x509"
	"encoding/pem"
)


func TestGrep(t *testing.T) {
	result := grep("helio", "helio loureiro")
	if result != true {
		t.Errorf("grep failed to find the parttern")
	}

	result = grep("something", "helio loureiro")
	if result != false {
		t.Errorf("grep failed and found a pattern")
	}
}


func TestSed(t *testing.T) {
	text := "helio loureiro"
	expected := "helioloureiro"
	response := sed(" ", "", text)
	if response != expected {
		t.Errorf("Expected: %s - got: %s", expected, response)
	}

	expected = "HeLio Loureiro"
	response = sed("h", "H", text)
	response = sed("l", "L", response)
	if response != expected {
		t.Errorf("Expected: %s - got: %s", expected, response)
	}
}

