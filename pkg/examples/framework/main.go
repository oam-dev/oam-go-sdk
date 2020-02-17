package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"reflect"

	"github.com/oam-dev/oam-go-sdk/apis/core.oam.dev/v1alpha1"
	"github.com/oam-dev/oam-go-sdk/pkg/client/clientset/versioned"
	"github.com/oam-dev/oam-go-sdk/pkg/oam"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	ctrl.SetLogger(zap.Logger(true))
	_ = v1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.Parse()
	options := ctrl.Options{Scheme: scheme, MetricsBindAddress: metricsAddr}
	// init
	oam.InitMgr(ctrl.GetConfigOrDie(), options)
	clientset, err := kubernetes.NewForConfig(ctrl.GetConfigOrDie())
	if err != nil {
		panic(err)
	}
	oamclient, err := versioned.NewForConfig(ctrl.GetConfigOrDie())
	if err != nil {
		log.Fatal("create client err: ", err)
	}
	// register workloadtpye & trait hooks and handlers
	oam.RegisterHandlers(oam.STypeApplicationConfiguration, &Handler{name: "my-handler", oamclient: oamclient, k8sclient: clientset})

	// reconcilers must register manualy
	// cloudnativeapp/oam-runtime/pkg/oam as a pkg should not do os.Exit(), instead of
	// panic or returning Error could be better
	err = oam.Run(oam.WithApplicationConfiguration())
	if err != nil {
		panic(err)
	}
}

type Handler struct {
	name      string
	oamclient *versioned.Clientset
	k8sclient *kubernetes.Clientset
}

func getManuelScale(traits []v1alpha1.TraitBinding) *int32 {
	var def int32 = 1
	for _, tr := range traits {
		if tr.Name != "manual-scaler" {
			continue
		}
		values := make(map[string]interface{})
		err := json.Unmarshal(tr.Properties.Raw, &values)
		if err != nil {
			log.Println("traits value spec error")
			continue
		}
		f, ok := values["replicaCount"]
		if !ok {
			log.Println("replicaCount didn't exist error")
			continue
		}
		ff, ok := f.(float64)
		if !ok {
			log.Println("replicaCount type is " + reflect.TypeOf(f).Name())
			continue
		}
		def = int32(ff)
	}
	return &def
}

func (s *Handler) Handle(ctx *oam.ActionContext, obj runtime.Object, eType oam.EType) error {
	ac, ok := obj.(*v1alpha1.ApplicationConfiguration)
	if !ok {
		return errors.New("type mismatch")
	}
	setupLog.Info("oam handler: " + s.name + " received ApplicationConfiguration " + ac.Name)
	for _, compConf := range ac.Spec.Components {
		comp, err := s.oamclient.CoreV1alpha1().ComponentSchematics(ac.Namespace).Get(compConf.ComponentName, v1.GetOptions{})
		if err != nil {
			return err
		}
		switch comp.Spec.WorkloadType {
		case "core.oam.dev/v1alpha1.Server":
			// for example, we create K8s deployment here for core.oam.dev/v1alpha1.Server workload
			deploymentsClient := s.k8sclient.AppsV1().Deployments(ac.Namespace)
			deployment := &appsv1.Deployment{
				ObjectMeta: v1.ObjectMeta{
					Name: compConf.InstanceName,
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: getManuelScale(compConf.Traits),
					Selector: &v1.LabelSelector{
						MatchLabels: map[string]string{
							"app": "demo",
						},
					},
					Template: apiv1.PodTemplateSpec{
						ObjectMeta: v1.ObjectMeta{
							Labels: map[string]string{
								"app": "demo",
							},
						},
						Spec: apiv1.PodSpec{
							Containers: []apiv1.Container{
								{
									Name:  comp.Spec.Containers[0].Name,
									Image: comp.Spec.Containers[0].Image,
									Ports: []apiv1.ContainerPort{
										{
											Name:          comp.Spec.Containers[0].Ports[0].Name,
											Protocol:      apiv1.Protocol(comp.Spec.Containers[0].Ports[0].Protocol),
											ContainerPort: comp.Spec.Containers[0].Ports[0].ContainerPort,
										},
									},
								},
							},
						},
					},
				},
			}
			fmt.Println("Creating deployment...")
			result, err := deploymentsClient.Create(deployment)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

		default:
			//You could launch you own CRD here according to workloadType
			return errors.New(comp.Spec.WorkloadType + " is undefined")
		}
	}
	return nil
}

func (s *Handler) Id() string {
	return "Handler"
}
