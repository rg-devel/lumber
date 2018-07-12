package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 2018-06-29 17:16:26,963 DEBUG: PCPrintService.cpp:1927 - No configuration changes required. [3932]
// TimeStr + TimeLevelDivStr + LevelStr + LevelSourceNameDivStr + SourceNameStr + SourceNameLineNumDivStr + LineNumStr + LineNumMsgDivStr + MsgStr + MsgTIDDivStr + TIDStr + EndStr

// 2016-11-04 17:13:47,341 INFO : Service installed successfully.
// TimeStr + TimeLevelDivStr + LevelStr + LevelMsgDivStr + MsgStr
const (
	StartStr                string = `^`
	TimeStr                 string = `\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2},\d{3}`
	LevelStr                string = `(?:DEV|DEBUG|INFO|ERROR)`
	SourceNameStr           string = `\S+`
	LineNumStr              string = `\d+`
	MsgStr                  string = `.+`
	TIDStr                  string = `\d+`
	TimeLevelDivStr         string = `\s+`
	LevelSourceNameDivStr   string = `\s*:\s+`
	SourceNameLineNumDivStr string = `:`
	LineNumMsgDivStr        string = `\s+-\s+`
	MsgTIDDivStr            string = `\s+\[`
	EndStr                  string = `\]$`
	LevelMsgDivStr          string = LevelSourceNameDivStr
)

type Line string

func (line Line) String() string {
	return fmt.Sprintf("%s", string(line))
}

func (line Line) Time() (time.Time, bool) {
	s := StartStr +
		`(` + TimeStr + `)` +
		TimeLevelDivStr
	token, _ := line.Token(s)
	t, err := time.Parse("2006-01-02 15:04:05.000", strings.Replace(token, ",", ".", 1))
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

func (line Line) Level() (level string, ok bool) {
	s := StartStr +
		TimeStr +
		TimeLevelDivStr +
		`(` + LevelStr + `)` +
		LevelSourceNameDivStr
	return line.Token(s)
}

func (line Line) SourceFile() (string, bool) {
	// 2018-06-29 17:16:26,963 DEBUG: PCPrintService.cpp:1927 - No configuration changes required. [3932]
	// 2016-11-04 17:13:47,341 INFO : Service installed successfully.
	s := StartStr +
		TimeStr +
		TimeLevelDivStr +
		LevelStr +
		LevelSourceNameDivStr +
		`(` + SourceNameStr + `)` +
		SourceNameLineNumDivStr
	return line.Token(s)
}

func (line Line) LineNumber() (int, bool) {
	s := StartStr +
		TimeStr +
		TimeLevelDivStr +
		LevelStr +
		LevelSourceNameDivStr +
		SourceNameStr +
		SourceNameLineNumDivStr +
		`(` + LineNumStr + `)` +
		LineNumMsgDivStr
	v, _ := line.Token(s)

	lineNumber, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return lineNumber, true
}

func (line Line) MessageWhenDebug() (string, bool) {

	s := StartStr +
		TimeStr +
		TimeLevelDivStr +
		LevelStr +
		LevelSourceNameDivStr +
		SourceNameStr +
		SourceNameLineNumDivStr +
		LineNumStr +
		LineNumMsgDivStr +
		`(` + MsgStr + `)` +
		MsgTIDDivStr
	return line.Token(s)
}

func (line Line) Message() (string, bool) {

	// 2016-11-04 17:13:47,341 INFO : Service installed successfully.
	s := StartStr +
		TimeStr +
		TimeLevelDivStr +
		LevelStr +
		LevelMsgDivStr +
		`(` + MsgStr + `)` +
		`$`
	return line.Token(s)
}

func (line Line) TID() (string, bool) {

	s := StartStr +
		TimeStr +
		TimeLevelDivStr +
		LevelStr +
		LevelSourceNameDivStr +
		SourceNameStr +
		SourceNameLineNumDivStr +
		LineNumStr +
		LineNumMsgDivStr +
		MsgStr +
		MsgTIDDivStr +
		`(` + TIDStr + `)` +
		EndStr
	return line.Token(s)
}

func (line Line) Token(pattern string) (string, bool) {
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(string(line))
	if len(matches) != 2 {
		return "", false
	}
	return matches[1], true
}
