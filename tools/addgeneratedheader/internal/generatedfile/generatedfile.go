package generatedfile

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
)

var genPrefix = []byte(`// Code generated `)
var genSuffix = []byte(`DO NOT EDIT.`)

//HasGeneratedComment returns true if there is a code generated comment in the file
func HasGeneratedComment(file string) (bool, error) {
	fl, err := os.Open(file) //nolint:gosec
	if err != nil {
		return false, err
	}
	scanner := bufio.NewScanner(fl)
	for scanner.Scan() {
		line := scanner.Bytes()
		if bytes.HasPrefix(line, genPrefix) && bytes.HasSuffix(line, genSuffix) {
			return true, scanner.Err()
		}
	}

	return false, scanner.Err()
}

//AddGeneratedComment adds a code generated comment if one doesn't already exist
func AddGeneratedComment(file string, byStatement string) error {
	alreadyGenerated, err := HasGeneratedComment(file)
	if err != nil {
		return err
	}
	if alreadyGenerated {
		return nil
	}
	content, err := ioutil.ReadFile(file) //nolint:gosec
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	mustWrite(&buf, genPrefix)
	mustWrite(&buf, []byte(byStatement))
	mustWrite(&buf, []byte(" "))
	mustWrite(&buf, genSuffix)
	mustWrite(&buf, []byte("\n"))
	mustWrite(&buf, content)
	return ioutil.WriteFile(file, buf.Bytes(), 0600)
}

func mustWrite(w io.Writer, p []byte) {
	_, err := w.Write(p)
	if err != nil {
		panic(err)
	}
}
