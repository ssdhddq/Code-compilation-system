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

var (
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "app",
			Name:      "http_requests_total",
			Help:      "Общее число HTTP-запросов",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "app",
			Name:      "http_request_duration_seconds",
			Help:      "Время обработки HTTP-запроса (в секундах)",
			Buckets:   []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2, 5},
		},
		[]string{"method", "path"},
	)
	DBTasksTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "app",
			Name:      "db_operations_total",
			Help:      "Общее число операций с БД",
		},
		[]string{"operation", "status"},
	)
)

func ObserveHTTP(method, path, status string, duration time.Duration) {
	HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
	HTTPRequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
}

func ObserveDB(operation string, err error) {
	status := "success"
	if err != nil {
		status = "error"
	}
	DBTasksTotal.WithLabelValues(operation, status).Inc()
}
