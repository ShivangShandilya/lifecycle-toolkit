package analysis

import (
	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

type AnalysisWorkerPool struct {
	*metricsapi.Analysis
	Objectives []metricsapi.Objective
	client.Client
	Log        logr.Logger
	numWorkers int
	numJobs    int
	jobs       chan metricsapi.Objective
	results    chan providers.AnalysisResult
}

func NewWorkerPool(analysis *metricsapi.Analysis, definition *metricsapi.AnalysisDefinition, numWorkers int, c client.Client, log logr.Logger) AnalysisWorkerPool {
	numJobs := len(definition.Spec.Objectives)
	return AnalysisWorkerPool{
		Analysis:   analysis,
		Objectives: definition.Spec.Objectives,
		Client:     c,
		Log:        log,
		numWorkers: numWorkers,
		numJobs:    numJobs,
		jobs:       make(chan metricsapi.Objective, numJobs),
		results:    make(chan providers.AnalysisResult, numJobs),
	}
}

func (aw AnalysisWorkerPool) dispatchObjectives(ctx context.Context) {
	for w := 1; w <= aw.numWorkers; w++ {
		go aw.queryWorker(ctx, w)
	}

	for _, obj := range aw.Objectives {
		aw.jobs <- obj
	}
	close(aw.jobs)
}

func (aw AnalysisWorkerPool) collectAnalysisResults() {
	for a := 1; a <= aw.numJobs; a++ {
		<-aw.results
	}
}

func generateQuery(query string, selectors map[string]string) string {
	for key, value := range selectors {
		query = strings.Replace(query, "$"+strings.ToUpper(key), value, -1)
	}
	return query
}

func (aw AnalysisWorkerPool) queryWorker(ctx context.Context, id int) {
	for j := range aw.jobs {
		aw.Log.Info("worker", "id:", id, "started  job:", j.AnalysisValueTemplateRef.Name)

		template := &metricsapi.AnalysisValueTemplate{}
		err := aw.Client.Get(ctx,
			types.NamespacedName{
				Name:      j.AnalysisValueTemplateRef.Name,
				Namespace: j.AnalysisValueTemplateRef.Namespace},
			template,
		)

		if err != nil {
			aw.results <- providers.AnalysisResult{Objective: &j, Result: providers.Result{Err: err.Error()}}
		}

		providerRef := &metricsapi.KeptnMetricsProvider{}
		err = aw.Client.Get(ctx,
			types.NamespacedName{
				Name:      j.AnalysisValueTemplateRef.Name,
				Namespace: j.AnalysisValueTemplateRef.Namespace},
			providerRef,
		)

		if err != nil {
			aw.Log.Error(err, "Failed to get Provider")
			aw.results <- providers.AnalysisResult{Objective: &j, Result: providers.Result{Err: err.Error()}}
		}

		provider, err2 := providers.NewProvider(providerRef.GetType(), aw.Log, aw.Client)
		if err2 != nil {
			aw.Log.Error(err2, "Failed to get the correct Provider")
			aw.results <- providers.AnalysisResult{Objective: &j, Result: providers.Result{Err: err2.Error()}}
		}

		templatedQuery := generateQuery(template.Spec.Query, aw.Analysis.Spec.Args)
		// TODO implement
		result := provider.RunAnalysis(ctx, templatedQuery, providerRef)

		aw.Log.Info("worker", "id:", id, "finished job:", j.AnalysisValueTemplateRef.Name)
		aw.results <- result
	}
}
