# Queue [![Build Status](https://travis-ci.org/kyani-inc/queue.svg)](https://travis-ci.org/kyani-inc/queue)&nbsp;[![godoc reference](https://godoc.org/github.com/kyani-inc/queue?status.png)](https://godoc.org/github.com/kyani-inc/queue)

Unified queue interface for several different backing technologies. 

- **SQS** ([Amazon SQS](https://aws.amazon.com/sqs/) backed store; production ready)
- **File** (uses a local file for queue; intended for dev only)
- **Local** (uses a memory based queue; intended for dev only)

# Usage and Examples

`go get -u gopkg.in/kyani-inc/queue.v1`

You can write your application to use the queue interface and then choose which backing technology to use
based on environment needs. For example, you could use `SQS` for a production sytem and `Local` for 
your development environment. 



