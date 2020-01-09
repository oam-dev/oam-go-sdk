package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"

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

	// register workloadtpye & trait hooks and handlers
	oam.RegisterPreHooks(oam.STypeComponent, &PreHook{name: "comp"})
	oam.RegisterHandlers(oam.STypeComponent, &Handler{name: "comp"})
	oam.RegisterPostHooks(oam.STypeComponent, &PostHook{name: "comp"})

	// reconcilers must register manualy
	// cloudnativeapp/oam-runtime/pkg/oam as a pkg should not do os.Exit(), instead of
	// panic or returning Error could be better
	err := oam.Run(oam.WithComponent())
	if err != nil {
		panic(err)
	}
}

type PreHook struct {
	name string
}

type Handler struct {
	name string
}

type PostHook struct {
	name string
}

func (p *PreHook) Exec(ctx *oam.ActionContext, comp runtime.Object, ev oam.EType) error {
	setupLog.Info("hello oam from pre hook: " + p.name)
	return nil
}

func (e *PreHook) Id() string {
	return "PreHook"
}

func (s *Handler) Handle(ctx *oam.ActionContext, comp runtime.Object, eType oam.EType) error {
	component, ok := comp.(*v1alpha1.ComponentSchematic)
	if !ok {
		return errors.New("type mismatch")
	}
	fmt.Printf("settings: %s\n", component.Spec.WorkloadSettings.Raw)
	//Note: this type should be consistent with workloadType
	settings := make(map[string]interface{})
	err := json.Unmarshal(component.Spec.WorkloadSettings.Raw, &settings)
	if err != nil {
		return err
	}
	for k, v := range settings {
		fmt.Printf("%s: %s\n", k, v)
	}
	return nil
}

func (s *Handler) Id() string {
	return "Handler"
}

func (p *PostHook) Exec(ctx *oam.ActionContext, comp runtime.Object, ev oam.EType) error {
	setupLog.Info("hello oam from post hook: " + p.name)
	return nil
}

func (e *PostHook) Id() string {
	return "PostHook"
}
