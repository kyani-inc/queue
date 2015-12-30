package local

import (
	"strconv"
)

var queue = []string{}
var counter = 0

type LocalQueue struct{}

func New() LocalQueue {
	return LocalQueue{}
}

func (l LocalQueue) Next(queueName string) (id string, msg string, err error) {
	if len(queue) < 1 {
		return id, msg, nil
	}

	msg, queue = queue[0], queue[1:]

	counter++
	return strconv.Itoa(counter), msg, nil
}

func (l LocalQueue) Append(queueName, msg string) error {
	queue = append(queue, msg)

	return nil
}

func (l LocalQueue) Complete(queueName, messageID string) error {
	return nil
}
