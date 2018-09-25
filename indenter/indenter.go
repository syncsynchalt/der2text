package indenter

import (
	"fmt"
	"io"
	"strings"
)

type Indenter struct {
	atEos  bool
	level  int
	writer io.Writer
}

func New(writer io.Writer) *Indenter {
	return &Indenter{true, 0, writer}
}

func (i *Indenter) Println(a ...interface{}) (n int, err error) {
	if i.atEos {
		fmt.Fprint(i.writer, strings.Repeat(" ", i.level))
	}
	i.atEos = true
	return fmt.Fprintln(i.writer, a...)
}

func (i *Indenter) Printf(format string, a ...interface{}) (n int, err error) {
	return i.Print(fmt.Sprintf(format, a...))
}

func (i *Indenter) Print(a ...interface{}) (n int, err error) {
	if i.atEos {
		fmt.Fprint(i.writer, strings.Repeat(" ", i.level))
	}
	s := fmt.Sprint(a...)
	if len(s) > 0 {
		i.atEos = s[len(s)-1] == '\n'
	}
	return fmt.Fprint(i.writer, s)
}

func (i *Indenter) Write(p []byte) (n int, err error) {
	if i.atEos {
		fmt.Fprint(i.writer, strings.Repeat(" ", i.level))
	}
	if len(p) > 0 {
		i.atEos = p[len(p)-1] == '\n'
	}
	return i.writer.Write(p)
}

func (i *Indenter) NextLevel() *Indenter {
	return &Indenter{true, i.level + 2, i.writer}
}
