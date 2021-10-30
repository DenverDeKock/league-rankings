package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanTxtAndPopulateTeamPoints(t *testing.T) {
	inputMap := make(map[string]int, 0)

	var tests = []struct {
		input       string
		description string
		expected    map[string]int
	}{
		{
			"input.txt",
			"test 1",
			map[string]int{"Tarantulas": 6, "Grouches": 0, "FC Awesome": 1, "Lions": 5, "Snakes": 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			_ = scanTxtAndPopulateTeamPoints(tt.input, inputMap)
			assert.Equal(t, tt.expected, inputMap, "they should be equal")
			inputMap = make(map[string]int, 0)
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
			"test 1",
			[]interface{}{
				"Tarantulas",
				1,
				"FC Awesome",
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
			"test 1",
			[]TeamPointsPair{
				{Key: "Tarantulas", Value: 6},
				{Key: "Lions", Value: 5},
				{Key: "FC Awesome", Value: 1},
				{Key: "Snakes", Value: 1},
				{Key: "Grouches", Value: 0},
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
