package ycmonitoringgo

import (
	"net/http"
)

type options struct {
	url    string
	client *http.Client
	logger Logger
}

func defaultOptions() *options {
	return &options{
		url:    "https://monitoring.api.cloud.yandex.net/monitoring",
		client: http.DefaultClient,
	}
}

type Option func(opts *options)

type Logger interface {
	Error(msg string, args ...any)
}

func WithHttpClient(client *http.Client) Option {
	return func(opts *options) {
		opts.client = client
	}
}

func WithLogger(logger Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func WithUrl(url string) Option {
	return func(opts *options) {
		opts.url = url
	}
}
