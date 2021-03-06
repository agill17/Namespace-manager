package autokill

import (
	"context"
	"github.com/sirupsen/logrus"

	agillv1alpha1 "github.com/agill17/namespace-manager/pkg/apis/agill/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_autokill")
const defaultTillerNamespace = "kube-system"


// Add creates a new AutoKill Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileAutoKill{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("autokill-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource AutoKill
	err = c.Watch(&source.Kind{Type: &agillv1alpha1.AutoKill{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileAutoKill{}

// ReconcileAutoKill reconciles a AutoKill object
type ReconcileAutoKill struct {
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a AutoKill object and makes changes based on the state read
// and what is in the AutoKill.Spec
func (r *ReconcileAutoKill) Reconcile(request reconcile.Request) (reconcile.Result, error) {

	cr := &agillv1alpha1.AutoKill{}
	err := r.client.Get(context.TODO(), request.NamespacedName, cr)
	if err != nil {
		if errors.IsNotFound(err) || errors.IsForbidden(err)  {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// set default tiller-namespace
	if cr.Spec.TillerNamespace == "" {
		cr.Spec.TillerNamespace = defaultTillerNamespace
		if err := r.client.Update(context.TODO(), cr); err != nil {
			logrus.Errorf("Failed to setup default tiller-namespace for CR: %v in Namespace", cr.Name, cr.Namespace)
		}
	}

	// get ns object
	ns, err := r.nsObject(cr)
	if err != nil {
		return reconcile.Result{},err
	}

	// run the main delete logic
	if err := r.decideAndDelete(cr, ns); err != nil {
		return reconcile.Result{}, err
	}



	return reconcile.Result{Requeue:true}, nil
}

