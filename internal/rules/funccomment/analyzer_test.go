package funccomment

import "testing"

func TestAnalyzerName(t *testing.T) {
	if Analyzer == nil || Analyzer.Name != "funccomment" {
		t.Fatal("invalid function analyzer")
	}
}
