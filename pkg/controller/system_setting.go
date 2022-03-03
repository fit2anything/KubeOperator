package controller

import (
	"errors"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/kolog"
	"github.com/KubeOperator/KubeOperator/pkg/controller/koregexp"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
)

var (
	RegistryAlreadyExistsErr = errors.New("REGISTRY_ALREADY_EXISTS")
)

type SystemSettingController struct {
	Ctx                  context.Context
	SystemSettingService service.SystemSettingService
}

func NewSystemSettingController() *SystemSettingController {
	return &SystemSettingController{
		SystemSettingService: service.NewSystemSettingService(),
	}
}

func (s SystemSettingController) Get() (interface{}, error) {
	item, err := s.SystemSettingService.List()
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (s SystemSettingController) GetBy(name string) (interface{}, error) {
	item, err := s.SystemSettingService.ListByTab(name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (s SystemSettingController) Post() ([]dto.SystemSetting, error) {
	var req dto.SystemSettingCreate
	err := s.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return nil, err
	}
	result, err := s.SystemSettingService.Create(req)
	if err != nil {
		return nil, err
	}

	go kolog.Save(s.Ctx, constant.CREATE_EMAIL, "-")

	return result, nil
}

func (s SystemSettingController) GetRegistry() (page.Page, error) {
	p, _ := s.Ctx.Values().GetBool("page")
	if p {
		num, _ := s.Ctx.Values().GetInt(constant.PageNumQueryKey)
		size, _ := s.Ctx.Values().GetInt(constant.PageSizeQueryKey)
		return s.SystemSettingService.PageRegistry(num, size)
	} else {
		var page page.Page
		items, err := s.SystemSettingService.ListRegistry()
		if err != nil {
			return page, err
		}
		page.Items = items
		page.Total = len(items)
		return page, nil
	}

}

func (s SystemSettingController) GetRegistryBy(arch string) (interface{}, error) {
	item, err := s.SystemSettingService.GetRegistryByArch(arch)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (s SystemSettingController) PostRegistry() (*dto.SystemRegistry, error) {
	var req dto.SystemRegistryCreate
	err := s.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	if err := validate.RegisterValidation("koip", koregexp.CheckIpPattern); err != nil {
		return nil, err
	}
	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	item, _ := s.SystemSettingService.GetRegistryByArch(req.Architecture)
	if item.ID != "" {
		return nil, RegistryAlreadyExistsErr
	}

	go kolog.Save(s.Ctx, constant.CREATE_REGISTRY, req.Architecture)

	return s.SystemSettingService.CreateRegistry(req)
}

func (s SystemSettingController) PatchRegistryBy(arch string) (*dto.SystemRegistry, error) {
	var req dto.SystemRegistryUpdate
	err := s.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}
	validate := validator.New()
	if err := validate.RegisterValidation("koip", koregexp.CheckIpPattern); err != nil {
		return nil, err
	}
	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	go kolog.Save(s.Ctx, constant.UPDATE_REGISTRY, req.Architecture)

	return s.SystemSettingService.UpdateRegistry(req)
}

func (s SystemSettingController) PostRegistryBatch() error {
	var req dto.SystemRegistryBatchOp
	err := s.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		return err
	}
	err = s.SystemSettingService.BatchRegistry(req)
	if err != nil {
		return err
	}
	delCres := ""
	for _, item := range req.Items {
		delCres += (item.Architecture + ",")
	}
	go kolog.Save(s.Ctx, constant.DELETE_REGISTRY, delCres)
	return err
}
