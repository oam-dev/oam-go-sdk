package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/oam-dev/oam-go-sdk/apis/common"

	"github.com/oam-dev/oam-go-sdk/pkg/client/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"log"

	"github.com/oam-dev/oam-go-sdk/apis/core.oam.dev/v1alpha1"
	"github.com/oam-dev/oam-go-sdk/pkg/oam"
	"k8s.io/apimachinery/pkg/runtime"
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
	client, err := versioned.NewForConfig(ctrl.GetConfigOrDie())
	if err != nil {
		log.Fatal("create client err: ", err)
	}
	// register workloadtpye & trait hooks and handlers
	oam.RegisterHandlers(oam.STypeApplicationConfiguration, &Handler{name: "app", client: client})

	// reconcilers must register manualy
	// cloudnativeapp/oam-runtime/pkg/oam as a pkg should not do os.Exit(), instead of
	// panic or returning Error could be better
	err = oam.Run(oam.WithApplicationConfiguration())
	if err != nil {
		panic(err)
	}
}

type Handler struct {
	client *versioned.Clientset
	name   string
}

func (s *Handler) Handle(ctx *oam.ActionContext, comp runtime.Object, eType oam.EType) error {
	appConfig, ok := comp.(*v1alpha1.ApplicationConfiguration)
	if !ok {
		return errors.New("type mismatch")
	}
	for _, comp := range appConfig.Spec.Components {
		compIns, err := s.client.CoreV1alpha1().ComponentSchematics(appConfig.Namespace).Get(comp.ComponentName, v1.GetOptions{})
		if err != nil {
			return fmt.Errorf("get component %s err %v", comp.ComponentName, err)
		}
		settings, err := common.ExtractParams(comp.ParameterValues, compIns.Spec.WorkloadSettings)
		if err != nil {
			return err
		}
		for k, v := range settings {
			fmt.Printf("%s: %s\n", k, v)
		}
	}
	return nil
}

func (s *Handler) Id() string {
	return "Handler"
}
