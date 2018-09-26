package indenter_test

import (
	. "github.com/syncsynchalt/der2text/read/indenter"
	"github.com/syncsynchalt/der2text/test"
	"strings"
	"testing"
)

func TestIndenterLevelZero(t *testing.T) {
	o := strings.Builder{}
	i := New(&o)
	test.Equals(t, "", o.String())

	i.Write([]byte("foo"))
	test.Equals(t, "foo", o.String())
	i.Write([]byte("bar"))
	test.Equals(t, "foobar", o.String())
	i.Write([]byte("\nbaz"))
	test.Equals(t, "foobar\nbaz", o.String())
	i.Write([]byte("\n"))
	test.Equals(t, "foobar\nbaz\n", o.String())
}

func TestIndenterLevelOne(t *testing.T) {
	o := strings.Builder{}
	i := New(&o)
	test.Equals(t, "", o.String())

	i.Write([]byte("foo\n"))
	test.Equals(t, "foo\n", o.String())

	i = i.NextLevel()
	test.Equals(t, "foo\n", o.String())

	i.Write([]byte("bar"))
	test.Equals(t, "foo\n  bar", o.String())
	i.Write([]byte("\nbaz"))
	test.Equals(t, "foo\n  bar\nbaz", o.String())
	i.Write([]byte("\n"))
	test.Equals(t, "foo\n  bar\nbaz\n", o.String())
	i.Write([]byte("a"))
	test.Equals(t, "foo\n  bar\nbaz\n  a", o.String())
	i.Write([]byte("b"))
	test.Equals(t, "foo\n  bar\nbaz\n  ab", o.String())
}

func TestIndenterLevelTwo(t *testing.T) {
	o := strings.Builder{}
	i := New(&o)

	i.Write([]byte("foo\n"))
	test.Equals(t, "foo\n", o.String())

	i = i.NextLevel()
	i.Write([]byte("bar\n"))
	test.Equals(t, "foo\n  bar\n", o.String())

	i = i.NextLevel()
	i.Write([]byte("baz\n"))
	test.Equals(t, "foo\n  bar\n    baz\n", o.String())
}

func TestIndenterLevelBackout(t *testing.T) {
	o := strings.Builder{}
	i1 := New(&o)
	test.Equals(t, "", o.String())

	i1.Write([]byte("foo\n"))
	test.Equals(t, "foo\n", o.String())

	i2 := i1.NextLevel()
	i2.Write([]byte("bar\n"))
	test.Equals(t, "foo\n  bar\n", o.String())

	i1.Write([]byte("baz\n"))
	test.Equals(t, "foo\n  bar\nbaz\n", o.String())
}
