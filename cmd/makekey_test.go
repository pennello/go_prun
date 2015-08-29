// chris 082915

package cmd

import (
	"strings"
	"testing"
)

func testMakeKeyExpect(t *testing.T, command string, args []string, expect string) {
	key := MakeKey(command, args)
	t.Logf("key %q\n", key)
	if key != expect {
		t.Errorf("MakeKey(%q, %q) != %q (got %q)\n", command, args, expect, key)
	}
}

func testMakeKeyLong(t *testing.T, command string, args []string) {
	key := MakeKey(command, args)
	t.Logf("long key %q\n", key)
	if len(key) > maxkeylen {
		t.Errorf("%q longer than %d\n", key, maxkeylen)
	}
}

func TestMakeKey(t *testing.T) {
	testMakeKeyExpect(t, "ls", []string{"-l"}, "ls_l")
	testMakeKeyExpect(t, "ls", []string{""}, "ls")
	testMakeKeyExpect(t, "ls", []string{"", ""}, "ls")
	testMakeKeyExpect(t, "", []string{}, "")
	testMakeKeyExpect(t, "", []string{}, "")
	testMakeKeyExpect(t, "", []string{"x"}, "_x")
	testMakeKeyExpect(t, "grep", []string{"-Rw", "blah", "."}, "grep_Rw_blah_")

	var longcommand string

	longcommand = strings.Repeat("a", maxkeylen - 2)
	testMakeKeyExpect(t, longcommand, []string{}, longcommand)
	testMakeKeyExpect(t, longcommand, []string{"x"}, longcommand + "_x")
	longcommand = strings.Repeat("b", maxkeylen - 1)
	testMakeKeyExpect(t, longcommand, []string{}, longcommand)
	longcommand = strings.Repeat("c", maxkeylen)
	testMakeKeyExpect(t, longcommand, []string{}, longcommand)

	testMakeKeyLong(t, strings.Repeat("d", maxkeylen + 1), []string{})
	testMakeKeyLong(t, strings.Repeat("e", 2 * maxkeylen), []string{})
}
