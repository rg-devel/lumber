package main

import (
	"time"
)

type Entry struct {
	Time       time.Time
	Level      string
	SourceFile string
	LineNumber int
	Message    string
	TID        string
	OK         bool
}

//func (p Entry) String() string {
//	return fmt.Sprintf("%v - %v", p.Message, p.OK)
//}
