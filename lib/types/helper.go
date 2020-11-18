package types

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strings"

	cp "github.com/otiai10/copy"
)

// CopyDir moves over the dir and skips git dirs
func CopyDir(src, dst string) (err error) {
	err = cp.Copy(src, dst,
		cp.Options{
			Skip: func(src string) (bool, error) {
				isGit := strings.Contains(src, "/.git")
				if isGit {
					log.Printf("Skip: %s", src)
				}
				return isGit, nil
			},
		},
	)
	return
}

func compVal(i1, i2 interface{}, fails []string) (res []string) {
	res = fails
	t1 := reflect.TypeOf(i1)
	t2 := reflect.TypeOf(i2)
	if t1 != t2 {
		fmt.Printf("%v :i1 // i2: %v", i1, i2)
		fmt.Printf(" ==> Wrong TYPE\n")
		res = append(fails, fmt.Sprintf("Type: i1:%s != i2:%s", t1, t2))
		return
	}
	if i1 != i2 {
		fmt.Printf("%v :i1 // i2: %v", i1, i2)
		fmt.Printf(" ==> Do not match\n")
		res = append(fails, fmt.Sprintf("Value: i1:%v != i2:%v", i1, i2))
	}
	return
}

func readFile(path string) ([]byte, error) {
	yamlFile, err := ioutil.ReadFile(path)
	return yamlFile, err
}

func updateMetaLines(inputLines, metaLines []string) (err error, outputLines []string) {
	sawDashes := 0
	for _, line := range inputLines {
		if strings.Contains(line, "---") {
			if sawDashes == 0 {
				outputLines = append(outputLines, line)
				outputLines = append(outputLines, metaLines...)
			}
			sawDashes++
		}
		switch {
		case sawDashes == 1:
			continue
		case sawDashes > 1:
			outputLines = append(outputLines, line)
			continue
		}
	}
	return
}
