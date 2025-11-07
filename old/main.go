package main

import (
	"context"
	"flag"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta1 "k8s.io/api/autoscaling/v2beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Error creating Kubernetes client: %v", err)
	}

	ctx := context.Background()
	namespace := "default"

	klog.Info("=== Creating Deployment ===")
	if err := createDeployment(ctx, clientset, namespace); err != nil {
		klog.Fatalf("Failed to create deployment: %v", err)
	}

	klog.Info("=== Creating HPA using deprecated autoscaling/v2beta1 API ===")
	if err := createHPA(ctx, clientset, namespace); err != nil {
		klog.Fatalf("Failed to create HPA: %v", err)
	}

	klog.Info("Demo completed successfully!")
}

func createDeployment(ctx context.Context, clientset *kubernetes.Clientset, namespace string) error {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-deployment",
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:1.20",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("100m"),
									corev1.ResourceMemory: resource.MustParse("128Mi"),
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := clientset.AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	klog.Info("Created Deployment: nginx-deployment")
	return nil
}

func createHPA(ctx context.Context, clientset *kubernetes.Clientset, namespace string) error {
	hpa := &autoscalingv2beta1.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-hpa",
			Namespace: namespace,
		},
		Spec: autoscalingv2beta1.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2beta1.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       "nginx-deployment",
			},
			MinReplicas: int32Ptr(2),
			MaxReplicas: 10,
			Metrics: []autoscalingv2beta1.MetricSpec{
				{
					Type: autoscalingv2beta1.ResourceMetricSourceType,
					Resource: &autoscalingv2beta1.ResourceMetricSource{
						Name:                     corev1.ResourceCPU,
						TargetAverageUtilization: int32Ptr(70),
					},
				},
			},
		},
	}

	_, err := clientset.AutoscalingV2beta1().HorizontalPodAutoscalers(namespace).Create(ctx, hpa, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	klog.Info("Created HPA using autoscaling/v2beta1 (deprecated API): nginx-hpa")
	return nil
}

func int32Ptr(i int32) *int32 {
	return &i
}
