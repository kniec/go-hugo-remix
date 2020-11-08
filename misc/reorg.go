package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"path/filepath"
	"regexp"
	"strings"
)

type Tupel struct {
	Path   *regexp.Regexp
	Weight string
	Prefix string
}

func (t *Tupel) Visit(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !!fi.IsDir() {
		return nil //
	}

	matched := t.Path.MatchString(path)

	if err != nil {
		panic(err)
		return err
	}

	if matched {
		fmt.Printf("%s / Going to apply weight:%s, prefix:%s\n", path, t.Weight, t.Prefix)
		read, err := ioutil.ReadFile(path)
		content := string(read)
		if err != nil {
			panic(err)
		}
		r := regexp.MustCompile("weight = .*")
		if r.MatchString(content) {
			log.Printf(`Found regex 'weight = .*' in %s`, path)
			content = r.ReplaceAllString(content, fmt.Sprintf("weight = %s", t.Weight))
		}
		content = r.ReplaceAllString(content, fmt.Sprintf("weight = %s", t.Weight))
		r = regexp.MustCompile("weight: .*")
		if r.MatchString(content) {
			log.Printf(`Found regex 'weight: .*' in %s`, path)
			content = r.ReplaceAllString(content, fmt.Sprintf("weight: %s", t.Weight))

		}
		r = regexp.MustCompile(`pre = ".*"`)
		if r.MatchString(content) {
			log.Printf(`Found regex 'pre = ".*"' in %s`, path)
			content = r.ReplaceAllString(content, fmt.Sprintf(`pre = %s`, t.Prefix))
		}

		r = regexp.MustCompile(`prefix: ".*"`)
		if r.MatchString(content) {
			log.Printf(`Found regex 'prefix: ".*"' in %s`, path)
			content = r.ReplaceAllString(content, fmt.Sprintf(`prefix: %s`, t.Prefix))
		}
		err = ioutil.WriteFile(path, []byte(content), 0)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

// END TUpel

func readLine(fpath string) (lines []string) {
	file, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			fmt.Printf("Skip line: %s\n", line)
			continue
		}
		if len(strings.Split(line, ";")) == 0 {
			fmt.Printf("Skip line: %s\n", line)
			continue
		}

		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return lines
}

func splitTupel(line string) (t Tupel, err error) {
	s := strings.Split(line, ";")
	t.Path = regexp.MustCompile(s[0])
	ws := strings.Split(s[1], ":")
	if len(ws) != 2 {
		panic("should split in two")
	}
	t.Weight = ws[1]
	ps := strings.Split(s[2], ":")
	if len(ps) != 2 {
		panic("should split in two")
	}
	t.Prefix = ps[1]
	return
}

func walkFiles(startPath string, t Tupel) {
	err := filepath.Walk(startPath, t.Visit)
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("### read: ", os.Args[1])
	lines := readLine(os.Args[1])
	for _, line := range lines {
		t, _ := splitTupel(line)
		walkFiles("./content", t)
	}
}
