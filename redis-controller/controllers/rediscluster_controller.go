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
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	clusterv1 "lostar.com/rediscluster/api/v1"
)

// RedisclusterReconciler reconciles a Rediscluster object
type RedisclusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cluster.lostar.com,resources=redisclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cluster.lostar.com,resources=redisclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cluster.lostar.com,resources=redisclusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Rediscluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *RedisclusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	reportlog := log.Log.WithValues("namespace", req.NamespacedName.Namespace, "Name", req.NamespacedName.Name)
	// TODO(user): your logic here
	var redisCluster clusterv1.Rediscluster

	if err := r.Get(ctx, req.NamespacedName, &redisCluster); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	switch redisCluster.Status.Phase {
	case clusterv1.PENDING:
		// Add finalizer flag if not exist
		if !controllerutil.ContainsFinalizer(&redisCluster, clusterv1.RedisClusterFinalizer) {
			controllerutil.AddFinalizer(&redisCluster, clusterv1.RedisClusterFinalizer)
			if err := r.Update(ctx, &redisCluster); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
		redisCluster.Status.Phase = clusterv1.CREATING
		if err := r.Status().Update(ctx, &redisCluster); err != nil {
			return ctrl.Result{}, err
		}
	case clusterv1.CREATING:
		for i := 0; int32(i) < redisCluster.Spec.Replicas; i++ {
			if err := r.CreatePod(ctx, &redisCluster); err != nil {
				return ctrl.Result{}, err
			}
		}

		redisCluster.Status.Phase = clusterv1.RUNNING
		if err := r.Status().Update(ctx, &redisCluster); err != nil {
			return ctrl.Result{}, err
		}
	case clusterv1.RUNNING:
		// List pod with labels
		count, err := r.getRedisclusterPodsCount(ctx, &redisCluster)
		if err != nil {
			reportlog.Error(err, "Failed to list pods")
			return ctrl.Result{}, err
		}
		reportlog.Info("Redis cluster relate pod", "count", count)

	case clusterv1.FAILED:
		reportlog.Info("Reconciling Redis Cluster")

	case clusterv1.DELETING:
		reportlog.Info("Reconciling Redis Cluster Deleting")
	default:
		redisCluster.Status.Phase = clusterv1.PENDING
		if err := r.Status().Update(ctx, &redisCluster); err != nil {
			return ctrl.Result{}, err
		}
	}

	//delete resource when kubectl delete
	if !redisCluster.ObjectMeta.DeletionTimestamp.IsZero() {

		if err := r.CleanActualResource(reportlog, &redisCluster); err != nil {
			return ctrl.Result{}, err
		}
		redisCluster.Status.Phase = clusterv1.DELETED
		if err := r.Status().Update(ctx, &redisCluster); err != nil {
			return ctrl.Result{}, err
		}
		// Remove finalizername
		controllerutil.RemoveFinalizer(&redisCluster, clusterv1.RedisClusterFinalizer)
		if err := r.Update(ctx, &redisCluster); err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RedisclusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&clusterv1.Rediscluster{}).
		Complete(r)
}

func (r *RedisclusterReconciler) calculateMD5(input string) string {
	hasher := md5.New()

	hasher.Write([]byte(input))

	hashBytes := hasher.Sum(nil)
	md5String := hex.EncodeToString(hashBytes)

	return md5String
}

func (r *RedisclusterReconciler) randomString(length int) string {
	const charset = "qwertyuiopasdfghjklzxcvbnm1234567890"

	rand.Seed(time.Now().UnixNano())

	randomBytes := make([]byte, length)

	for i := 0; i < length; i++ {
		randomBytes[i] = charset[rand.Intn(len(charset))]
	}
	return string(randomBytes)
}

func (r *RedisclusterReconciler) CreatePod(ctx context.Context, redisCluster *clusterv1.Rediscluster) error {
	redisClusterPodFlag := r.calculateMD5(redisCluster.Spec.Name)[:9]
	randomFlag := r.randomString(5)
	// Create redis cluster
	redisPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      redisCluster.Name + "-" + redisClusterPodFlag + "-" + randomFlag,
			Namespace: redisCluster.Namespace,
			Labels:    redisCluster.Labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:      redisCluster.Spec.Name,
					Image:     redisCluster.Spec.Image,
					Resources: redisCluster.Spec.Resources,
					Env:       redisCluster.Spec.Env,
					EnvFrom:   redisCluster.Spec.EnvFrom,
				},
			},
		},
	}

	if err := r.Create(ctx, redisPod); err != nil {
		return err
	}
	return nil
}

func (r *RedisclusterReconciler) getRedisclusterPodsCount(ctx context.Context, redisCluster *clusterv1.Rediscluster) (int32, error) {
	// List pod with labels
	podList := &corev1.PodList{}
	if err := r.List(ctx, podList, client.InNamespace(redisCluster.Namespace), client.MatchingLabels(redisCluster.Labels)); err != nil {
		return 0, err
	}
	return int32(len(podList.Items)), nil
}

func (r *RedisclusterReconciler) CleanActualResource(reportlog logr.Logger, redisCluster *clusterv1.Rediscluster) error {
	ctx := context.Background()
	redisCluster.Status.Phase = clusterv1.DELETING
	if err := r.Status().Update(ctx, redisCluster); err != nil {
		return err
	}
	// Delete the associated Pod with the label
	if err := r.DeleteAllOf(ctx, &corev1.Pod{}, client.InNamespace(redisCluster.Namespace), client.MatchingLabels(redisCluster.Labels)); err != nil {
		reportlog.Error(err, "Failed to delete rediscluster resource")
		return err
	}
	reportlog.Info("Deleted rediscluster resource")
	return nil
}
