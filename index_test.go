package gostruct

import (
	"fmt"
	"testing"
)

func TestIndex(t *testing.T) {
	index := NewIndex()

	ok1 := index.HasElem("test")
	if ok1 {
		t.Error(fmt.Sprintf("test doesn't eixst in index"))
	}

	index.AddElem("test")
	ok2 := index.HasElem("test")
	if !ok2 {
		t.Error(fmt.Sprintf("test eixsts in index"))
	}

	id := index.GetId("test")
	if id != 0 {
		t.Error(fmt.Sprintf("got %d want 1", id))
	}
}
