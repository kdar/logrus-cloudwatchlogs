package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/kdar/logrus-cloudwatchlogs"
	"github.com/sirupsen/logrus"
)

func main() {
	key := os.Getenv("AWS_ACCESS_KEY")
	secret := os.Getenv("AWS_SECRET_KEY")
	group := os.Getenv("AWS_CLOUDWATCHLOGS_GROUP_NAME")
	stream := os.Getenv("AWS_CLOUDWATCHLOGS_STREAM_NAME")

	// logs.us-east-1.amazonaws.com
	cred := credentials.NewStaticCredentials(key, secret, "")
	cfg := aws.NewConfig().WithRegion("us-east-1").WithCredentials(cred)

	hook, err := logrus_cloudwatchlogs.NewHook(group, stream, cfg)
	if err != nil {
		log.Fatal(err)
	}

	l := logrus.New()
	l.Hooks.Add(hook)
	l.Out = ioutil.Discard
	l.Formatter = logrus_cloudwatchlogs.NewProdFormatter()

	l.WithFields(logrus.Fields{
		"event": "testevent",
		"topic": "testtopic",
		"key":   "testkey",
	}).Fatal("Some fatal event")
}
