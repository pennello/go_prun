// chris 082915 Keys for unique command line-invocation identification.

package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"crypto/md5"
)

// MaxKeyLength is the maximum key length used by MakeKey, over which it
// will use a hash at its end to stay within the length limit, but still
// uniquely and deterministically identify the given command and its
// arguments.
const MaxKeyLength = 128

var nonwordre = regexp.MustCompile(`[^\w]+`)
var undscre = regexp.MustCompile(`_{2,}`)

func subnonwords(x string) string {
	x = nonwordre.ReplaceAllLiteralString(x, "_")
	x = undscre.ReplaceAllLiteralString(x, "_")
	return x
}

// MakeKey produces a "key" for a given command and its arguments.
//
// The intent of the key is to produce a deterministic and reasonably
// human-readable string that identifies the command being run.  MakeKey
// consolidates all non-word characters between the command and its
// arguments and replaces them with underscores.
//
// If the length of the key exceeds MaxKeyLength, then it will be
// truncated and the suffix of the key will be a hash of the full key so
// as to stay within the length limit, but still uniquely and
// deterministically identify the given command and its arguments.
func MakeKey(command string, args []string) string {
	key := command
	a := subnonwords(strings.Join(args, "_"))
	if len(a) > 0 && a != "_" {
		key += "_" + a
	}
	key = subnonwords(key)
	key = strings.Trim(key, "_")
	if len(key) > MaxKeyLength {
		hash := fmt.Sprintf("%x", md5.Sum([]byte(key)))
		key = fmt.Sprintf("%s%s", key[:MaxKeyLength-len(hash)], hash)
	}
	return key
}
