package ibmserviceoperator

import (
	"context"
	"reflect"
	"os"

	ibmv1 "projects/simple-operator/pkg/apis/ibm/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_ibmserviceoperator")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new IBMserviceOperator Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileIBMserviceOperator{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("ibmserviceoperator-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource IBMserviceOperator
	err = c.Watch(&source.Kind{Type: &ibmv1.IBMserviceOperator{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner IBMserviceOperator
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &ibmv1.IBMserviceOperator{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileIBMserviceOperator implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileIBMserviceOperator{}

// ReconcileIBMserviceOperator reconciles a IBMserviceOperator object
type ReconcileIBMserviceOperator struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a IBMserviceOperator object and makes changes based on the state read
// and what is in the IBMserviceOperator.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileIBMserviceOperator) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling IBMserviceOperator")

	// Fetch the IBMserviceOperator instance
	instance := &ibmv1.IBMserviceOperator{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Info("IBMserviceOperator resource not found. Ignoring since object must be deleted")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Check if the deployment already exists, if not create a new one
	found := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := r.deploymentForIBMserviceOperator(instance)
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return reconcile.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Deployment")
		return reconcile.Result{}, err
	}

	// Ensure the deployment size is the same as the spec
	size := instance.Spec.Size
        if *found.Spec.Replicas != size {
                found.Spec.Replicas = &size
                err = r.client.Update(context.TODO(), found)
                if err != nil {
                        reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
                        return reconcile.Result{}, err
                }
                // Spec updated - return and requeue
                return reconcile.Result{Requeue: true}, nil
        }

	// Update the IBMserviceOperator status with the pod names
	// List the pods for this IBMserviceOperator's deployment
	reqLogger.Info("Updating node status", "looking for:", labelsForIBMserviceOperator(instance.Name))
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(instance.Namespace),
		client.MatchingLabels(labelsForIBMserviceOperator(instance.Name)),
	}
	if err = r.client.List(context.TODO(), podList, listOpts...); err != nil {
		reqLogger.Error(err, "Failed to list pods", "IBMserviceOperator.Namespace", instance.Namespace, "IBMserviceOperator.Name", instance.Name)
		return reconcile.Result{}, err
	}
	podNames := getPodNames(podList.Items)
	reqLogger.Info("Got pod names", "Pods:", podNames)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, instance.Status.Nodes) {
		reqLogger.Info("Instance pods", "before update Pods:", instance.Status.Nodes)
                instance.Status.Nodes = podNames
                err := r.client.Status().Update(context.TODO(), instance)
                if err != nil {
                        reqLogger.Error(err, "Failed to update IBMserviceOperator status")
                        return reconcile.Result{}, err
                }
		reqLogger.Info("Done updating operator status")
		return reconcile.Result{}, nil
        }

	reqLogger.Info("Nothing to update, requeue-ing")
	return reconcile.Result{}, nil
}

// deploymentForIBMserviceOperator returns a instance Deployment object
func (r *ReconcileIBMserviceOperator) deploymentForIBMserviceOperator(m *ibmv1.IBMserviceOperator) *appsv1.Deployment {
	ls := labelsForIBMserviceOperator(m.Name)
	replicas := m.Spec.Size

	// Propagate the horizon env vars to the service.
        hznConfigMap := os.Getenv("HZN_ENV_VARS")

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:   "openhorizon/ibm.helloworld_amd64:1.0.0",
						Name:    "ibm-operator",
						EnvFrom:[]corev1.EnvFromSource{{
                                                        ConfigMapRef: &corev1.ConfigMapEnvSource{
                                                                LocalObjectReference: corev1.LocalObjectReference{
                                                                        Name: hznConfigMap,
                                                                },
                                                        },
                                                }},
						Command: []string{"/service.sh"},
						ImagePullPolicy: corev1.PullAlways,
						Ports: []corev1.ContainerPort{{
							ContainerPort: 8347,
							Name:          "ibm-operator",
						}},
					}},
				},
			},
		},
	}
	// Set IBMserviceOperator instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

// labelsForIBMserviceOperator returns the labels for selecting the resources
// belonging to the given ibmserviceoperator CR name.
func labelsForIBMserviceOperator(name string) map[string]string {
	return map[string]string{"app": "ibm-operator", "ibm-operator_cr": name}
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

