package logrus_cloudwatchlogs

import (
	"os"
	"sync"
	"testing"

	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/sirupsen/logrus"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

func TestHook(t *testing.T) {
	a := assertions.New(t)

	key := os.Getenv("AWS_ACCESS_KEY")
	secret := os.Getenv("AWS_SECRET_KEY")
	group := os.Getenv("AWS_CLOUDWATCHLOGS_GROUP_NAME")
	stream := os.Getenv("AWS_CLOUDWATCHLOGS_STREAM_NAME")

	if key == "" {
		t.Skip("skipping test; AWS_ACCESS_KEY not set")
	}
	if secret == "" {
		t.Skip("skipping test; AWS_SECRET_KEY not set")
	}
	if group == "" {
		t.Skip("skipping test; AWS_CLOUDWATCHLOGS_GROUP_NAME not set")
	}
	if stream == "" {
		t.Skip("skipping test; AWS_CLOUDWATCHLOGS_STREAM_NAME not set")
	}

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

func TestConcurrentHook(t *testing.T) {
	a := assertions.New(t)

	key := os.Getenv("AWS_ACCESS_KEY")
	secret := os.Getenv("AWS_SECRET_KEY")
	group := os.Getenv("AWS_CLOUDWATCHLOGS_GROUP_NAME")
	stream := os.Getenv("AWS_CLOUDWATCHLOGS_STREAM_NAME")

	if key == "" {
		t.Skip("skipping test; AWS_ACCESS_KEY not set")
	}
	if secret == "" {
		t.Skip("skipping test; AWS_SECRET_KEY not set")
	}
	if group == "" {
		t.Skip("skipping test; AWS_CLOUDWATCHLOGS_GROUP_NAME not set")
	}
	if stream == "" {
		t.Skip("skipping test; AWS_CLOUDWATCHLOGS_STREAM_NAME not set")
	}

	// logs.us-east-1.amazonaws.com
	cred := credentials.NewStaticCredentials(key, secret, "")
	cfg := aws.NewConfig().WithRegion("us-east-1").WithCredentials(cred)

	hook, err := NewHook(group, stream, cfg)
	a.So(err, should.BeNil)
	a.So(hook, should.NotBeNil)

	l := logrus.New()
	l.Hooks.Add(hook)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := hook.Write([]byte("TestConcurrentHook"))
			a.So(err, should.BeNil)
		}()
	}

	wg.Wait()
}

func TestBatching(t *testing.T) {
	a := assertions.New(t)

	key := os.Getenv("AWS_ACCESS_KEY")
	secret := os.Getenv("AWS_SECRET_KEY")
	group := os.Getenv("AWS_CLOUDWATCHLOGS_GROUP_NAME")
	stream := os.Getenv("AWS_CLOUDWATCHLOGS_STREAM_NAME")

	if key == "" {
		t.Skip("skipping test; AWS_ACCESS_KEY not set")
	}
	if secret == "" {
		t.Skip("skipping test; AWS_SECRET_KEY not set")
	}
	if group == "" {
		t.Skip("skipping test; AWS_CLOUDWATCHLOGS_GROUP_NAME not set")
	}
	if stream == "" {
		t.Skip("skipping test; AWS_CLOUDWATCHLOGS_STREAM_NAME not set")
	}

	// logs.us-east-1.amazonaws.com
	cred := credentials.NewStaticCredentials(key, secret, "")
	cfg := aws.NewConfig().WithRegion("us-east-1").WithCredentials(cred)

	hook, err := NewBatchingHook(group, stream, cfg, 100*time.Millisecond)
	a.So(err, should.BeNil)
	a.So(hook, should.NotBeNil)

	l := logrus.New()
	l.Hooks.Add(hook)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := hook.Write([]byte("TestConcurrentHook"))
			a.So(err, should.BeNil)
		}()
	}

	wg.Wait()
}
