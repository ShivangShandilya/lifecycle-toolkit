/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package analysis

import (
	"context"
	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AnalysisReconciler reconciles a Analysis object
type AnalysisReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	Log        logr.Logger
	numWorkers int //maybe 2 or 4 as def
}

//+kubebuilder:rbac:groups=metrics.keptn.sh,resources=analyses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=metrics.keptn.sh,resources=analyses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=metrics.keptn.sh,resources=analyses/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Analysis object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the usea.
//
// For more details, check Reconcile and its AnalysisResult here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (a *AnalysisReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	a.Log.Info("Reconciling Analysis")
	analysis := &metricsapi.Analysis{}

	//retrieve analysis
	if err := a.Client.Get(ctx, req.NamespacedName, analysis); err != nil {
		if errors.IsNotFound(err) {
			// taking down all associated K8s resources is handled by K8s
			a.Log.Info("Analysis resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		a.Log.Error(err, "Failed to get the Analysis")
		return ctrl.Result{}, nil
	}

	//find AnalysisDefinition to have the collection of Objectives
	analysisDef := &metricsapi.AnalysisDefinition{}
	err := a.Client.Get(ctx,
		types.NamespacedName{
			Name:      analysis.Spec.AnalysisDefinition.Name,
			Namespace: analysis.Spec.AnalysisDefinition.Namespace},
		analysisDef,
	)

	if err != nil {
		if errors.IsNotFound(err) {
			a.Log.Info(err.Error() + ", ignoring error since object must be deleted")
			return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
		}
		a.Log.Error(err, "Failed to retrieve the AnalysisDefinition")
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	//create multiple workers handling the Objectives
	wp := NewWorkerPool(analysis, analysisDef, a.numWorkers, a.Client, a.Log)

	wp.dispatchObjectives(ctx)
	wp.collectAnalysisResults()

	//TODO compute score and update status

	//TODO update final status making sure analysis did not change!!
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Managea.
func (a *AnalysisReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&metricsapi.Analysis{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Complete(a)
}
