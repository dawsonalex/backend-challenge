package aggregator

import (
	"sort"
	"testing"
	"time"
)

func TestAddInput(t *testing.T) {
	agg := New()

	chan1 := make(chan Operation)
	chan2 := make(chan Operation)
	agg.AddInput(chan1)
	agg.AddInput(chan2)

	agg.Start()

	chan2 <- Operation{
		Filename: "file2.txt",
		OpType:   AddOperation,
	}
	chan1 <- Operation{
		Filename: "file1.txt",
		OpType:   AddOperation,
	}

	time.Sleep(time.Millisecond * 1)
	files := agg.FileList()
	expected := []string{"file1.txt", "file2.txt"}
	if !equals(files, expected) {
		t.Errorf("Expected %v, got %v", expected, files)
	}
}

func equals(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Strings(a)
	sort.Strings(b)

	for _, v := range a {
		i := sort.SearchStrings(b, v)
		if i > len(b) || b[i] != v {
			return false
		}
	}
	return true
}

func Test_Aggregator_Remove(t *testing.T) {

}
