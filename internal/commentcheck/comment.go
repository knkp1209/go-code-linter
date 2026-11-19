package commentcheck

import (
	"go/ast"
	"strings"
	"unicode"
	"unicode/utf8"
)

func HasNamePrefixedComment(doc *ast.CommentGroup, name string) bool {
	if doc == nil || len(doc.List) == 0 {
		return false
	}

	comment := strings.TrimSpace(doc.Text())
	firstLine := comment
	if lineEnd := strings.IndexByte(comment, '\n'); lineEnd >= 0 {
		firstLine = comment[:lineEnd]
	}
	if !strings.HasPrefix(firstLine, name) {
		return false
	}

	rawRemainder := strings.TrimPrefix(firstLine, name)
	if rawRemainder != "" {
		firstRemainder, _ := utf8.DecodeRuneInString(rawRemainder)
		if !unicode.IsSpace(firstRemainder) && firstRemainder != ':' && firstRemainder != '：' {
			return false
		}
	}

	remainder := strings.TrimSpace(strings.TrimPrefix(comment, name))
	remainder = strings.TrimSpace(strings.TrimPrefix(remainder, ":"))
	remainder = strings.TrimSpace(strings.TrimPrefix(remainder, "："))
	if remainder == "" {
		return false
	}
	return strings.ContainsFunc(remainder, func(r rune) bool {
		return unicode.Is(unicode.Han, r)
	})
}
