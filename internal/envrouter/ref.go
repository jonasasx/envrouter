package envrouter

import (
	"errors"
	"github.com/ghodss/yaml"
	"gitlab.com/jonasasx/envrouter/internal/envrouter/api"
	"gitlab.com/jonasasx/envrouter/internal/envrouter/k8s"
	"sort"
)

type RefService interface {
	SaveBinding(refBinding *api.RefBinding) (*api.RefBinding, error)
	FindAllBindings(environmentFilter *string, applicationFilter *string, refFilter *string) ([]*api.RefBinding, error)
}

type refService struct {
	dataStorage        k8s.ConfigMapDataStorage
	environmentService EnvironmentService
	applicationService ApplicationService
}

func NewRefService(
	dataStorage k8s.ConfigMapDataStorage,
	environmentService EnvironmentService,
	applicationService ApplicationService,
) RefService {
	return &refService{
		dataStorage,
		environmentService,
		applicationService,
	}
}

func (r *refService) SaveBinding(refBinding *api.RefBinding) (*api.RefBinding, error) {
	if !r.environmentService.ExistsByName(refBinding.Environment) {
		return nil, errors.New("Environment " + refBinding.Environment + " is not found")
	}
	if !r.applicationService.ExistsByName(refBinding.Application) {
		return nil, errors.New("Application " + refBinding.Application + " is not found")
	}
	value, err := yaml.Marshal(refBinding)
	if err != nil {
		return nil, err
	}
	return refBinding, r.dataStorage.Save(refBinding.Environment+"-"+refBinding.Application, string(value))
}

func (r *refService) FindAllBindings(environmentFilter *string, applicationFilter *string, refFilter *string) ([]*api.RefBinding, error) {
	data, err := r.dataStorage.GetAll()
	if err != nil {
		return nil, err
	}
	bindings := map[string]string{}
	for _, v := range data {
		item := api.RefBinding{}
		err := yaml.Unmarshal([]byte(v), &item)
		if err != nil {
			return nil, err
		}
		bindings[item.Environment+"-"+item.Application] = item.Ref
	}
	environments, err := r.environmentService.FindAll()
	if err != nil {
		return nil, err
	}
	applications, err := r.applicationService.FindAll()
	if err != nil {
		return nil, err
	}
	result := []*api.RefBinding{}
	for _, environment := range environments {
		if environmentFilter != nil && *environmentFilter != "" && *environmentFilter != environment.Name {
			continue
		}
		for _, application := range applications {
			if applicationFilter != nil && *applicationFilter != "" && *applicationFilter != application.Name {
				continue
			}
			var ref string
			if v, ok := bindings[environment.Name+"-"+application.Name]; ok {
				ref = v
			} else {
				ref = DefaultRef
			}
			if refFilter != nil && *refFilter != "" && *refFilter != ref {
				continue
			}
			item := api.RefBinding{
				Environment: environment.Name,
				Application: application.Name,
				Ref:         ref,
			}
			result = append(result, &item)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Environment < result[j].Environment && result[i].Application < result[j].Application
	})
	return result, nil
}
