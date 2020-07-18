package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
}

const (
	usage = `firstly enter you git repository, execute:
  $ git log --author='xxx' --since=2020-3-1 --until=2020-6-30 --stat > ~/git-20h1.commits
secondly, run this script with this command:
  $ go run git_contributions.go git-20h1.commits
`
	regInsertions = `(\d+) insertions\(\+\)`
	regDeletions  = `(\d+) deletions\(-\)`
)

var (
	reInsertions = regexp.MustCompile(regInsertions)
	reDeletions  = regexp.MustCompile(regDeletions)
	reCommit     = regexp.MustCompile(`^commit [a-f0-9]{40}`)
	reTitle      = regexp.MustCompile(`^    \S+`)
	reAuthor     = regexp.MustCompile(`^Author: `)
	reDate       = regexp.MustCompile(`^Date:`)
)

func main() {
	if len(os.Args) != 2 || len(os.Args) > 1 && os.Args[1] == "help" {
		fmt.Println(usage)
		os.Exit(0)
	}
	// open file
	fn := os.Args[1]
	in, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	// parse line by line
	rd := bufio.NewReader(in)
	var (
		insertions, deletions int
		changed               int
		commits               int
		largeCommits          int
		matches               []string
		commit, title         string
		author, date          string
	)
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}
		// commit and title
		if reCommit.MatchString(line) {
			commit = line
			commits++
		} else if reTitle.MatchString(line) {
			title = line
		} else if reAuthor.MatchString(line) {
			author = line
		} else if reDate.MatchString(line) {
			date = line
		}
		// insertions and deletions
		if matches = reInsertions.FindStringSubmatch(line); len(matches) > 1 {
			n, err := strconv.Atoi(matches[1])
			if err != nil {
				log.Fatalf("parse line '%s' error %s", line, err)
			}
			insertions += n
			changed += n
		}
		if matches = reDeletions.FindStringSubmatch(line); len(matches) > 1 {
			n, err := strconv.Atoi(matches[1])
			if err != nil {
				log.Fatalf("parse line '%s' error %s", line, err)
			}
			deletions += n
			changed += n
		}
		if changed >= 1000 {
			largeCommits++
			fmt.Println("large commit", largeCommits)
			fmt.Println(strings.TrimSpace(commit))
			fmt.Println(strings.TrimSpace(title))
			fmt.Println(strings.TrimSpace(author))
			fmt.Println(strings.TrimSpace(date))
			fmt.Println(strings.TrimSpace(line))
			fmt.Println()
		}
		changed = 0
	}
	fmt.Printf("commits: %d, total lines: %d, additions: %d, deletions: %d\n",
		commits, insertions+deletions, insertions, deletions)
}
