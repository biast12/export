package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ApiRequests       = promauto.NewCounterVec(counterOpsApi("requests"), []string{"method", "route"})
	ApiTokenExchanges = promauto.NewCounter(counterOpsApi("token_exchanges"))

	RequestsCreated          = promauto.NewCounterVec(counterOpsApi("requests_created"), []string{"type"})
	ArtifactsDownloaded      = promauto.NewCounterVec(counterOpsApi("artifacts_downloaded"), []string{"type"})
	ArtifactsDownloadedBytes = promauto.NewCounterVec(counterOpsApi("artifacts_downloaded_bytes"), []string{"type"})

	RequestsProcessed      = promauto.NewCounterVec(counterOpsWorker("requests_processed"), []string{"type", "status"})
	ArtifactsUploaded      = promauto.NewCounterVec(counterOpsWorker("artifacts_uploaded"), []string{"type"})
	ArtifactsUploadedBytes = promauto.NewCounterVec(counterOpsWorker("artifacts_uploaded_bytes"), []string{"type"})
)

func counterOpsApi(name string) prometheus.CounterOpts {
	return prometheus.CounterOpts{
		Namespace: "tickets",
		Subsystem: "export_api",
		Name:      name,
	}
}

func counterOpsWorker(name string) prometheus.CounterOpts {
	return prometheus.CounterOpts{
		Namespace: "tickets",
		Subsystem: "export_worker",
		Name:      name,
	}
}
