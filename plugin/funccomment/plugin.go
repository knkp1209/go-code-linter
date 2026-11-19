package plugin

import (
	codelint "github.com/knkp1209/go-code-linter"
	plugininternal "github.com/knkp1209/go-code-linter/plugin/internal"
)

func init() {
	plugininternal.Register("funccomment", codelint.FunctionCommentAnalyzer)
}
