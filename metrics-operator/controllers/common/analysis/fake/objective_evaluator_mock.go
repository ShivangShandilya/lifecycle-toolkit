// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package fake

import (
	"github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"sync"
)

// IObjectiveEvaluatorMock is a mock implementation of analysis.IObjectiveEvaluator.
//
//	func TestSomethingThatUsesIObjectiveEvaluator(t *testing.T) {
//
//		// make and configure a mocked analysis.IObjectiveEvaluator
//		mockedIObjectiveEvaluator := &IObjectiveEvaluatorMock{
//			EvaluateFunc: func(values map[string]string, objective v1alpha3.Objective) v1alpha3.ObjectiveResult {
//				panic("mock out the Evaluate method")
//			},
//		}
//
//		// use mockedIObjectiveEvaluator in code that requires analysis.IObjectiveEvaluator
//		// and then make assertions.
//
//	}
type IObjectiveEvaluatorMock struct {
	// EvaluateFunc mocks the Evaluate method.
	EvaluateFunc func(values map[string]string, objective v1alpha3.Objective) v1alpha3.ObjectiveResult

	// calls tracks calls to the methods.
	calls struct {
		// Evaluate holds details about calls to the Evaluate method.
		Evaluate []struct {
			// Values is the values argument value.
			Values map[string]string
			// Objective is the objective argument value.
			Objective v1alpha3.Objective
		}
	}
	lockEvaluate sync.RWMutex
}

// Evaluate calls EvaluateFunc.
func (mock *IObjectiveEvaluatorMock) Evaluate(values map[string]string, objective v1alpha3.Objective) v1alpha3.ObjectiveResult {
	if mock.EvaluateFunc == nil {
		panic("IObjectiveEvaluatorMock.EvaluateFunc: method is nil but IObjectiveEvaluator.Evaluate was just called")
	}
	callInfo := struct {
		Values    map[string]string
		Objective v1alpha3.Objective
	}{
		Values:    values,
		Objective: objective,
	}
	mock.lockEvaluate.Lock()
	mock.calls.Evaluate = append(mock.calls.Evaluate, callInfo)
	mock.lockEvaluate.Unlock()
	return mock.EvaluateFunc(values, objective)
}

// EvaluateCalls gets all the calls that were made to Evaluate.
// Check the length with:
//
//	len(mockedIObjectiveEvaluator.EvaluateCalls())
func (mock *IObjectiveEvaluatorMock) EvaluateCalls() []struct {
	Values    map[string]string
	Objective v1alpha3.Objective
} {
	var calls []struct {
		Values    map[string]string
		Objective v1alpha3.Objective
	}
	mock.lockEvaluate.RLock()
	calls = mock.calls.Evaluate
	mock.lockEvaluate.RUnlock()
	return calls
}
