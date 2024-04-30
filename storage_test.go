package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	return NewStore(opts)
}

func teardown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Fatal(err)
	}
}

func TestStore(t *testing.T) {

	s := newStore()

	defer teardown(t, s)

	data := []byte("some jpg bytes")
	key := "fooanbar"

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if ok := s.Has(key); !ok {
		t.Errorf("expected to have key %s", key)
	}

	r, err := s.read(key)
	if err != nil {
		t.Error(err)
	}
	b, err := io.ReadAll(r)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(b))
	if string(b) != string(data) {
		t.Errorf("want %s have %s", string(data), string(b))
	}
}

func TestPathTransformFunc(t *testing.T) {
	key := "fooandbar"
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
