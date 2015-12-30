package file

import (
	"encoding/json"
	"errors"
	"strconv"
	"sync"

	"gopkg.in/kyani-inc/storage.v1"
)

var lock sync.Mutex
var queues = map[string]storage.Storage{}

type row struct {
	ID      string `json:"i"`
	Msg     string `json:"m"`
	Pending bool   `json:"p"`
	Done    bool   `json:"d"`
}

type File struct {
	path string
}

func New(path string) File {
	lock.Lock()
	defer lock.Unlock()

	_, ok := queues[path]

	if !ok {
		queues[path], _ = storage.Folder(path)
	}

	return File{path: path}
}

func (f File) Next(queue string) (id string, msg string, err error) {
	lock.Lock()
	defer lock.Unlock()

	q, err := f.get(queue)

	if err != nil {
		return "", "", err
	}

	for i, node := range q {
		if !node.Done && !node.Pending {
			q[i].Pending = true

			if err := f.put(queue, q); err != nil {
				return "", "", err
			}

			return node.ID, node.Msg, err
		}
	}

	return "", "", nil
}

func (f File) Append(queue, msg string) error {
	lock.Lock()
	defer lock.Unlock()

	q, err := f.get(queue)

	if err != nil {
		return err
	}

	q = append(q, row{
		ID:  strconv.Itoa(len(q) + 1),
		Msg: msg,
	})

	return f.put(queue, q)
}

func (f File) Complete(queue, id string) error {
	lock.Lock()
	defer lock.Unlock()

	q, err := f.get(queue)

	if err != nil {
		return err
	}

	for i, node := range q {
		if node.ID == id {
			q[i].Done = true
			q[i].Pending = false

			if err := f.put(queue, q); err != nil {
				return err
			}

			return nil
		}
	}

	return nil
}

func (f File) Flush(queue string) error {
	lock.Lock()
	defer lock.Unlock()

	return f.put(queue, []row{})
}

func (f File) get(key string) ([]row, error) {
	q := []row{}

	store, ok := queues[f.path]

	if !ok {
		return q, errors.New("file queue: path missing")
	}

	data := store.Get(key)

	if len(data) < 1 {
		return q, nil
	}

	err := json.Unmarshal(data, &q)

	return q, err
}

func (f File) put(key string, q []row) error {
	data, err := json.Marshal(q)

	if err != nil {
		return err
	}

	store, ok := queues[f.path]

	if !ok {
		return errors.New("file queue: path missing")
	}

	return store.Put(key, data)
}
