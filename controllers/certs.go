package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	certmanagerv1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1"
	api "github.com/ludovicus3/foreman-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	SelfSignedName           string = "foreman-selfsign"
	ForemanCAName            string = "foreman-ca"
	ForemanCertificateName   string = "foreman-tls"
	PulpcoreCertificateName  string = "pulpcore-tls"
	CandlepinCertificateName string = "candlepin-tls"
)

func (r *ForemanReconciler) reconcileCerts(ctx context.Context, log logr.Logger, foreman *api.Foreman) error {
	// reconcile self-signed Issuer
	if err := reconcileSelfSignedIssuer(ctx, log, foreman); err != nil {
		return err
	}
	// reconcile foreman-ca Certificate
	if err := reconcileCACertificate(ctx, log, foreman); err != nil {
		return err
	}
	// reconcile foreman-ca Secret
	caSecret := &corev1.Secret{}
	if err := r.Get(ctx, types.NamespacedName{Name: ForemanCAName, Namespace: foreman.Namespace}, caSecert); err != nil {
		return err
	}
	// reconcile foreman-ca Issuer
	if err := reconcileCAIssuer(ctx, log, foreman); err != nil {
		return err
	}
	// reconcile foreman-client-tls Certificate
	if err := reconcileClient(ctx, log, foreman); err != nil {
		return err
	}
	// reconcile foreman-client-tls Secret
	clientSecret := &corev1.Secret{}
	if err := r.Get(ctx, types.NamespacedName{Name: ForemanCertificateName, Namespace: foreman.Namespace}, clientSecret); err != nil {
		return err
	}
	// reconcile for each service
	for _, service := range foreman.Certs.Services {
		//   reconcile service-tls Certificate
		if err := reconcileCertificate(ctx, log, foreman); err != nil {
			return err
		}
		//   reconcile service-tls Secret
		if err := reconcileTLSSecret(ctx, service.SecretName, foreman); err != nil {
			return err
		}
	}
	// end services
}

func (r *ForemanReconciler) reconcileSelfSignedIssuer(ctx context.Context, log logr.Logger, foreman *api.Foreman) error {
	issuer := r.newSelfServeIssuer(foreman)
	existingIssuer := &certmanagerv1.Issuer{}
	err := r.Get(ctx, types.NamespacedName{Name: SelfSignedName, Namespace: foreman.Namespace}, existingIssuer)
	if err != nil {
		if errors.IsNotFound(err) {
			if err = ctrl.SetControllerReference(foreman, issuer, r.Scheme); err != nil {
				return &certmanagerv1.Issuer{}, fmt.Errorf("Unable to set controller reference on issuer %s, %v", SelfSignedName, err)
			}
			return issuer, r.Create(ctx, issuer)
		}
		return &certmanagerv1.Issuer{}, fmt.Errorf("Unabled to query existing issuer %s, %v", SelfSignedName)
	}
	return existingIssuer, nil
}

func (r *ForemanReconciler) newSelfSignedIssuer(selfSigned bool, foreman *api.Foreman) *certmanagerv1.Issuer {
	issuer := &certmanagerv1.Issuer{
		TypeMeta: metav1.TypeMeta{
			APIVersion: certmanagerv1.SchemeGroupVersion.String(),
			Kind:       certmanagerv1.IssuerKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      SelfSignedIssuerName,
			Namespace: foreman.Namespace,
		},
		Spec: certmanagerv1.IssuerSpec{
			IssuerConfig: newIssuerConfig(nil),
		},
	}

	return issuer
}

func newIssuerConfig(secret *corev1.Secret) {
	if secret != nil {
		return certmanagerv1.IssuerConfig{
			CA: &certmanagerv1.CAIssuer{
				SecertName: secret.Name,
			}
		}
	} else {
		return certmanagerv1.IssuerConfig{
			SelfSigned: &certmanagerv1.SelfSignedIssuer{},
		}
	}
}

func (r *ForemanReconciler) reconcileCACert(ctx context.Context, log logr.Logger, foreman *api.Foreman) {

}

func (r *ForemanReconciler) reconcileTLSSecret(ctx context.Context, name string, foreman *api.Foreman) error {
	secret := &corev1.Secret{}
	return r.Get(ctx, types.NamespacedName{Name: name, Namespace: foreman.Namespace}, secret)
}
