package io

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	learn "github.com/RenatoGeh/gospn/src/learn"
)

func GetPath(relpath string) string {
	rp, err := filepath.Abs(relpath)

	if err != nil {
		fmt.Printf("Error retrieving path \"%s\".\n", relpath)
		panic(err)
	}

	return rp
}

// Reads from a file named filename and returns a matrix of
func ParseData(filename string) (map[int]learn.Variable, []map[int]int) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error. Could not open file [%s].\n", filename)
		panic(err)
	}
	defer file.Close()

	sc := make(map[int]learn.Variable)

	scanner := bufio.NewScanner(file)

	var line string = ""

	// Get variable definitions.
	for {
		if !scanner.Scan() {
			break
		}
		line = scanner.Text()
		if line[0] != 'v' {
			break
		}
		var varid, cats int
		fmt.Sscanf(line, "var %d %d", &varid, &cats)
		sc[varid] = learn.Variable{varid, cats}
	}

	n := len(sc)
	var data [][]int

	regex := regexp.MustCompile("[\\,\\s]+")
	// We assume complete data.
	for i := 0; scanner.Scan(); i++ {
		data = append(data, make([]int, n))
		s := regex.Split(line, -1)
		for j := 0; j < n; j++ {
			data[i][j], err = strconv.Atoi(s[j])
			if err != nil {
				fmt.Printf("Invalid string \"%s\" found in data file [%s].\n", s[j], filename)
				panic(err)
			}
		}
		line = scanner.Text()
	}

	m, n := len(data), len(data[0])
	cvntmap := make([]map[int]int, m)
	for i := 0; i < m; i++ {
		cvntmap[i] = make(map[int]int)
		for j := 0; j < n; j++ {
			cvntmap[i][j] = data[i][j]
		}
	}

	return sc, cvntmap
}
