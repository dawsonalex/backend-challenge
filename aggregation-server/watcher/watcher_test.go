package watcher

import (
	"testing"

	"github.com/google/uuid"
)

// TestAddNode checks that nodes can be added to the registry and
// their initial file state can be set.
func TestAddNode(t *testing.T) {

	reg := NewRegistry()

	filechan, done, added := reg.AddNode(uuid.New())
	if added {
		filechan <- "file1.txt"
		filechan <- "file2.txt"
		close(filechan)
		<-done
		fileCount := len(reg.ListFiles())
		if fileCount != 2 {
			t.Errorf("expected 2 files, got %d", fileCount)
			return
		}
		return
	}
	t.Error("node not identified as being added to registry.")
}
