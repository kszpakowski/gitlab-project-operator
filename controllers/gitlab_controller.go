/*
Copyright 2022.

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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	gitlabv1alpha1 "github.com/kszpakowski/gitlab-project-operator/api/v1alpha1"
	gitlabapi "github.com/xanzy/go-gitlab"
)

// GitlabReconciler reconciles a Gitlab object
type GitlabReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=gitlab.kszpakowski.com,resources=gitlabs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=gitlab.kszpakowski.com,resources=gitlabs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=gitlab.kszpakowski.com,resources=gitlabs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Gitlab object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *GitlabReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// your logic here
	gitlab := &gitlabv1alpha1.Gitlab{}
	err := r.Get(ctx, req.NamespacedName, gitlab)
	if err != nil {
		log.Error(err, "Failed to get Gitlab")
		return ctrl.Result{}, err
	}

	log.Info("Found gitlab instance:", "Instance name", gitlab.Spec.Name)

	//git, err := gitlab.NewClient("yourtokengoeshere", gitlab.WithBaseURL("https://git.mydomain.com/api/v4"))
	apiClient, err := gitlabapi.NewClient(gitlab.Spec.ApiToken)
	if err != nil {
		log.Error(err, "Unable to connect to gitlab")
	}

	version, _, err := apiClient.Version.GetVersion()
	if err != nil {
		log.Error(err, "Unable to get gitlab version")
	}

	log.Info(version.String())

	gitlab.Status.Version = version.String()

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GitlabReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gitlabv1alpha1.Gitlab{}).
		Complete(r)
}
