package logrus_cloudwatchlogs

import (
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

func TestHook(t *testing.T) {
	a := assertions.New(t)

	key := os.Getenv("AWS_ACCESS_KEY")
	secret := os.Getenv("AWS_SECRET_KEY")
	group := os.Getenv("AWS_CLOUDWATCHLOGS_GROUP_NAME")
	stream := os.Getenv("AWS_CLOUDWATCHLOGS_STREAM_NAME")

	// logs.us-east-1.amazonaws.com
	cred := credentials.NewStaticCredentials(key, secret, "")
	cfg := aws.NewConfig().WithRegion("us-east-1").WithCredentials(cred)

	hook, err := NewHook(group, stream, cfg)
	a.So(err, should.BeNil)
	a.So(hook, should.NotBeNil)

	l := logrus.New()
	l.Hooks.Add(hook)

	for _, level := range hook.Levels() {
		if len(l.Hooks[level]) != 1 {
			t.Errorf("CloudWatchLogs hook was not added. The length of l.Hooks[%v]: %v", level, len(l.Hooks[level]))
		}
	}
}
