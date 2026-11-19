package mixedreceiver

import "testing"

func TestAnalyzerName(t *testing.T) {
	if Analyzer == nil || Analyzer.Name != "mixedreceiver" {
		t.Fatal("invalid receiver analyzer")
	}
}
