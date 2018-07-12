package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {

	parameters := parseFlags()
	fmt.Println(*parameters.InFile)
	fmt.Println(*parameters.Action)

	// Read the input log file and parse the contents
	contents, err := os.Open(*parameters.InFile)
	if err != nil {
		panic(err)
	}
	logEntries := ParseFile(contents)
	for _, entry := range logEntries {
		//if !entry.OK {
		fmt.Printf("%+v\n", entry.Message)
		//}
	}

}

func ParseFile(r io.Reader) (entries []Entry) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		s := scanner.Text()
		line := Line(s)
		var entry Entry
		entry.Level, entry.OK = line.Level()
		entry.Time, entry.OK = line.Time()

		if entry.Level == "DEV" || entry.Level == "DEBUG" {
			entry.SourceFile, entry.OK = line.SourceFile()
			entry.LineNumber, entry.OK = line.LineNumber()
			entry.Message, entry.OK = line.MessageWhenDebug()
			entry.TID, entry.OK = line.TID()
		} else {
			entry.Message, entry.OK = line.Message()
		}

		entries = append(entries, entry)
	}
	return
}
