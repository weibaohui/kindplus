package installer

import (
	"fmt"

	"github.com/weibaohui/kom/kom"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/klog/v2"
)

type Installer struct {
	// 集群配置
	Config *Config
	// 集群运行时
	Runtime *Runtime
}

type Config struct {
	BaseDomain string // Kind 外部访问Kind集群使用的域名, 如：dev.k8m.site,最终集群内svc访问域名就是cluster-app-svc.cluster-name.dev.k8m.site
	Port       int    // Kind APIServer 端口
	Name       string // Kind 集群名称
	Namespace  string // kind 集群安装在哪个命名空间
}

func (i *Installer) Deploy() {
	// Create Deployment
	name := i.Config.Name
	ns := i.Config.Namespace
	deployment := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			Labels: map[string]string{
				"app": name,
			},
		},
		Spec: v1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            name,
							Image:           i.Runtime.BaseImage,
							ImagePullPolicy: corev1.PullAlways,
							SecurityContext: &corev1.SecurityContext{
								Privileged: new(bool),
							},
							Env: []corev1.EnvVar{
								{
									Name:  "KIND_CLUSTER_NAME",
									Value: name,
								}, {
									Name:  "KIND_CLUSTER_PORT",
									Value: fmt.Sprintf("%v", i.Config.Port),
								}, {
									Name:  "DOMAIN",
									Value: fmt.Sprintf("%s.%s", name, i.Config.BaseDomain),
								},
								{
									Name: "KIND_CLUSTER_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	*deployment.Spec.Template.Spec.Containers[0].SecurityContext.Privileged = true

	err := kom.DefaultCluster().Resource(deployment).Create(deployment).Error
	if err != nil {
		klog.Errorf("Failed to create Deployment: %v", err)
	}
	klog.Infof("Deployment created")

	// Create Service
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       "web-80",
					Port:       32480,
					TargetPort: intstr.FromInt32(32480),
				},
				{
					Name:       "web-443",
					Port:       32443,
					TargetPort: intstr.FromInt32(32443),
				},
				{
					Name:       "apiserver",
					Port:       int32(i.Config.Port),
					TargetPort: intstr.FromInt32(int32(i.Config.Port)),
				},
			},
			Selector: map[string]string{
				"app": name,
			},
		},
	}
	err = kom.DefaultCluster().Resource(service).Create(service).Error
	if err != nil {
		klog.Errorf("Failed to create Service: %v", err)
	}
	klog.Infof("Service created")

	// Create Ingress
	prefix := networkingv1.PathTypePrefix
	host := fmt.Sprintf("*.%s", i.Config.BaseDomain)
	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: new(string),
			Rules: []networkingv1.IngressRule{
				{
					Host: host,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &prefix,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: name,
											Port: networkingv1.ServiceBackendPort{
												Number: 32480,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	*ingress.Spec.IngressClassName = "nginx"
	err = kom.DefaultCluster().Resource(ingress).Create(ingress).Error
	if err != nil {
		klog.Errorf("Failed to create Ingress: %v", err)
	}
	klog.Infof("Ingress created")

}
func (i *Installer) Clean() {
	// Create Deployment
	name := i.Config.Name
	n := i.Config.Namespace

	err := kom.DefaultCluster().Resource(&v1.Deployment{}).Namespace(n).Name(name).Delete().Error
	if err != nil {
		klog.Errorf("Failed to delete Deployment: %v", err)
	}
	klog.Infof("Deployment deleted")

	err = kom.DefaultCluster().Resource(&corev1.Service{}).Namespace(n).Name(name).Delete().Error
	if err != nil {
		klog.Errorf("Failed to delete Service: %v", err)
	}
	klog.Infof("Service deleted")

	err = kom.DefaultCluster().Resource(&networkingv1.Ingress{}).Namespace(n).Name(name).Delete().Error
	if err != nil {
		klog.Errorf("Failed to delete Ingress: %v", err)
	}
	klog.Infof("Ingress deleted")

}
