package autokill

import (
	"context"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"time"
	agillv1alpha1 "github.com/agill17/namespace-manager/pkg/apis/agill/v1alpha1"

)

func (r *ReconcileAutoKill) nsObject(cr *agillv1alpha1.AutoKill) (*v1.Namespace, error){
	ns := &v1.Namespace{}
	if err := r.client.Get(context.TODO(), types.NamespacedName{Namespace:"", Name:cr.Namespace}, ns); err != nil {
		logrus.Errorf("Error while getting namepace object: %v", err)
		return nil, err
	}
	return ns, nil
}

// calcTimeLeft will return currentAge of ns by subtracting currentTime with nsCreationTime in Hours
func calcTimeLeft(ns *v1.Namespace) int {
	return int(time.Now().Sub(ns.CreationTimestamp.Time).Hours())
}


//deleteNs takes a Namespace object and performs a Delete
func (r *ReconcileAutoKill)deleteNs(ns *v1.Namespace) error {
	if ns != nil {
		logrus.Warnf("Namespace: %v | Msg: Deleting Namespace", ns.Name )
		return r.client.Delete(context.TODO(), ns)
	}
	return nil
}

