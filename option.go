package main

import "flag"

type Parameters struct {
	InFile *string
	Action *string
}

func parseFlags() Parameters {
	parameters := Parameters{
		InFile: flag.String("i", "print-provider.log", "input log file"),
		Action: flag.String("a", "list-job", "list all print jobs"),
	}
	flag.Parse()

	return parameters
}
