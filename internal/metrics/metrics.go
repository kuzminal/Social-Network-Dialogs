package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	SendMessage = prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "http_send_message",
		Help: "A summary of the HTTP request durations in seconds.",
		Objectives: map[float64]float64{
			0.5:  0.05,  // 50th percentile with a max. absolute error of 0.05.
			0.9:  0.01,  // 90th percentile with a max. absolute error of 0.01.
			0.99: 0.001, // 99th percentile with a max. absolute error of 0.001.
		},
	},
	)
	GetMessage = prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "http_get_message",
		Help: "A summary of the HTTP request durations in seconds.",
		Objectives: map[float64]float64{
			0.5:  0.05,  // 50th percentile with a max. absolute error of 0.05.
			0.9:  0.01,  // 90th percentile with a max. absolute error of 0.01.
			0.99: 0.001, // 99th percentile with a max. absolute error of 0.001.
		},
	},
	)
	MarkAsRead = prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "http_mark_as_read",
		Help: "A summary of the HTTP request durations in seconds.",
		Objectives: map[float64]float64{
			0.5:  0.05,  // 50th percentile with a max. absolute error of 0.05.
			0.9:  0.01,  // 90th percentile with a max. absolute error of 0.01.
			0.99: 0.001, // 99th percentile with a max. absolute error of 0.001.
		},
	},
	)
)
