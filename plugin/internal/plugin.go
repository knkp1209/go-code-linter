package internal

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

type linterPlugin struct {
	analyzer *analysis.Analyzer
}

func Register(name string, analyzer *analysis.Analyzer) {
	register.Plugin(name, newPlugin(analyzer))
}

func newPlugin(analyzer *analysis.Analyzer) register.NewPlugin {
	return func(any) (register.LinterPlugin, error) {
		return &linterPlugin{analyzer: analyzer}, nil
	}
}

func (p *linterPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{p.analyzer}, nil
}

func (*linterPlugin) GetLoadMode() string {
	return register.LoadModeSyntax
}
