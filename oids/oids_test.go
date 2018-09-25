package oids_test

import (
	"github.com/syncsynchalt/der2text/oids"
	"github.com/syncsynchalt/der2text/test"
	"testing"
)

func TestOids(t *testing.T) {
	test.Equals(t, "", oids.Name("1.2.3.4"))
	test.Equals(t, "Initials", oids.Name("2.5.4.43"))
}
