package plugin

import (
	"github.com/knkp1209/go-code-linter/internal/rules/structcomment"
	plugininternal "github.com/knkp1209/go-code-linter/plugin/internal"
)

func init() {
	plugininternal.Register("structcomment", structcomment.Analyzer)
}
