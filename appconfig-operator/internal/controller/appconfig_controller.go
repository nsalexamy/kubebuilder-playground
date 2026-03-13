/*
Copyright 2026.

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

package controller

import (
	"context"
	"reflect"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	// nsalexamyv1alpha1 "github.com/nsalexamy/kubebuilder-playground/appconfig-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	nsalexamycomv1alpha1 "github.com/nsalexamy/kubebuilder-playground/appconfig-operator/api/v1alpha1"
)

// AppConfigReconciler reconciles a AppConfig object
type AppConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=nsalexamy.com,resources=appconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=nsalexamy.com,resources=appconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=nsalexamy.com,resources=appconfigs/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AppConfig object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.23.1/pkg/reconcile
func (r *AppConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logf.FromContext(ctx)

	// Fetch the AppConfig resource
	var appConfig nsalexamycomv1alpha1.AppConfig
	err := r.Get(ctx, req.NamespacedName, &appConfig)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Check if ConfigMap already exists
	var configMap corev1.ConfigMap
	err = r.Get(ctx, req.NamespacedName, &configMap)

	if err != nil && apierrors.IsNotFound(err) {

		// Create ConfigMap if it doesn't exist
		cm := corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      appConfig.Name,
				Namespace: appConfig.Namespace,
			},
			Data: appConfig.Spec.Data,
		}

		// Set owner reference
		err = controllerutil.SetControllerReference(&appConfig, &cm, r.Scheme)
		if err != nil {
			return ctrl.Result{}, err
		}

		err = r.Create(ctx, &cm)
		return ctrl.Result{}, err
	}

	// Update ConfigMap if spec changes
	if !reflect.DeepEqual(configMap.Data, appConfig.Spec.Data) {
		configMap.Data = appConfig.Spec.Data
		err = r.Update(ctx, &configMap)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AppConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&nsalexamycomv1alpha1.AppConfig{}).
		Named("appconfig").
		Complete(r)
}
