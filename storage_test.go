package main

import (
	"bytes"
	"io"
	"testing"
)

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)

	data := []byte("some jpg bytes")
	key := "moms_best_picture"

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}
	r, err := s.read(key)
	if err != nil {
		t.Error(err)
	}
	b, err := io.ReadAll(r)
	if err != nil {
		t.Error(err)
	}

	if string(b) != string(data) {
		t.Errorf("want %s have %s", string(data), string(b))
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
	if PathKey.Filename != expectedOriginKey {
		t.Errorf("have %s want %s", PathKey.Filename, expectedOriginKey)
	}
}

func TestDeleteKey(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "moms_best_picture"
	data := []byte("some jpg bytes")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}
	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}
