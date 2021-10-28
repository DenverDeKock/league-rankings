package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	var (
		err error
	)

	teamPoints := make(map[string]int)

	err = scanTxtAndPopulateTeamPoints(os.Args[1], teamPoints)
	if err != nil {
		log.Fatal(err)
	}

	orderedPoints := orderTeamPoints(teamPoints)
	generateOutput(orderedPoints)
}

func scanTxtAndPopulateTeamPoints(fileName string, teamPoints map[string]int) error {
	var (
		scanner *bufio.Scanner
		file    *os.File
		err     error
	)

	file, err = os.Open(fileName)
	if err != nil {
		return err
	}

	defer file.Close()

	scanner = bufio.NewScanner(file)

	for scanner.Scan() {
		team1Name, team1Points, team2Name, team2Points, err := deriveTeamsFromText(strings.TrimSpace(scanner.Text()))

		if err != nil {
			fmt.Printf("%+v", err)
		}

		switch {
		case team1Points > team2Points:
			teamPoints[team1Name] += 3
			teamPoints[team2Name] += 0
		case team2Points > team1Points:
			teamPoints[team2Name] += 3
			teamPoints[team2Name] += 0
		case team1Points == team2Points:
			teamPoints[team1Name]++
			teamPoints[team2Name]++
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func deriveTeamsFromText(entry string) (string, int, string, int, error) {
	var (
		comma, whiteSpace1, whiteSpace2, team1Points, team2Points int
		t1str, t2str                                              string
		err                                                       error
	)

	comma = strings.Index(entry, ",")
	t1str = strings.TrimSpace(entry[:comma])
	t2str = strings.TrimSpace(entry[comma+1:])
	whiteSpace1 = strings.LastIndex(t1str, " ")
	whiteSpace2 = strings.LastIndex(t2str, " ")

	team1Points, err = strconv.Atoi(t1str[whiteSpace1+1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	team2Points, err = strconv.Atoi(t2str[whiteSpace2+1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	return t1str[:whiteSpace1], team1Points, t2str[:whiteSpace2], team2Points, err
}

func orderTeamPoints(teamPoints map[string]int) TeamPointsPairList {
	var (
		i int
	)

	keys := make([]string, 0, len(teamPoints))

	for k := range teamPoints {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	points := make(TeamPointsPairList, len(keys))

	for _, v := range keys {
		points[i] = TeamPointsPair{v, teamPoints[v]}
		i++
	}

	sort.Sort(points)

	return points
}

type TeamPointsPair struct {
	Key   string
	Value int
}

type TeamPointsPairList []TeamPointsPair

func (t TeamPointsPairList) Len() int {
	return len(t)
}

func (t TeamPointsPairList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t TeamPointsPairList) Less(i, j int) bool {
	return t[i].Value > t[j].Value
}

func generateOutput(orderedPoints TeamPointsPairList) {
	i := 0
	k := 1

	for j, v := range orderedPoints {
		i++
		if j > 0 && (v.Value != orderedPoints[j-1].Value) {
			k = i
		}

		if v.Value == 1 {
			fmt.Printf("%d. %s, %d pt\n", k, v.Key, v.Value)

			continue
		}

		fmt.Printf("%d. %s, %d pts\n", k, v.Key, v.Value)
	}
}
