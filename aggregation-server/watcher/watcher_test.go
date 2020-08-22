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

func TestAddExistingNode(t *testing.T) {

	reg := NewRegistry()

	id := uuid.New()
	reg.AddNode(id)
	reg.AddNode(id)

	nodeCount := len(reg.nodes)
	if nodeCount != 1 {
		t.Errorf("registry should contain 1 nodes, got %d", nodeCount)
	}
}

func TestRemoveNode(t *testing.T) {

	reg := NewRegistry()

	id := uuid.New()
	filechan, done, added := reg.AddNode(id)
	if added {
		filechan <- "file1.txt"
		filechan <- "file2.txt"
		close(filechan)
		<-done

		reg.RemoveNode(id)
		nodeCount := len(reg.nodes)
		if nodeCount != 0 {
			t.Errorf("registry should contain 0 nodes, got %d", nodeCount)
		}
	}
}

func TestRemoveNonExistingNode(t *testing.T) {

	reg := NewRegistry()

	id := uuid.New()
	reg.AddNode(id)
	reg.RemoveNode(uuid.New())

	nodeCount := len(reg.nodes)
	if nodeCount != 1 {
		t.Errorf("registry should contain 1 nodes, got %d", nodeCount)
	}
}

func TestAddOperation(t *testing.T) {
	reg := NewRegistry()

	id := uuid.New()
	reg.AddNode(id)
	op := Operation{
		Type:     "add",
		SeqNo:    1,
		Filename: "file1.txt",
	}
	reg.Node(id).Do(op)
	fileCount := len(reg.ListFiles())
	if fileCount != 1 {
		t.Errorf("expected 1 files, got %d", fileCount)
	}
}

func TestRemoveOperation(t *testing.T) {
	reg := NewRegistry()

	id := uuid.New()
	reg.AddNode(id)
	addOp := Operation{
		Type:     "add",
		SeqNo:    1,
		Filename: "file1.txt",
	}
	reg.Node(id).Do(addOp)

	remOp := Operation{
		Type:     "remove",
		SeqNo:    2,
		Filename: "file1.txt",
	}
	reg.Node(id).Do(remOp)

	fileCount := len(reg.ListFiles())
	if fileCount != 0 {
		t.Errorf("expected 0 files, got %d", fileCount)
	}
}
