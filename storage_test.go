package main

import (
	"bytes"
	"testing"
)

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: DefaultPathTransformFunc,
	}
	s := NewStore(opts)

	data := bytes.NewReader([]byte("some jpg bytes"))

	if err := s.WriteStream("myspecialpicture", data); err != nil {
		t.Error(err)
	}
}

func TestPathTransformFunc(t *testing.T) {
	key := "moms_best_picture"
	PathKey := CASPathTransformFunc(key)
	expectedOriginKey := "841e17de0d42f33b85acdb9ede1a47bfa8b9ef9f"
	expectedPathName := "841e1/7de0d/42f33/b85ac/db9ed/e1a47/bfa8b/9ef9f"
	if PathKey.PathName != expectedPathName {
		t.Errorf("hava %s want %s", PathKey.PathName, expectedPathName)
	}
	if PathKey.Original != expectedOriginKey {
		t.Errorf("have %s want %s", PathKey.Original, expectedOriginKey)
	}
}
