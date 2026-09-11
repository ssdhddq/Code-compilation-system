package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	TaskDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "codeprocessor",
			Name:      "task_duration_seconds",
			Help:      "Время выполнения задачи (в секундах)",
			Buckets:   []float64{0.1, 0.5, 1, 2, 5, 10, 30, 60},
		},
		[]string{"translator", "status"},
	)
	TranslationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "codeprocessor",
			Name:      "translations_total",
			Help:      "Общее число запусков кода по трансляторам",
		},
		[]string{"translator", "status"},
	)
	TasksInProgress = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "codeprocessor",
			Name:      "tasks_in_progress",
			Help:      "Текущее количество обрабатываемых задач",
		},
	)
)

func ObserveTask(translator, status string, duration time.Duration) {
	seconds := duration.Seconds()
	TaskDuration.WithLabelValues(translator, status).Observe(seconds)
	TranslationsTotal.WithLabelValues(translator, status).Inc()
}
