package local

import (
	"sync"
	"testing"
)

func TestLocal(t *testing.T) {
	q := New()
	qn := "test"

	q.Append(qn, "first")
	q.Append(qn, "second")
	q.Append(qn, "third")

	_, a, _ := q.Next(qn)
	_, b, _ := q.Next(qn)
	_, c, _ := q.Next(qn)

	if a != "first" {
		t.Errorf("expected 'first' got %s", a)
	}

	if b != "second" {
		t.Errorf("expected 'second' got %s", b)
	}

	if c != "third" {
		t.Errorf("expected 'third' got %s", c)
	}
}

func TestLocalRace(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func(wg *sync.WaitGroup, t *testing.T) {
		defer wg.Done()
		TestLocal(t)
	}(&wg, t)

	go func(wg *sync.WaitGroup, t *testing.T) {
		defer wg.Done()
		TestLocal(t)
	}(&wg, t)

	wg.Wait()
}
