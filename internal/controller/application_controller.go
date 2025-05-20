/*
Copyright 2025.

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
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	almcv1 "github.com/plimbore/applicationmgr/api/v1"
)

func int32Ptr(i int32) *int32 { return &i }

// ApplicationReconciler reconciles a Application object
type ApplicationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=almc.applicationmgr.io,resources=applications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=almc.applicationmgr.io,resources=applications/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=almc.applicationmgr.io,resources=applications/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete

func (r *ApplicationReconciler) getDeployment(ctx context.Context, namespace string, name string) (*appsv1.Deployment, error) {
	deployment := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, deployment)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}
	return deployment, nil
}

func (r *ApplicationReconciler) getService(ctx context.Context, namespace string, name string) (*corev1.Service, error) {
	service := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, service)
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}
	return service, nil
}

func (r *ApplicationReconciler) getIngress(ctx context.Context, namespace string, name string) (*networkingv1.Ingress, error) {
	ingress := &networkingv1.Ingress{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, ingress)
	if err != nil {
		return nil, fmt.Errorf("failed to get ingress: %w", err)
	}
	return ingress, nil
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.4/pkg/reconcile
func (r *ApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Load application by name
	var application almcv1.Application
	if err := r.Get(ctx, req.NamespacedName, &application); err != nil {
		log.Error(err, "unable to fetch Application")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	} else {
		log.Info("Application found")
	}

	// Create/Update deployment
	deploymentName := req.Name + "-deployment"
	deploymentNamespace := req.Namespace
	log.Info("Finding deployment with name: " + deploymentName + " in namespace: " + deploymentNamespace)
	deployment, err := r.getDeployment(ctx, deploymentNamespace, deploymentName)
	if err != nil {
		// Handle the error appropriately, e.g., log it or requeue the request
		log.Info("Deployment " + deploymentName + " not found in namespace: " + deploymentNamespace + " or unable to fetch")
	}

	if deployment != nil {
		log.Info("Updating deployment: " + deploymentName + " in namespace: " + deploymentName)

		deployment.Spec.Template.Spec.Containers[0].Name = req.Name
		deployment.Spec.Template.Spec.Containers[0].Image = application.Spec.Image.Repository + ":" + application.Spec.Image.Tag
		deployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = int32(application.Spec.Service.Port)

	} else {
		// Deployment does not exist, create deployment
		deployment = &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      deploymentName,
				Namespace: deploymentNamespace,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: int32Ptr(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": req.Name,
					},
				},
				Template: apiv1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": req.Name,
						},
					},
					Spec: apiv1.PodSpec{
						Containers: []apiv1.Container{
							{
								Name:  req.Name,
								Image: application.Spec.Image.Repository + ":" + application.Spec.Image.Tag,
								Ports: []apiv1.ContainerPort{
									{
										Name:          "http",
										Protocol:      apiv1.ProtocolTCP,
										ContainerPort: int32(application.Spec.Service.Port),
									},
								},
							},
						},
					},
				},
			},
		}

		// Set the custom resource instance as the owner and controller
		if err := controllerutil.SetControllerReference(&application, deployment, r.Scheme); err != nil {
			log.Error(err, "Setting owner for service failed")
		}

		// Create Deployment
		log.Info("Creating deployment" + deploymentName + " in namespace: " + deploymentName)
		err = r.Create(context.TODO(), deployment) //, metav1.CreateOptions{})
		if err != nil {
			log.Error(err, "Creating deployment failed")
		} else {
			log.Info("Created deployment\n")
		}
	}

	log.Info("Deployment create/update finished")

	// Create/Update service
	serviceName := req.Name + "-service"
	serviceNamespace := req.Namespace
	log.Info("Finding service with name: " + serviceName + " in namespace: " + serviceNamespace)
	service, err := r.getService(ctx, serviceNamespace, serviceName)
	if err != nil {
		// Handle the error appropriately, e.g., log it or requeue the request
		log.Info("Service " + serviceName + " not found in namespace: " + serviceNamespace + " or unable to fetch")
	}

	if service != nil {
		log.Info("Updating service: " + serviceName + " in namespace: " + serviceName)

		service.Spec.Ports[0].Port = int32(application.Spec.Service.Port)

	} else {
		// Service does not exist, create service
		service = &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      serviceName,
				Namespace: serviceNamespace,
			},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{
					"app": req.Name,
				},
				Ports: []corev1.ServicePort{
					{
						Port:     int32(application.Spec.Service.Port),
						Protocol: apiv1.ProtocolTCP,
					},
				},
			},
		}

		// Set the custom resource instance as the owner and controller
		if err := controllerutil.SetControllerReference(&application, service, r.Scheme); err != nil {
			log.Error(err, "Setting owner for service failed")
		}

		// Create Service
		log.Info("Creating service" + serviceName + " in namespace: " + serviceName)
		err = r.Create(context.TODO(), service) //, metav1.CreateOptions{})
		if err != nil {
			log.Error(err, "Creating service failed")
		} else {
			log.Info("Created service\n")
		}
	}

	log.Info("Service create/update finished")

	// Create/Update ingress
	ingressName := req.Name + "-ingress"
	ingressNamespace := req.Namespace
	log.Info("Finding ingress with name: " + ingressName + " in namespace: " + ingressNamespace)
	ingress, err := r.getIngress(ctx, ingressNamespace, ingressName)
	if err != nil {
		// Handle the error appropriately, e.g., log it or requeue the request
		log.Info("Ingress " + ingressName + " not found in namespace: " + ingressNamespace + " or unable to fetch")
	}

	// Create HTTP ingress service backend for later use
	httpIngressServiceBackend := &networkingv1.IngressServiceBackend{
		Name: serviceName,
		Port: networkingv1.ServiceBackendPort{
			Number: int32(application.Spec.Service.Port),
		},
	}

	var ingressPathType = networkingv1.PathTypeImplementationSpecific

	// Create ingress rules object for later use
	newIngressRulesLen := len(application.Spec.Ingress.Hosts)
	ingressRules := make([]networkingv1.IngressRule, newIngressRulesLen)
	for i := 0; i < newIngressRulesLen; i++ {
		// Create ingress rules, http paths object for later use
		pathsLen := len(application.Spec.Ingress.Hosts[i].Paths)
		httpIngressPaths := make([]networkingv1.HTTPIngressPath, pathsLen)
		for j := 0; j < pathsLen; j++ {
			// Set http ingress path
			httpIngressPaths[j] = networkingv1.HTTPIngressPath{
				Path:     application.Spec.Ingress.Hosts[i].Paths[j].Path,
				PathType: &ingressPathType,
				Backend: networkingv1.IngressBackend{
					Service: httpIngressServiceBackend,
				},
			}
		}

		// Set ingress rule
		ingressRules[i] = networkingv1.IngressRule{
			Host: application.Spec.Ingress.Hosts[i].Host,
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{
					Paths: httpIngressPaths,
				},
			},
		}
	}

	if ingress != nil {
		log.Info("Updating ingress: " + ingressName + " in namespace: " + ingressName)

		ingress.Spec.Rules = ingressRules

	} else {
		// Ingress does not exist, create ingress
		ingress = &networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      ingressName,
				Namespace: ingressNamespace,
			},
			Spec: networkingv1.IngressSpec{
				Rules: ingressRules,
			},
		}

		// Set the custom resource instance as the owner and controller
		if err := controllerutil.SetControllerReference(&application, ingress, r.Scheme); err != nil {
			log.Error(err, "Setting owner for ingress failed")
		}

		// Create Ingress
		log.Info("Creating ingress" + ingressName + " in namespace: " + ingressName)
		err = r.Create(context.TODO(), ingress) //, metav1.CreateOptions{})
		if err != nil {
			log.Error(err, "Creating ingress failed")
		} else {
			log.Info("Created ingress\n")
		}
	}

	log.Info("Ingress create/update finished")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&almcv1.Application{}).
		Named("application").
		Complete(r)
}
