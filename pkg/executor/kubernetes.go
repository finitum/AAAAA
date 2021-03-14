package executor

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesExecutor struct {
	cs          *kubernetes.Clientset
	namespace   string
	runnerImage string
}

func NewKubernetesExecutor(runnerImage, namespace, cfgpath string) (*KubernetesExecutor, error) {
	config, err := clientcmd.BuildConfigFromFlags("", cfgpath)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't build kubernetes config")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't build kubernetes clientset")
	}

	return &KubernetesExecutor{clientset, namespace, runnerImage}, nil
}

func (k *KubernetesExecutor) PrepareBuild(ctx context.Context) error {
	return nil
}

func (k *KubernetesExecutor) BuildPackage(ctx context.Context, cfg *Config) error {
	job, err := k.buildJobSpec(cfg)
	if err != nil {
		return errors.Wrap(err, "building job spec")
	}

	batch := k.cs.BatchV1().Jobs(k.namespace)

	_, err = batch.Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		return errors.Wrap(err, "creating job")
	}

	return nil
}

func (k *KubernetesExecutor) buildJobSpec(cfg *Config) (*batchv1.Job, error) {
	cfgb, err := json.Marshal(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "marshal config")
	}

	// name is the suffix used for naming jobs/pods/containers
	name := cfg.Package.Name
	// backoffLimit specifies how many times a job should be retried
	backoffLimit := int32(1)
	// deadlineSeconds specifies how long a job (including retries) may run
	deadlineSeconds := int64((time.Hour * 4).Seconds())
	// ttlAfterFinishSeconds specifies how many a seconds a job should be kept around after completion (failure or success)
	ttlAfterFinishSeconds := int32((time.Minute * 15).Seconds())
	return &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "batch/v1",
			Kind:       "Job",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "aaaaa-job-" + name,
			Namespace: k.namespace,
			Labels: map[string]string{
				"name":      "aaaaa-builder",
				"part-of":   "aaaaa",
				"component": "builder",
			},
		},
		Spec: batchv1.JobSpec{
			ActiveDeadlineSeconds:   &deadlineSeconds,
			BackoffLimit:            &backoffLimit,
			TTLSecondsAfterFinished: &ttlAfterFinishSeconds,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "aaaaa-pod-" + name,
					Namespace: k.namespace,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "aaaaa-container-" + name,
							Image: k.runnerImage,
							Env: []corev1.EnvVar{
								{
									Name:  "CONFIG",
									Value: string(cfgb),
								},
							},
						},
					},
					RestartPolicy: "Never",
				},
			},
		},
	}, nil
}
