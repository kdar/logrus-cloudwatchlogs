package logrus_cloudwatchlogs

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/kdar/gtest"
	"github.com/sirupsen/logrus"
	"github.com/smartystreets/assertions/should"
)

func TestDevFormatterHTTPRequest(t *testing.T) {
	a := gtest.New(t)

	req, err := http.NewRequest("GET", "http://rxmanagement.net", nil)
	a.So(err, should.BeNil).ElseFatal()

	buf := &bytes.Buffer{}

	l := logrus.New()
	l.Out = buf
	formatter := &DevFormatter{}
	formatter.HTTPRequestKey = "request"
	l.Formatter = formatter

	l.WithFields(logrus.Fields{
		"request": req,
	}).Print("with http request")

	a.So(buf.String(), should.NotContainSubstring, "rxmanagement").ElseFatal()
}
