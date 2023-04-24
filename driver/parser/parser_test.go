package parser

import (
	"strings"
	"testing"
)

func TestParseEmpty(t *testing.T) {
	_, err := ParseBytes([]byte{})
	if err == nil {
		t.Fatal("parsing of empty input should have failed")
	}
}

var minimalExample = `
name: Minimal Example
`

func TestParseMinimalExample(t *testing.T) {
	_, err := ParseBytes([]byte(minimalExample))
	if err != nil {
		t.Fatalf("parsing of the minimal example should have worked, got %v", err)
	}
}

var unknownKeyExample = minimalExample + `
some_other_key: with a value
`

func TestParseFailsOnUnknownKey(t *testing.T) {
	_, err := ParseBytes([]byte(unknownKeyExample))
	if err == nil {
		t.Fatalf("parsing of the example with unknown key should have failed")
	}
	if !strings.Contains(err.Error(), "some_other_key") {
		t.Errorf("error message should have named the invalid key")
	}
}

// smallExample defines a small scenario including instances of most
// configuration options.
var smallExample = `
name: Small Test
nodes:
  - name: A
    instances: 10
    features:
      - validator
      - archive
    start: 5
    end: 7.5

sources:
  - application: lottery
    instances: 10
    start: 7
    end: 10
    rate:
      constant: 8

  - application: my_coin
    rate:
      slope: 10

  - application: game
    rate:
      wave:
        min: 10
        max: 20
        period: 120
`

func TestParseSmallExampleWorks(t *testing.T) {
	_, err := ParseBytes([]byte(smallExample))
	if err != nil {
		t.Fatalf("parsing of input failed: %v", err)
	}
}
