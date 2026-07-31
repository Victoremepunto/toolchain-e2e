package resources

import (
	"context"
	"fmt"

	ctemplate "github.com/codeready-toolchain/toolchain-common/pkg/template"
	"github.com/codeready-toolchain/toolchain-e2e/setup/templates"
	"github.com/codeready-toolchain/toolchain-e2e/setup/wait"

	templatev1 "github.com/openshift/api/template/v1"
	"k8s.io/apimachinery/pkg/runtime"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const userNSParam = "CURRENT_USER_NAMESPACE"

var tmpls map[string]*templatev1.Template = make(map[string]*templatev1.Template)

func CreateUserResourcesFromTemplateFiles(ctx context.Context, cl runtimeclient.Client, s *runtime.Scheme, username string, templatePaths []string) error {
	// load and validate all templates first so errors surface before the space wait
	for _, templatePath := range templatePaths {
		if _, ok := tmpls[templatePath]; !ok {
			var err error
			if tmpls[templatePath], err = templates.GetTemplateFromFile(templatePath); err != nil {
				return fmt.Errorf("invalid template file: '%s': %w", templatePath, err)
			}
		}
	}

	resolvedName, err := wait.ForSpaceWithName(cl, username)
	if err != nil {
		return err
	}
	userNS := fmt.Sprintf("%s-dev", resolvedName)
	combinedObjsToProcess := []runtimeclient.Object{}
	for _, templatePath := range templatePaths {
		tmpl := tmpls[templatePath]
		processor := ctemplate.NewProcessor(s)
		objsToProcess, err := processor.Process(tmpl.DeepCopy(), map[string]string{
			userNSParam: userNS,
		})
		if err != nil {
			return err
		}
		combinedObjsToProcess = append(combinedObjsToProcess, objsToProcess...)
	}

	if len(combinedObjsToProcess) == 0 {
		return fmt.Errorf("no objects found in templates %v", templatePaths)
	}

	return templates.ApplyObjectsConcurrently(ctx, cl, combinedObjsToProcess, templates.NamespaceModifier(userNS))
}
