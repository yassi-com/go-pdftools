package pdftk

import (
	"reflect"
	"testing"
)

func TestCatArgs(t *testing.T) {
	tests := []struct {
		name       string
		fileNames  InputFileMap
		pageRanges []PageRange
		want       []string
	}{
		{
			name:      "no page ranges concatenates every input",
			fileNames: NewInputFileMap("first.pdf", "second.pdf"),
			want:      []string{"A=first.pdf", "B=second.pdf", "cat", "output", "-"},
		},
		{
			name:      "page ranges are emitted in request order",
			fileNames: InputFileMap{"A": "first.pdf"},
			pageRanges: []PageRange{
				{FileHandleName: "A", BeginPage: 2, EndPage: 4},
				{FileHandleName: "A", BeginPage: 1, EndPage: 1},
			},
			want: []string{"A=first.pdf", "cat", "A2-4", "A1", "output", "-"},
		},
		{
			name:      "page ranges span inputs",
			fileNames: NewInputFileMap("first.pdf", "second.pdf"),
			pageRanges: []PageRange{
				{FileHandleName: "B", Qualifier: Odd},
				{FileHandleName: "A", BeginPage: 4, Rotation: East},
			},
			want: []string{"A=first.pdf", "B=second.pdf", "cat", "Bodd", "A4-endeast", "output", "-"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := catArgs(tt.fileNames, tt.pageRanges); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("catArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
