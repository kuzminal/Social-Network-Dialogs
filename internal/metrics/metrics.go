package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	SendMessage = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_send_message",
		Help:    "A summary of the HTTP request durations in seconds.",
		Buckets: prometheus.LinearBuckets(0, 10, 20),
	},
	)
	GetMessage = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_get_message",
		Help:    "A summary of the HTTP request durations in seconds.",
		Buckets: prometheus.LinearBuckets(0, 10, 20),
	},
	)
	MarkAsRead = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_mark_as_read",
		Help:    "A summary of the HTTP request durations in seconds.",
		Buckets: prometheus.LinearBuckets(0, 10, 20),
	},
	)
)
