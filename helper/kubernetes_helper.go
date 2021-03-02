package helper

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/k8s"
	batch "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetJobE get the job and status.
func GetJobE(t *testing.T, options *k8s.KubectlOptions, name string) (*batch.Job, error) {
	output, err := k8s.RunKubectlAndGetOutputE(t, options, "get", "job", name, "-o", "json")
	if err != nil {
		return nil, err
	}
	var job batch.Job
	err = json.Unmarshal([]byte(output), &job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// GetPodsByNameE returns all pods that match a particular name
func GetPodsByNameE(t *testing.T, options *k8s.KubectlOptions, name string) ([]corev1.Pod, error) {
	pods, err := k8s.ListPodsE(t, options, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	result := make([]corev1.Pod, 0)
	for _, pod := range pods {
		if strings.Contains(pod.Name, name) {
			result = append(result, pod)
		}
	}
	return result, nil
}
