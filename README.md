# Description

This library should be used if you are printing to stdout and want to automatically display the output via `less`.

# Basic Usage

Pass a Pager object to a function that accepts the io.Writer interface. All output will be written to a buffer. When you wish to display the output, call the Page method. If the number of output lines exceeds the height of the terminal, it will be piped to `less`:

```go
import (
    "fmt"
    "github.com/walkert/pager"
)

func ManyLines(output io.Writer) {
    for i := 0; i < 100; i++ {
        fmt.Fprintf(output, "Hello, world!\n")
    }
}

func main() {
    p := pager.New()
    ManyLines(p)
    p.Page()
}
```
