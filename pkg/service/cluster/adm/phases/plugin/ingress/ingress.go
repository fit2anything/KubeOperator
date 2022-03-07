package ingress

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/facts"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	ingressPlaybook = "14-ingress-controller.yml"
)

type ControllerPhase struct {
	IngressControllerType string
}

func (ControllerPhase) Name() string {
	return "IngressController"
}

func (c ControllerPhase) Run(b kobe.Interface, fileName string) error {
	if c.IngressControllerType != "" {
		b.SetVar(facts.IngressControllerTypeFactName, c.IngressControllerType)
	}
	return phases.RunPlaybookAndGetResult(b, ingressPlaybook, "", fileName)
}
