// chris 082915 Keys for unique command line-invocation identification.

package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"crypto/md5"
)

const maxkeylen = 128

var nonwordre = regexp.MustCompile(`[^\w]+`)
var undscre = regexp.MustCompile(`_{2,}`)

func subnonwords(x string) string {
	x = nonwordre.ReplaceAllLiteralString(x, "_")
	x = undscre.ReplaceAllLiteralString(x, "_")
	return x
}

func MakeKey(command string, args []string) string {
	key := command
	a := subnonwords(strings.Join(args, "_"))
	if len(a) > 0 && a != "_" {
		key += "_" + a
	}
	key = subnonwords(key)
	if len(key) > maxkeylen {
		hash := fmt.Sprintf("%x", md5.Sum([]byte(key)))
		key = fmt.Sprintf("%s%s", key[:maxkeylen - len(hash)], hash)
	}
	return key
}
