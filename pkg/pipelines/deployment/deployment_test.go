package deployment

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/google/go-cmp/cmp"
	"github.com/redhat-developer/kam/pkg/pipelines/meta"
)

const (
	testComponent       = "nginx-deployment"
	testImage           = "nginx:1.7.9"
	testComponentPartOf = "nginx-deployment-operator"
)

func TestCreate(t *testing.T) {
	d := Create(testComponentPartOf, "", testComponent, testImage, ContainerPort(80))

	want := &appsv1.Deployment{
		TypeMeta: meta.TypeMeta("Deployment", "apps/v1"),
		ObjectMeta: meta.ObjectMeta(meta.NamespacedName("", testComponent), meta.AddLabels(
			map[string]string{
				KubernetesAppNameLabel: testComponent,
				KubernetesPartOfLabel:  testComponentPartOf,
			},
		)),
		Spec: appsv1.DeploymentSpec{
			Replicas: ptr32(1),
			Selector: LabelSelector(testComponent, testComponentPartOf),
			Template: podTemplate(testComponentPartOf, testComponent, testImage, ContainerPort(80)),
		},
	}

	if diff := cmp.Diff(want, d); diff != "" {
		t.Fatalf("deployment diff:\n%s", diff)
	}
}

func TestDefaultPodTemplate(t *testing.T) {
	testComponent := "test-svc"
	want := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				KubernetesAppNameLabel: testComponent,
				KubernetesPartOfLabel:  testComponentPartOf,
			},
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: "default",
			Containers: []corev1.Container{
				{
					Name:            testComponent,
					Image:           testImage,
					ImagePullPolicy: corev1.PullAlways,
				},
			},
		},
	}

	spec := podTemplate(testComponentPartOf, testComponent, testImage)

	if diff := cmp.Diff(want, spec); diff != "" {
		t.Fatalf("podTemplate diff: %s", diff)
	}
}

func TestPodTemplateEnv(t *testing.T) {
	env := []corev1.EnvVar{{Name: "FOO_BAR_SERVICE_HOST", Value: "1.2.3.4"}}

	spec := podTemplate(testComponentPartOf, testComponent, testImage, Env(env))

	want := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				KubernetesAppNameLabel: testComponent,
				KubernetesPartOfLabel:  testComponentPartOf,
			},
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: "default",
			Containers: []corev1.Container{
				{
					Name:            testComponent,
					Image:           testImage,
					ImagePullPolicy: corev1.PullAlways,
					Env:             env,
				},
			},
		},
	}
	if diff := cmp.Diff(want, spec); diff != "" {
		t.Fatalf("podTemplate diff: %s", diff)
	}
}

func TestPodTemplateCommand(t *testing.T) {
	spec := podTemplate(testComponentPartOf, testComponent, testImage, Command([]string{"/usr/local/bin/test"}))

	want := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				KubernetesAppNameLabel: testComponent,
				KubernetesPartOfLabel:  testComponentPartOf,
			},
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: "default",
			Containers: []corev1.Container{
				{
					Name:            testComponent,
					Image:           testImage,
					ImagePullPolicy: corev1.PullAlways,
					Command:         []string{"/usr/local/bin/test"},
				},
			},
		},
	}
	if diff := cmp.Diff(want, spec); diff != "" {
		t.Fatalf("podTemplate diff: %s", diff)
	}

}

func TestPodTemplateContainerPort(t *testing.T) {
	spec := podTemplate(testComponentPartOf, testComponent, testImage, ContainerPort(80))

	want := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				KubernetesAppNameLabel: testComponent,
				KubernetesPartOfLabel:  testComponentPartOf,
			},
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: "default",
			Containers: []corev1.Container{
				{
					Name:            testComponent,
					Image:           testImage,
					ImagePullPolicy: corev1.PullAlways,
					Ports: []corev1.ContainerPort{
						{ContainerPort: 80},
					},
				},
			},
		},
	}
	if diff := cmp.Diff(want, spec); diff != "" {
		t.Fatalf("podTemplate diff: %s", diff)
	}
}
