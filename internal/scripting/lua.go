package scripting

import (
	"fmt"
	"time"

	"github.com/bansikah22/kswp/internal/kubernetes"
	"github.com/bansikah22/kswp/internal/scanner"
	"github.com/bansikah22/kswp/pkg/models"
	lua "github.com/yuin/gopher-lua"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Execute(script string, client kubernetes.Client) error {
	L := lua.NewState()
	defer L.Close()

	// This is a temporary solution to avoid an import cycle.
	// TODO: Refactor ScanResources into a separate package.
	var resources []models.Resource
	listOptions := metav1.ListOptions{}
	unusedConfigMaps, err := scanner.GetUnusedConfigMaps(client.Clientset(), "", listOptions)
	if err != nil {
		return fmt.Errorf("error getting unused configmaps: %w", err)
	}
	resources = append(resources, unusedConfigMaps...)
	unusedSecrets, err := scanner.GetUnusedSecrets(client.Clientset(), "", listOptions)
	if err != nil {
		return fmt.Errorf("error getting unused secrets: %w", err)
	}
	resources = append(resources, unusedSecrets...)
	orphanServices, err := scanner.GetOrphanServices(client.Clientset(), "", listOptions)
	if err != nil {
		return fmt.Errorf("error getting orphan services: %w", err)
	}
	resources = append(resources, orphanServices...)
	oldReplicaSets, err := scanner.GetOldReplicaSets(client.Clientset(), "", listOptions)
	if err != nil {
		return fmt.Errorf("error getting old replicasets: %w", err)
	}
	resources = append(resources, oldReplicaSets...)
	completedJobs, err := scanner.GetCompletedJobs(client.Clientset(), 24*time.Hour, "", listOptions)
	if err != nil {
		return fmt.Errorf("error getting completed jobs: %w", err)
	}
	resources = append(resources, completedJobs...)
	failedPods, err := scanner.GetFailedPods(client.Clientset(), "", listOptions)
	if err != nil {
		return fmt.Errorf("error getting failed pods: %w", err)
	}
	resources = append(resources, failedPods...)
	completedPods, err := scanner.GetCompletedPods(client.Clientset(), 24*time.Hour, "", listOptions)
	if err != nil {
		return fmt.Errorf("error getting completed pods: %w", err)
	}
	resources = append(resources, completedPods...)
	if err != nil {
		return err
	}

	luaResources := L.NewTable()
	for _, res := range resources {
		resTable := L.NewTable()
		resTable.RawSetString("name", lua.LString(res.Name))
		resTable.RawSetString("namespace", lua.LString(res.Namespace))
		resTable.RawSetString("kind", lua.LString(res.Kind))
		resTable.RawSetString("reason", lua.LString(res.Reason))
		luaResources.Append(resTable)
	}

	L.SetGlobal("resources", luaResources)

	if err := L.DoString(script); err != nil {
		return err
	}

	return nil
}
