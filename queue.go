package queue

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kyani-inc/queue/file"
	"github.com/kyani-inc/queue/local"
	"github.com/kyani-inc/queue/sqs"
)

type Queue interface {
	Next(queue string) (messageID, result string, err error)
	Append(queue, msg string) error
	Complete(queue, messageID string) error
}

func Local() Queue {
	return local.New()
}

func SQS(secret, access, url, region string, messageTimeout int) Queue {
	return sqs.New(sqs.Config{
		Secret:         secret,
		Access:         access,
		URL:            url,
		Region:         region,
		MessageTimeout: messageTimeout,
	})
}

func SQSNoCreds(url, region string, messageTimeout int) Queue {
	return sqs.NewNoCreds(url, region, messageTimeout)
}

// SQSWithSession returns an instance of SQS using a passed session
func SQSWithSession(url, region string, timeout int, sess *session.Session) Queue {
	return sqs.NewWithSession(url, region, timeout, sess)
}

func File(path string) Queue {
	return file.New(path)
}
