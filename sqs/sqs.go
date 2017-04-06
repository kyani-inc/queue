package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQS struct {
	cfg        Config
	sqsService *sqs.SQS
}

// New returns an instance of SQS
func New(cfg Config) SQS {
	var instance SQS
	instance.cfg = cfg

	creds := credentials.NewStaticCredentials(cfg.Access, cfg.Secret, "")

	awsConfig := aws.Config{}
	awsConfig.WithCredentials(creds)
	awsConfig.Region = aws.String(cfg.Region)

	session := session.New()
	instance.sqsService = sqs.New(session, &awsConfig)

	return instance
}

// NewUseRole returns an instance of SQS using default credentials, fall back to env...
func NewNoCreds(url, region string, timeout int) SQS {
	awsConfig := aws.Config{}
	awsConfig.Region = aws.String(region)

	session := session.New()
	svc := sqs.New(session, &awsConfig)
	return SQS{
		sqsService: svc,
		cfg: Config{
			URL: url,
			MessageTimeout: timeout,
		},
	}
}

// Next returns the next SQS message.
func (instance SQS) Next(queue string) (messageID, result string, err error) {
	params := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(instance.cfg.URL + "/" + queue),
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(int64(instance.cfg.MessageTimeout)),
	}

	resp, err := instance.sqsService.ReceiveMessage(params)

	if err != nil {
		return "", "", err
	}

	if len(resp.Messages) > 0 {
		return *resp.Messages[0].ReceiptHandle, *resp.Messages[0].Body, nil
	}

	return messageID, "", nil
}

func (instance SQS) Append(queue, msg string) (err error) {
	params := &sqs.SendMessageInput{
		QueueUrl:    aws.String(instance.cfg.URL + "/" + queue),
		MessageBody: aws.String(msg),
	}

	_, err = instance.sqsService.SendMessage(params)
	if err != nil {
		return err
	}

	return nil
}

// Complete lets SQS know that a message was successfully
// processed.
func (instance SQS) Complete(queue, messageID string) (err error) {
	params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(instance.cfg.URL + "/" + queue),
		ReceiptHandle: aws.String(messageID),
	}

	_, err = instance.sqsService.DeleteMessage(params)
	if err != nil {
		return err
	}

	return nil
}
