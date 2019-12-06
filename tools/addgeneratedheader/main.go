package main

import (
	"fmt"
	"os"

	"github.com/killa-beez/gopkgs/tools/addgeneratedheader/internal/generatedfile"
)

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Printf(`usage: %s file [attribution statement]

prefixes file with a line like 
// Code generated <attribution statement> DO NOT EDIT.

Does nothing if file already has a line that starts with "// Code generated" and ends with "DO NOT EDIT."

`, os.Args[0])
		os.Exit(2)
	}

	file := os.Args[1]
	statement := ""
	if len(os.Args) == 3 {
		statement = os.Args[2]
	}
	err := generatedfile.AddGeneratedComment(file, statement)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
