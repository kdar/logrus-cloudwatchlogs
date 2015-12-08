package logrus_cloudwatchlogs

import "github.com/Sirupsen/logrus"

type DevFormatter struct {
	HTTPRequestKey string
	*logrus.TextFormatter
}

func (f *DevFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if _, ok := entry.Data[f.HTTPRequestKey]; ok {
		delete(entry.Data, f.HTTPRequestKey)
	}

	if f.TextFormatter == nil {
		f.TextFormatter = &logrus.TextFormatter{}
	}

	return f.TextFormatter.Format(entry)
}
