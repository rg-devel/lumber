package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

type Parameters struct {
	InFile *string
	Action *string
}

func main() {

	parameters := parseFlags()
	fmt.Println(*parameters.InFile)
	fmt.Println(*parameters.Action)

	// Read the input log file and parse the contents
	contents, err := os.Open(*parameters.InFile)
	if err != nil {
		panic(err)
	}
	entries := ReadFile(contents)
	for _, v := range entries {
		fmt.Println(v)
	}

}

func parseFlags() Parameters {
	parameters := Parameters{
		InFile: flag.String("i", "print-provider.log", "input log file"),
		Action: flag.String("a", "list-job", "list all print jobs"),
	}
	flag.Parse()

	return parameters
}

func ReadFile(r io.Reader) (LogEntries []*Entry) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		t := scanner.Text()
		entry := NewEntry()
		entry.Parse(t)
		LogEntries = append(LogEntries, entry)
	}
	return
}
