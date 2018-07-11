package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type LogLevel int

const (
	DEV LogLevel = 0
	DBG LogLevel = 1
	INF LogLevel = 2
	ERR LogLevel = 3
)

type Entry struct {
	Time  time.Time
	Level LogLevel
	File  string
	Line  int
	Text  string
	TID   int
	OK    bool
}

func NewEntry() *Entry {
	return &Entry{}
}

func (v *Entry) Parse(line string) error {
	// 2018-07-02 11:16:40,944 DEBUG: xmlrpc.c:329 - Call to print-provider.registerPrinter3 returned: SUCCESS [2876]
	s := `(?P<Time>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2},\d{3})\s+`
	s += `(?P<Level>\w+):\s+`
	s += `(?P<File>\S+)\s*:\s*`
	s += `(?P<Line>\d+)\s+-\s+`
	s += `(?P<Text>.+)\s+\[`
	s += `(?P<TID>.+)\]`
	re := regexp.MustCompile(s)
	names := re.SubexpNames()
	matches := re.FindStringSubmatch(line)
	if matches == nil || len(matches) != 7 {
		v.OK = false
		return errors.New("parse error")
	}
	md := map[string]string{}
	for i, n := range matches {
		fmt.Printf("%d. match='%s'\tname='%s'\n", i, n, names[i])
		md[names[i]] = n
	}
	v.Time = parseDateTime(md["Time"])
	v.Level = parseLevel(md["Level"])
	v.OK = true
	return nil
}

func (p Entry) String() string {
	return fmt.Sprintf("%v - %v", p.Text, p.OK)
}

func parseDateTime(s string) (t time.Time) {
	t , err := time.Parse("2006-01-02 15:04:05.000", strings.Replace(s, ",", ".", 1))
	if err != nil {
		t = time.Time{}
	}
	return
}

func parseLevel(s string) (l LogLevel) {
	switch s {
	case "DEV":
		l = DEV
	default:
		l = DBG
	}
	return
}