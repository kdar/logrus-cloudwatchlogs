package logrus_cloudwatchlogs

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/kdar/gtest"
	"github.com/sirupsen/logrus"
	"github.com/smartystreets/assertions/should"
)

//Example is a struct that implements Marshaler. It marshals as a map where "Content" is set to the value of A, and B is discarded
type Example struct {
	A, B string
}

func (e Example) MarshalLog() map[string]interface{} {
	out := make(map[string]interface{})
	out["Content"] = e.A
	return out
}

func TestProdFormatter(t *testing.T) {
	a := gtest.New(t)

	l := logrus.New()
	l.Out = ioutil.Discard
	prodFormatter := NewProdFormatter()
	l.Formatter = prodFormatter

	buf := &bytes.Buffer{}
	l.Hooks.Add(NewWriterHook(buf))

	l.WithFields(logrus.Fields{
		"event":       "testevent",
		"topic":       "testtopic",
		"key":         "testkey",
		"marshaltest": Example{A: "hello", B: "this gets dropped"},
	}).Info("Some event")

	// split := bytes.SplitN(buf.Bytes(), []byte(" "), 2)
	// a.So(len(split), should.Equal, 2).ElseFatal()

	var v map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &v)
	a.So(err, should.BeNil).ElseFatal()

	a.So(v["event"], should.Equal, "testevent").ElseFatal()
	a.So(v["topic"], should.Equal, "testtopic").ElseFatal()
	a.So(v["key"], should.Equal, "testkey").ElseFatal()
	a.So(v["msg"], should.Equal, "Some event").ElseFatal()
	a.So(v["marshaltest"].(map[string]interface{})["Content"], should.Equal, "hello").ElseFatal()
}

func TestProdFormatterCaller(t *testing.T) {
	a := gtest.New(t)

	l := logrus.New()
	l.Out = ioutil.Discard
	prodFormatter := NewProdFormatter()
	l.Formatter = prodFormatter

	buf := &bytes.Buffer{}
	l.Hooks.Add(NewWriterHook(buf))

	modifiers := []func(*logrus.Logger) *logrus.Entry{
		func(logger *logrus.Logger) *logrus.Entry {
			return logrus.NewEntry(logger)
		},
		func(logger *logrus.Logger) *logrus.Entry {
			return logger.WithFields(logrus.Fields{
				"event": "testevent",
				"topic": "testtopic",
				"key":   "testkey",
			})
		},
		func(logger *logrus.Logger) *logrus.Entry {
			return logger.WithFields(logrus.Fields{
				"event": "testevent",
				"topic": "testtopic",
				"key":   "testkey",
			}).WithFields(logrus.Fields{
				"event1": "testevent",
				"topic1": "testtopic",
				"key1":   "testkey",
			}).WithFields(logrus.Fields{
				"event2": "testevent",
				"topic2": "testtopic",
				"key2":   "testkey",
			})
		},
		func(logger *logrus.Logger) *logrus.Entry {
			return logger.WithError(errors.New("some err"))
		},
	}

	for i, mod := range modifiers {
		buf.Reset()
		modl := mod(l)

		modl.Info("Some event")
		// split := bytes.SplitN(buf.Bytes(), []byte(" "), 2)
		// a.So(len(split), should.Equal, 2).Else(func(msg string) {
		// 	t.Fatalf("\nfailed at index %d\n%s", i, msg)
		// })

		var v map[string]interface{}
		err := json.Unmarshal(buf.Bytes(), &v)
		a.So(err, should.BeNil).Else(func(msg string) {
			t.Fatalf("\nfailed at index %d\n%s", i, msg)
		})

		// extra := v[prodFormatter.extraKey].(map[string]interface{})
		// a.So(extra["file"], should.Equal, "formatter_test.go").Else(func(msg string) {
		// 	t.Fatalf("\nfailed at index %d\n%s", i, msg)
		// })
		// a.So(extra["func"], should.Equal, "github.com/kdar/logrus-cloudwatchlogs.TestProdFormatterCaller").Else(func(msg string) {
		// 	t.Fatalf("\nfailed at index %d\n%s", i, msg)
		// })
	}
}

func TestProdFormatterHTTPRequest(t *testing.T) {
	a := gtest.New(t)

	req, err := http.NewRequest("GET", "http://rxmanagement.net", nil)
	a.So(err, should.BeNil).ElseFatal()

	buf := &bytes.Buffer{}

	l := logrus.New()
	l.Out = buf
	prodFormatter := NewProdFormatter(HTTPRequest("request"))
	l.Formatter = prodFormatter

	l.WithField("request", req).Print("with http request")

	var v struct {
		Request struct {
			Host   string
			Method string
		}
	}
	err = json.Unmarshal(buf.Bytes(), &v)
	a.So(err, should.BeNil).ElseFatal()

	a.So(v.Request.Host, should.Equal, "rxmanagement.net").ElseFatal()
	a.So(v.Request.Method, should.Equal, "GET").ElseFatal()
}
