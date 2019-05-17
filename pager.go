package pager

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	tm "github.com/buger/goterm"
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
		// Setup less as a pager
		pager := exec.Command("less")
		// Get a pipe that will be connected to stdin on start-up
		w, _ := pager.StdinPipe()
		pager.Stdout = p.Output
		pager.Stderr = os.Stderr
		pager.Start()
		// Write the string representation of Buffer to the stdin pipe for 'less'
		fmt.Fprintf(w, bufString)
		// Close the stdin pipe (causing 'less' to exit cleanly)
		w.Close()
		// Wait for the process to complete
		pager.Wait()
	}
}

// New returns a new Pager object that will write to os.Stdout
func New() Pager {
	return Pager{Output: os.Stdout, Buffer: &bytes.Buffer{}}
}
