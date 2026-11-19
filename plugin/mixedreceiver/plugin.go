package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	codelint "github.com/knkp1209/go-code-linter"
	"golang.org/x/tools/go/analysis"
)

type linterPlugin struct{}

func init() { register.Plugin("mixedreceiver", New) }

func New(any) (register.LinterPlugin, error) { return &linterPlugin{}, nil }

func (*linterPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{codelint.MixedReceiverAnalyzer}, nil
}

func (*linterPlugin) GetLoadMode() string { return register.LoadModeSyntax }
