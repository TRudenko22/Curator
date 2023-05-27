package controllers

import (
	"context"
	"os"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	curatorv1alpha1 "github.com/TRudenko22/Curator/api/v1alpha1"
)

// FetchDataReconciler reconciles a FetchData object
type FetchDataReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=curator.curator,resources=fetchdata,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=curator.curator,resources=fetchdata/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=curator.curator,resources=fetchdata/finalizers,verbs=update
//+kubebuilder:rbac:groups=batch,resources=cronjobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=cronjobs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=batch,resources=cronjobs/finalizers,verbs=update
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get
//+kubebuilder:rbac:groups="core",resources=persistentvolumeclaims,verbs=get;watch;list;create;update;delete
//+kubebuilder:rbac:groups="core",resources=persistentvolumeclaims,verbs=get;watch;list
//+kubebuilder:rbac:groups="core",resources=persistentvolumeclaims/status,verbs=get;update;patch

func (r *FetchDataReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	FetchData := &curatorv1alpha1.FetchData{}

	err := r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, FetchData)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			l.Info("FetchData resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
	}

	if err := r.createCronJob(ctx, FetchData); err != nil {
		l.Error(err, "failed to create the CronJob resource")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *FetchDataReconciler) createCronJob(ctx context.Context, m *curatorv1alpha1.FetchData) error {
	if _, err := FetchCronJob(ctx, m.Name, m.Namespace, r.Client); err != nil {
		if err := r.Client.Create(ctx, NewCronJob(m, r.Scheme)); err != nil {
			return err
		}
	}

	return nil
}

func FetchCronJob(ctx context.Context, name, namespace string, client client.Client) (*batchv1.CronJob, error) {
	cronJob := &batchv1.CronJob{}
	err := client.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, cronJob)
	return cronJob, err
}

func NewCronJob(m *curatorv1alpha1.FetchData, scheme *runtime.Scheme) *batchv1.CronJob {
	//fmt.Println("Name and Namespace", m.Namespace, m.Name)

	cronjob := &batchv1.CronJob{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CronJob",
			APIVersion: "batch/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Spec.CronjobNamespace,
		},
		Spec: batchv1.CronJobSpec{
			Schedule:          m.Spec.Schedule,
			ConcurrencyPolicy: batchv1.ForbidConcurrent,
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Volumes: []corev1.Volume{
								{
									Name: "koku-metrics-operator-data",
									VolumeSource: corev1.VolumeSource{
										PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
											ClaimName: "koku-metrics-operator-data",
										},
									},
								},
							},
							Containers: []corev1.Container{
								{
									Name:  "crdunzip",
									Image: "docker.io/surbhi0129/crd_unzip:latest",
									Env: []corev1.EnvVar{
										{Name: "BACKUP_SRC", Value: m.Spec.BackupSrc},
										{Name: "UNZIP_DIR", Value: m.Spec.UnzipDir},
										{Name: "DATABASE_NAME", Value: os.Getenv("DATABASE_NAME")},
										{Name: "DATABASE_USER", Value: os.Getenv("DATABASE_USER")},
										{Name: "DATABASE_PASSWORD", Value: os.Getenv("DATABASE_PASSWORD")},
										{Name: "DATABASE_HOST_NAME", Value: os.Getenv("DATABASE_HOST_NAME")},
										{Name: "PORT_NUMBER", Value: os.Getenv("PORT_NUMBER")},
									},
									Command: []string{"python3"},
									Args:    []string{"scripts/unzip_backup.py"},
									VolumeMounts: []corev1.VolumeMount{
										{
											Name:      "koku-metrics-operator-data",
											MountPath: "/tmp/koku-metrics-operator-data",
										},
									},
								},
							},
							RestartPolicy: corev1.RestartPolicyNever,
						},
					},
				},
			},
		},
	}

	return cronjob
}

// SetupWithManager sets up the controller with the Manager.
func (r *FetchDataReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&curatorv1alpha1.FetchData{}).
		Complete(r)
}
