package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanTxtAndPopulateTeamPoints(t *testing.T) {
	inputMap := make(map[string]int)

	var tests = []struct {
		input       string
		description string
		expected    map[string]int
	}{
		{
			"inputs/input1.txt",
			"test 1 - basic input",
			map[string]int{"Tarantulas": 6, "Grouches": 0, "FC Awesome": 1, "Lions": 5, "Snakes": 1},
		},
		{
			"inputs/input2.txt",
			"test 2 - matching scores",
			map[string]int{"aa team": 6, "ab team": 6, "ba team": 3, "bb team": 3, "bc team": 3, "ca team": 1, "cb team": 1},
		},
		{
			"inputs/input3.txt",
			"test 3 - large input",
			map[string]int{"FC Awesome": 14336, "Grouches": 0, "Lions": 71680, "Snakes": 14336, "Tarantulas": 86016},
		},
		{
			"inputs/input4.txt",
			"test 4 - more matching scores",
			map[string]int{"aa team": 3, "ab team": 1, "ba team": 6, "cb team": 3, "cc team": 1, "da team": 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			_ = scanTxtAndPopulateTeamPoints(tt.input, inputMap)
			assert.Equal(t, tt.expected, inputMap, "they should be equal")

			inputMap = make(map[string]int)
		})
	}
}

func TestDeriveTeamsFromText(t *testing.T) {
	var tests = []struct {
		input       string
		description string
		expected    []interface{}
	}{
		{
			"Tarantulas 1, FC Awesome 0",
			"test 1 - basic names",
			[]interface{}{
				"Tarantulas",
				1,
				"FC Awesome",
				0,
			},
		},
		{
			"Best in the Game 5, Worst in the game 0",
			"test 2 - more white space",
			[]interface{}{
				"Best in the Game",
				5,
				"Worst in the game",
				0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			g1, g2, g3, g4, _ := deriveTeamsFromText(tt.input)
			assert.Equal(t, tt.expected[0], g1, "they should be equal")
			assert.Equal(t, tt.expected[1], g2, "they should be equal")
			assert.Equal(t, tt.expected[2], g3, "they should be equal")
			assert.Equal(t, tt.expected[3], g4, "they should be equal")
		})
	}
}

func TestOrderTeamPoints(t *testing.T) {
	var tests = []struct {
		input       map[string]int
		description string
		expected    TeamPointsPairList
	}{
		{
			map[string]int{"Tarantulas": 6, "Grouches": 0, "FC Awesome": 1, "Lions": 5, "Snakes": 1},
			"test 1 - basic input",
			[]TeamPointsPair{
				{Key: "Tarantulas", Value: 6},
				{Key: "Lions", Value: 5},
				{Key: "FC Awesome", Value: 1},
				{Key: "Snakes", Value: 1},
				{Key: "Grouches", Value: 0},
			},
		},
		{
			map[string]int{"bb team": 3, "ab team": 6, "ca team": 1, "aa team": 6, "bc team": 3, "cb team": 1, "ba team": 3},
			"test 2 - matching scores and alphabetic ordering",
			[]TeamPointsPair{
				{Key: "aa team", Value: 6},
				{Key: "ab team", Value: 6},
				{Key: "ba team", Value: 3},
				{Key: "bb team", Value: 3},
				{Key: "bc team", Value: 3},
				{Key: "ca team", Value: 1},
				{Key: "cb team", Value: 1},
			},
		},
		{
			map[string]int{"aa team": 3, "ab team": 1, "ba team": 6, "cb team": 3, "cc team": 1, "da team": 6},
			"test 3 - more matching scores and alphabetic ordering",
			[]TeamPointsPair{
				{Key: "ba team", Value: 6},
				{Key: "da team", Value: 6},
				{Key: "aa team", Value: 3},
				{Key: "cb team", Value: 3},
				{Key: "ab team", Value: 1},
				{Key: "cc team", Value: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got := orderTeamPoints(tt.input)
			assert.Equal(t, tt.expected, got, "they should be equal")
		})
	}
}

func TestGenerateOutput(t *testing.T) {

	var tests = []struct {
		input       TeamPointsPairList
		description string
		expected    string
	}{
		{
			[]TeamPointsPair{
				{Key: "Tarantulas", Value: 6},
				{Key: "Lions", Value: 5},
				{Key: "FC Awesome", Value: 1},
				{Key: "Snakes", Value: 1},
				{Key: "Grouches", Value: 0},
			},
			"test 1",
			"1. Tarantulas, 6 pts\n2. Lions, 5 pts\n3. FC Awesome, 1 pt\n3. Snakes, 1 pt\n5. Grouches, 0 pts\n",
		},
		{
			[]TeamPointsPair{
				{Key: "aa team", Value: 6},
				{Key: "ab team", Value: 6},
				{Key: "ba team", Value: 3},
				{Key: "bb team", Value: 3},
				{Key: "bc team", Value: 3},
				{Key: "ca team", Value: 1},
				{Key: "cb team", Value: 1},
			},
			"test 2",
			"1. aa team, 6 pts\n1. ab team, 6 pts\n3. ba team, 3 pts\n3. bb team, 3 pts\n3. bc team, 3 pts\n6. ca team, 1 pt\n6. cb team, 1 pt\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			generateOutput(tt.input)
			w.Close()

			out, _ := ioutil.ReadAll(r)
			os.Stdout = rescueStdout

			assert.Equal(t, tt.expected, string(out), "they should be equal")
		})
	}
}
