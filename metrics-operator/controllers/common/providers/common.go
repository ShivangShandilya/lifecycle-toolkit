package providers

import metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"

const DynatraceProviderType = "dynatrace"
const DynatraceDQLProviderType = "dql"
const PrometheusProviderType = "prometheus"
const DataDogProviderType = "datadog"

type AnalysisResult struct {
	*metricsapi.Objective
	Result
}

type Result struct {
	Value string
	Raw   []byte
	Err   string
}
