package sqs

type Config struct {
	Secret         string
	Access         string
	URL            string
	Region         string
	MessageTimeout int
}
