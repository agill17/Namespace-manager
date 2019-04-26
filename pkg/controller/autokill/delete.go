package autokill

import (
	agillv1alpha1 "github.com/agill17/namespace-manager/pkg/apis/agill/v1alpha1"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
)

import "k8s.io/api/core/v1"

func (r *ReconcileAutoKill) decideAndDelete(cr *agillv1alpha1.AutoKill, ns *v1.Namespace) error {

	currentAge := calcCurrentAge(ns)
	hasLivedOverPolicyAge := currentAge >= cr.Spec.DeleteNamespaceAfter


	if !cr.Spec.Disable && hasLivedOverPolicyAge && ns.Status.Phase != v1.NamespaceTerminating {

		logrus.Warnf("Namespace: %v | Namespace Age: %v | CR Disabled: %v | Policy Age: %v", cr.Namespace, currentAge, cr.Spec.Disable, cr.Spec.DeleteNamespaceAfter)

		// delete helm releases ?
		if cr.Spec.DeleteAssociatedHelmReleases {

			// get helm releases associated to this namespace
			helmReleases , err := getReleasesForNs(cr.Namespace, cr.Spec.TillerNamespace)
			if err != nil {
				return err
			}

			// delete helm releases first
			if len(helmReleases) > 1 {
				if err := deleteReleases(cr.Spec.TillerNamespace, helmReleases); err != nil {
					return  err
				}
			} else {
				logrus.Infof("Namespace: %v | Msg: Could not find any helm releases associated to this namespace.", cr.Namespace)
			}

		}


		// delete namespace
		logrus.Warnf("Deleting Namespace: %v", cr.Namespace)
		if err := r.deleteNs(ns); err != nil && !errors.IsForbidden(err) {
			return err
		}

	}
	return nil
}

