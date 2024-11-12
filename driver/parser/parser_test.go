// Copyright 2024 Fantom Foundation
// This file is part of Norma System Testing Infrastructure for Sonic.
//
// Norma is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Norma is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Norma. If not, see <http://www.gnu.org/licenses/>.

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
num_validators: 5
nodes:
  - name: A
    instances: 10
    features:
      - validator
      - archive
    start: 5
    end: 7.5

applications:
  - name: lottery
    instances: 10
    start: 7
    end: 10
    rate:
      constant: 8

  - name: my_coin
    rate:
      slope:
        start: 5
        increment: 1

  - name: game
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

// withClientType defines an example with client specification
var withClientType = `
name: Small Test
num_validators: 5
nodes:
  - name: A
    instances: 10
    features:
      - validator
      - archive
    start: 5
    end: 7.5
    client:
      imagename: main
      type: validator

applications:
  - name: lottery
    instances: 10
    start: 7
    end: 10
    rate:
      constant: 8

  - name: my_coin
    rate:
      slope:
        start: 5
        increment: 1

  - name: game
    rate:
      wave:
        min: 10
        max: 20
        period: 120
`

func TestParseWithClientTypeWorks(t *testing.T) {
	_, err := ParseBytes([]byte(withClientType))
	if err != nil {
		t.Fatalf("parsing of input failed: %v", err)
	}
}

var withCheats = smallExample + `

cheats:
  - name: hello
    start: 8
`

func TestParseExampleWithCheats(t *testing.T) {
	_, err := ParseBytes([]byte(withCheats))
	if err != nil {
		t.Fatalf("parsing of input failed: %v", err)
	}
}
