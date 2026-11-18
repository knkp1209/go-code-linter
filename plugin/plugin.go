package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	funccomment "github.com/knkp1209/function-comment-linter"
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
		funccomment.Analyzer,
	}, nil
}

func (p *functionCommentPlugin) GetLoadMode() string {
	return register.LoadModeSyntax
}
