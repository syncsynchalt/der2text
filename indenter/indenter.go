package indenter

import (
	"fmt"
	"strings"
)

type Indenter struct {
	atEos bool
	level int
}

func New() Indenter {
	return Indenter{true, 0}
}

func (i *Indenter) Println(a ...interface{}) (n int, err error) {
	if i.atEos {
		fmt.Print(strings.Repeat(" ", i.level))
	}
	i.atEos = true
	return fmt.Println(a...)
}

func (i *Indenter) Printf(format string, a ...interface{}) (n int, err error) {
	return i.Print(fmt.Sprintf(format, a...))
}

func (i *Indenter) Print(a ...interface{}) (n int, err error) {
	if i.atEos {
		fmt.Print(strings.Repeat(" ", i.level))
	}
	s := fmt.Sprint(a...)
	if len(s) > 0 {
		i.atEos = s[len(s)-1] == '\n'
	}
	return fmt.Print(s)
}

func (i *Indenter) NextLevel() Indenter {
	return Indenter{true, i.level + 2}
}
