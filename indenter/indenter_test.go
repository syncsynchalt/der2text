package indenter_test

import (
	. "github.com/syncsynchalt/der2text/indenter"
	"github.com/syncsynchalt/der2text/test"
	"testing"
)

// an io.Writer that builds a string
type stringWriter struct {
    str string
}
func (s *stringWriter) Write(p []byte) (n int, err error) {
    s.str += string(p)
    return len(p), nil
}

func TestIndenterLevelZero(t *testing.T) {
	o := stringWriter{}
	i := New(&o)
	test.Equals(t, "", o.str)

	i.Write([]byte("foo"))
	test.Equals(t, "foo", o.str)
	i.Write([]byte("bar"))
	test.Equals(t, "foobar", o.str)
	i.Write([]byte("\nbaz"))
	test.Equals(t, "foobar\nbaz", o.str)
	i.Write([]byte("\n"))
	test.Equals(t, "foobar\nbaz\n", o.str)
}

func TestIndenterLevelOne(t *testing.T) {
	o := stringWriter{}
	i := New(&o)
	test.Equals(t, "", o.str)

	i.Write([]byte("foo\n"))
	test.Equals(t, "foo\n", o.str)

	i = i.NextLevel()
	test.Equals(t, "foo\n", o.str)

	i.Write([]byte("bar"))
	test.Equals(t, "foo\n  bar", o.str)
	i.Write([]byte("\nbaz"))
	test.Equals(t, "foo\n  bar\nbaz", o.str)
	i.Write([]byte("\n"))
	test.Equals(t, "foo\n  bar\nbaz\n", o.str)
	i.Write([]byte("a"))
	test.Equals(t, "foo\n  bar\nbaz\n  a", o.str)
	i.Write([]byte("b"))
	test.Equals(t, "foo\n  bar\nbaz\n  ab", o.str)
}

func TestIndenterLevelTwo(t *testing.T) {
	o := stringWriter{}
	i := New(&o)

	i.Write([]byte("foo\n"))
	test.Equals(t, "foo\n", o.str)

	i = i.NextLevel()
	i.Write([]byte("bar\n"))
	test.Equals(t, "foo\n  bar\n", o.str)

	i = i.NextLevel()
	i.Write([]byte("baz\n"))
	test.Equals(t, "foo\n  bar\n    baz\n", o.str)
}

func TestIndenterLevelBackout(t *testing.T) {
	o := stringWriter{}
	i1 := New(&o)
	test.Equals(t, "", o.str)

	i1.Write([]byte("foo\n"))
	test.Equals(t, "foo\n", o.str)

	i2 := i1.NextLevel()
	i2.Write([]byte("bar\n"))
	test.Equals(t, "foo\n  bar\n", o.str)

	i1.Write([]byte("baz\n"))
	test.Equals(t, "foo\n  bar\nbaz\n", o.str)
}
