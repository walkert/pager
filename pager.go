package pager

import (
	"bytes"
	"fmt"
	tm "github.com/buger/goterm"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Pager implements the io.Writer interface and holds an output
// destination and buffer object
type Pager struct {
	Output io.Writer
	Buffer *bytes.Buffer
}

// Write implements the Write method of the Buffer
func (p Pager) Write(b []byte) (n int, err error) {
	return p.Buffer.Write(b)
}

// Page writes the string representation of Buffer directly
// to Output or via 'less' if the line count is > terminal height
func (p Pager) Page() {
	bufString := p.Buffer.String()
	lines := strings.Split(bufString, "\n")
	if len(lines) <= tm.Height() {
		fmt.Fprintf(p.Output, bufString)
	} else {
		// Open a new pipe
		pipeRead, pipeWrite := io.Pipe()
		pager := exec.Command("less")
		// Attach the read-side to the stdin of 'less'
		pager.Stdin = pipeRead
		pager.Stdout = p.Output
		pager.Stderr = os.Stderr
		pager.Start()
		// Write the string representation of Buffer to write-side
		// of the pipe
		fmt.Fprintf(pipeWrite, bufString)
		// Close the read-side (causing 'less' to exit cleanly)
		pipeRead.Close()
		// Wait for the process to complete
		pager.Wait()
	}
}

// New returns a new Pager object that will write to os.Stdout
func New() Pager {
	return Pager{Output: os.Stdout, Buffer: &bytes.Buffer{}}
}
