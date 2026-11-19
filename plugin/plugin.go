package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	codelint "github.com/knkp1209/go-code-linter"
	"golang.org/x/tools/go/analysis"
)

type analyzerPlugin struct {
	analyzers []*analysis.Analyzer
}

func init() {
	register.Plugin("funccomment", newPlugin(codelint.Analyzer))
	register.Plugin("constcomment", newPlugin(codelint.ConstAnalyzer))
	register.Plugin("structcomment", newPlugin(codelint.StructAnalyzer))
	register.Plugin("mixedreceiver", newPlugin(codelint.MixedReceiverAnalyzer))
}

func newPlugin(analyzer *analysis.Analyzer) func(any) (register.LinterPlugin, error) {
	return func(settings any) (register.LinterPlugin, error) {
		return &analyzerPlugin{analyzers: []*analysis.Analyzer{analyzer}}, nil
	}
}

func (p *analyzerPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return p.analyzers, nil
}

func (p *analyzerPlugin) GetLoadMode() string {
	return register.LoadModeSyntax
}
