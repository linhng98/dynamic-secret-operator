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
	"math/rand"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	secretv1alpha1 "github.com/linhng98/dynamic-secret-operator/api/v1alpha1"
	"github.com/linhng98/dynamic-secret-operator/vars"
	"k8s.io/apimachinery/pkg/api/errors"
)

// PlaintextReconciler reconciles a Plaintext object
type PlaintextReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=secret.linhng98.com,resources=plaintexts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=secret.linhng98.com,resources=plaintexts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=secret.linhng98.com,resources=plaintexts/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Plaintext object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *PlaintextReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the Plaintext crd
	plaintextCrd := &secretv1alpha1.Plaintext{}
	err := r.Get(ctx, req.NamespacedName, plaintextCrd)

	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get crd plaintext")
		return ctrl.Result{}, err
	}

	defer func() {
		err = r.Client.Status().Update(context.Background(), plaintextCrd)
		if err != nil {
			log.Error(err, "Status update failed")
		}
	}()

	// Check if the secret already exists, if not create a new one
	sec := &corev1.Secret{}
	err = r.Get(ctx, req.NamespacedName, sec)
	if err != nil && errors.IsNotFound(err) {
		// Define a new Secret
		opaqueSec := r.generateOpaqueSecret(plaintextCrd)
		log.Info("Creating a new Secret", "Secret.Namespace", opaqueSec.Namespace, "Secret.Name", opaqueSec.Name)
		err = r.Create(ctx, opaqueSec)
		if err != nil {
			plaintextCrd.Status.Phase = secretv1alpha1.PhaseError
			log.Error(err, "Failed to create new Secret", "Secret.Namespace", opaqueSec.Namespace, "Secret.Name", opaqueSec.Name)
			return ctrl.Result{}, err
		}
		// Secret created successfully - set status then return and requeue

		plaintextCrd.Status.Phase = secretv1alpha1.PhaseReady
		plaintextCrd.Status.LastSyncedDate = metav1.Now()
		err = r.Client.Status().Update(context.Background(), plaintextCrd)
		if err != nil {
			log.Error(err, "Status update failed")
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Secret")
		plaintextCrd.Status.Phase = secretv1alpha1.PhaseReady
		_ = r.Client.Status().Update(context.Background(), plaintextCrd)

		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *PlaintextReconciler) generateOpaqueSecret(pt *secretv1alpha1.Plaintext) *corev1.Secret {
	mapData := map[string]string{}
	for _, v := range pt.Spec.Secrets {
		mapData[v.Key] = generateTextSecret(v.Whitelist, v.Len, v.Prefix, v.Postfix)
	}

	sec := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pt.Name,
			Namespace: pt.Namespace,
			Labels:    pt.Labels,
		},
		StringData: mapData,
	}

	ctrl.SetControllerReference(pt, sec, r.Scheme)
	return sec
}

func generateTextSecret(whitelist string, lenght int, prefix string, postfix string) string {
	if whitelist == "" {
		whitelist = vars.DefaultWhitelist
	}

	if lenght <= 0 {
		lenght = vars.DefaultLen
	}

	chars := strings.Split(whitelist, "")
	s := rand.NewSource(time.Now().UnixNano())
	n := len(chars)
	textSec := ""

	for i := 0; i < lenght; i++ {
		num := rand.New(s).Int()
		textSec += chars[num%n]
	}

	return prefix + textSec + postfix
}

// SetupWithManager sets up the controller with the Manager.
func (r *PlaintextReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&secretv1alpha1.Plaintext{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}
