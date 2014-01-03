package pt

import (
	"bufio"
	"fmt"
	"github.com/monochromegane/the_platinum_searcher/util"
	"os"
	"path/filepath"
	"strings"
)

type Searcher struct {
	Root, Pattern string
}

type GrepArgument struct {
	Path, Pattern string
}

type PrintArgument struct {
        Match string
}

func (self *Searcher) Search() {
	grep := make(chan *GrepArgument, 2)
	match := make(chan *PrintArgument, 2)
	done := make(chan bool)
	go self.find(grep)
	go self.grep(grep, match)
	go self.print(match, done)
	<-done
}

func (self *Searcher) find(grep chan *GrepArgument) {
	filepath.Walk(self.Root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		fileType := pt.IdentifyFileType(path)
		if fileType == pt.BINARY {
			return nil
		}
		grep <- &GrepArgument{path, self.Pattern}
		return nil
	})
	grep <- nil
}

func (self *Searcher) grep(grep chan *GrepArgument, match chan *PrintArgument) {
	for {
		arg := <-grep
		if arg == nil {
			break
		}

		fh, err := os.Open(arg.Path)
		f := bufio.NewReader(fh)
		if err != nil {
			panic(err)
		}
		buf := make([]byte, 1024)

		for {
			buf, _, err = f.ReadLine()
			if err != nil {
				break
			}

			s := string(buf)
			if strings.Contains(s, arg.Pattern) {
				match <- &PrintArgument{s}
			}
		}
		fh.Close()

	}
	match <- nil
}

func (self *Searcher) print(match chan *PrintArgument, done chan bool) {
	for {
		matched := <-match
		if matched == nil {
			break
		}
		fmt.Printf("matched: %s\n", matched.Match)
	}
	done <- true
}
