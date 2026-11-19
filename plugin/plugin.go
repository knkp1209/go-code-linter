package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	codelint "github.com/knkp1209/go-code-linter"
	"golang.org/x/tools/go/analysis"
)

type functionCommentPlugin struct{}

func init() {
	register.Plugin("funccomment", New)
}

func New(settings any) (register.LinterPlugin, error) {
	return &functionCommentPlugin{}, nil
}

func (p *functionCommentPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		codelint.Analyzer,
		codelint.ConstAnalyzer,
		codelint.StructAnalyzer,
		codelint.MixedReceiverAnalyzer,
	}, nil
}

func (p *functionCommentPlugin) GetLoadMode() string {
	return register.LoadModeSyntax
}
