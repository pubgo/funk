package env

import (
	"os"
	"strings"

	"github.com/pubgo/funk/v2/result"
)

// Expand substitutes environment variable references in value.
//
// Supported syntax:
//   - ${NAME}            value of NAME, empty when unset
//   - ${NAME:default}    default when NAME is unset or empty
//   - ${NAME:-default}   same as ${NAME:default}
//   - ${NAME:"default"}  default with surrounding quotes stripped
//   - $NAME              bare form, name starts with a letter or '_'
//
// Names are resolved against the raw environment first, then against the
// normalized key (see keyHandler), so variables renamed by Set/Reload
// still resolve: ${T2HOST} finds T2_HOST stored via env.Set("T2HOST", ...).
func Expand(value string) result.Result[string] {
	return result.OK(expandString(value))
}

func expandString(s string) string {
	if !strings.Contains(s, "$") {
		return s
	}

	var b strings.Builder
	b.Grow(len(s))

	for i := 0; i < len(s); {
		c := s[i]
		if c != '$' {
			b.WriteByte(c)
			i++
			continue
		}

		// "$$" escapes a single dollar.
		if i+1 < len(s) && s[i+1] == '$' {
			b.WriteByte('$')
			i += 2
			continue
		}

		if i+1 < len(s) && s[i+1] == '{' {
			end := strings.IndexByte(s[i+2:], '}')
			if end < 0 {
				b.WriteString(s[i:])
				break
			}

			body := s[i+2 : i+2+end]
			name, def := expandSplitDefault(body)
			b.WriteString(expandLookupOrDefault(name, def))
			i += len(body) + 3
			continue
		}

		j := i + 1
		for j < len(s) && isExpandNameChar(s[j]) {
			j++
		}
		// Bare form requires a name starting with a letter or '_'.
		if j == i+1 || isDigit(s[i+1]) {
			b.WriteByte('$')
			i++
			continue
		}

		b.WriteString(expandLookupOrDefault(s[i+1:j], ""))
		i = j
	}

	return b.String()
}

// lookupEnv resolves name against the process environment, trying the raw
// name first, then the keyHandler-normalized name, then the normalized
// uppercase name. Set/Reload store keys through keyHandler, which renames at
// digit->uppercase boundaries in a case-sensitive way (T2HOST -> T2_HOST,
// while sha256sum normalizes to SHA256SUM but SHA256SUM to SHA256_SUM), so
// the extra candidates make raw, mixed-case and lowercase spellings resolve.
func lookupEnv(name string) (string, bool) {
	for _, key := range envKeyCandidates(name) {
		if v, ok := os.LookupEnv(key); ok {
			return v, true
		}
	}
	return "", false
}

// envKeyCandidates returns the lookup order for a user-supplied key.
func envKeyCandidates(name string) []string {
	return []string{name, keyHandler(name), keyHandler(strings.ToUpper(name))}
}

// expandLookupOrDefault resolves name from the environment and falls back to
// def when no candidate is set to a non-empty value.
func expandLookupOrDefault(name, def string) string {
	for _, key := range envKeyCandidates(name) {
		if v, ok := os.LookupEnv(key); ok && v != "" {
			return v
		}
	}
	return expandStripQuotes(def)
}

func expandSplitDefault(body string) (name, def string) {
	before, after, ok := strings.Cut(body, ":")
	if !ok {
		return body, ""
	}
	// "${NAME:-default}" is the same operator as "${NAME:default}".
	def = strings.TrimPrefix(after, "-")
	return before, def
}

// expandStripQuotes removes one layer of matching surrounding quotes.
func expandStripQuotes(s string) string {
	if len(s) >= 2 && (s[0] == s[len(s)-1]) && (s[0] == '"' || s[0] == '\'') {
		return s[1 : len(s)-1]
	}
	return s
}

func isExpandNameChar(c byte) bool {
	return c == '_' || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9')
}

func isDigit(c byte) bool { return '0' <= c && c <= '9' }
