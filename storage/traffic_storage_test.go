package storage

import (
	"fmt"
	"testing"
)

func TestStorage_put (t *testing.T) {
	storage := NewVolativeTrafficStorage
	trace := storage.CreateTrace()
	fmt.Printf("trace created successfully with id " + trace.ID)//storage.store(trace)
}



