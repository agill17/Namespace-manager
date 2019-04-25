package autokill

import (
	"context"
	"fmt"
	agillv1alpha1 "github.com/agill17/namespace-manager/pkg/apis/agill/v1alpha1"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"os/exec"
	"strings"
	"time"
)

// gets the k8s namespace object
func (r *ReconcileAutoKill) nsObject(cr *agillv1alpha1.AutoKill) (*v1.Namespace, error){
	ns := &v1.Namespace{}
	if err := r.client.Get(context.TODO(), types.NamespacedName{Namespace:"", Name:cr.Namespace}, ns); err != nil {
		logrus.Errorf("Error while getting namespace object: %v", err)
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

//deleteReleases deletes list of helm releases
func deleteReleases(tillerNamespace string, releaseNames []string) {
	logrus.Warnf("Helm Releases to delete: %v", releaseNames)
	releases := strings.Join(releaseNames, " ")
	delCmd := fmt.Sprintf("helm del --purge --tiller-namespace=%v %v", tillerNamespace, releases)
	logrus.Infof("Delete Command: %v", delCmd)
	out, _ := runACommand(delCmd)
	logrus.Infof("Helm delete output: %v", out)

}

//getReleasesForNs gets list of releases that belongs to a namespace
func getReleasesForNs(ns ,tillerNamespace string) ([]string,error) {
	helmLsCmd := fmt.Sprintf("helm ls --tiller-namespace=%v --namespace=%v --short", tillerNamespace, ns)
	helmLsOut, err := runACommand(helmLsCmd)
	if err != nil {
		logrus.Errorf("Helm ls exec command failed: %v -- %v", err, helmLsOut)
		return nil, err
	}
	return strings.Split(string(helmLsOut), "\n"), nil

}

// yea whatever... wish there was a better package
func runACommand(cmd string) (string, error) {
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head,parts...).Output()
	if err != nil {
		logrus.Infof("%s", err)
	}
	return string(out), err
}