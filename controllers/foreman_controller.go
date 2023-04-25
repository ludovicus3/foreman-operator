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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	api "github.com/ludovicus3/foreman-operator/api/v1alpha1"
)

// ForemanReconciler reconciles a Foreman object
type ForemanReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=foreman.theforeman.org,resources=foremen,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=foreman.theforeman.org,resources=foremen/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=foreman.theforeman.org,resources=foremen/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ForemanReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	foreman := &api.Foreman{}

	if err := r.Get(ctx, req.NamespacedName, foreman); err != nil {
		if errors.IsNotFound(err) {
			log.Info("Foreman resource not found. Ignoring.")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if foreman.Status.Conditions == nil || len(foreman.Status.Conditions) == 0 {
		err := r.StatusUpdate(ctx, req, log, foreman, api.ConditionPending,
			metav1.ConditionFalse,
			api.ReasonReqNotMet,
			"Foreman progressing")
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ForemanReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&api.Foreman{}).
		Complete(r)
}

// StatusUpdate update the status to reflect the current state.
func (r *ForemanReconciler) StatusUpdate(ctx context.Context, req ctrl.Request, log logr.Logger, foreman *api.Foreman, condition string, status metav1.ConditionStatus, reason string, message string) error {
	meta.SetStatusCondition(&foreman.Status.Conditions, metav1.Condition{
		Status:             status,
		Reason:             reason,
		Message:            message,
		Type:               condition,
		ObservedGeneration: foreman.GetGeneration(),
	})
	if err := r.Status().Update(ctx, foreman); err != nil {
		log.Error(err, "Failed to update Foreman status")
		return err
	}

	if err := r.Get(ctx, req.NamespacedName, foreman); err != nil {
		log.Error(err, "Failed to re-fetch Foreman")
		return err
	}
	return nil
}
