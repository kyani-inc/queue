package file

import (
	"io/ioutil"
	"sync"
	"testing"
)

func TestFile(t *testing.T) {
	p, err := ioutil.TempDir("", "queue-file-test")

	if err != nil {
		t.Fatal(err.Error())
	}

	q := New(p)
	q.Flush("queue_test")

	if err := q.Append("queue_test", "first"); err != nil {
		t.Errorf("error on append: %s", err.Error())
		return
	}

	id, msg, err := q.Next("queue_test")

	if err != nil {
		t.Errorf("error on next: %s", err.Error())
		return
	}

	if msg != "first" {
		t.Errorf("expected 'first' but got '%s'", msg)
		return
	}

	q.Complete("queue_test", id)
}

func TestFileRace(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func(wg *sync.WaitGroup, t *testing.T) {
		defer wg.Done()
		TestFile(t)
	}(&wg, t)

	go func(wg *sync.WaitGroup, t *testing.T) {
		defer wg.Done()
		TestFile(t)
	}(&wg, t)

	wg.Wait()
}
