package structcomment

import "testing"

func TestAnalyzerName(t *testing.T) {
	if Analyzer == nil || Analyzer.Name != "structcomment" {
		t.Fatal("invalid struct analyzer")
	}
}
