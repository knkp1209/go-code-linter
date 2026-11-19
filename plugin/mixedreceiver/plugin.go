package plugin

import (
	"github.com/knkp1209/go-code-linter/internal/rules/mixedreceiver"
	plugininternal "github.com/knkp1209/go-code-linter/plugin/internal"
)

func init() {
	plugininternal.Register("mixedreceiver", mixedreceiver.Analyzer)
}
