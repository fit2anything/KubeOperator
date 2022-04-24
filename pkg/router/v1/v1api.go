package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KubeOperator/KubeOperator/pkg/controller"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/middleware"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/mvc"
)

var AuthScope iris.Party
var WhiteScope iris.Party
var log = logger.Default

func V1(parent iris.Party) {
	v1 := parent.Party("/v1")
	authParty := v1.Party("/auth")
	authParty.Use(middleware.LanguageMiddleware)
	mvc.New(authParty.Party("/session")).HandleError(ErrorHandler).Handle(controller.NewSessionController())
	AuthScope = v1.Party("/")
	authParty.Use(middleware.LanguageMiddleware)
	AuthScope.Use(middleware.SessionMiddleware)
	AuthScope.Use(middleware.RBACMiddleware())
	AuthScope.Use(middleware.PagerMiddleware)
	AuthScope.Use(middleware.ForceMiddleware)
	helperAuthScope := AuthScope.Party("/")
	helperAuthScope.Use(middleware.HelperMiddleware)
	mvc.New(helperAuthScope.Party("/clusters")).HandleError(ErrorHandler).Handle(controller.NewClusterHelperController())
	mvc.New(AuthScope.Party("/clusters")).HandleError(ErrorHandler).Handle(controller.NewClusterController())
	mvc.New(AuthScope.Party("/credentials")).HandleError(ErrorHandler).Handle(controller.NewCredentialController())
	mvc.New(AuthScope.Party("/hosts")).HandleError(ErrorHandler).Handle(controller.NewHostController())
	mvc.New(AuthScope.Party("/users")).HandleError(ErrorHandler).Handle(controller.NewUserController())
	mvc.New(AuthScope.Party("/settings")).HandleError(ErrorHandler).Handle(controller.NewSystemSettingController())
	mvc.New(AuthScope.Party("/logs")).HandleError(ErrorHandler).Handle(controller.NewSystemLogController())
	mvc.New(AuthScope.Party("/projects")).HandleError(ErrorHandler).Handle(controller.NewProjectController())
	mvc.New(AuthScope.Party("/clusters/istio")).HandleError(ErrorHandler).Handle(controller.NewClusterIstioController())
	mvc.New(AuthScope.Party("/clusters/kubernetes")).HandleError(ErrorHandler).Handle(controller.NewKubernetesController())
	mvc.New(AuthScope.Party("/backupaccounts")).HandleError(ErrorHandler).Handle(controller.NewBackupAccountController())
	mvc.New(AuthScope.Party("/clusters/backup")).HandleError(ErrorHandler).Handle(controller.NewClusterBackupStrategyController())
	mvc.New(AuthScope.Party("/license")).Handle(ErrorHandler).Handle(controller.NewLicenseController())
	mvc.New(AuthScope.Party("/clusters/backup/files")).HandleError(ErrorHandler).Handle(controller.NewClusterBackupFileController())
	mvc.New(AuthScope.Party("/manifests")).HandleError(ErrorHandler).Handle(controller.NewManifestController())
	mvc.New(AuthScope.Party("/vm/configs")).HandleError(ErrorHandler).Handle(controller.NewVmConfigController())
	mvc.New(AuthScope.Party("/events")).HandleError(ErrorHandler).Handle(controller.NewClusterEventController())
	mvc.New(AuthScope.Party("/project/resources")).HandleError(ErrorHandler).Handle(controller.NewProjectResourceController())
	mvc.New(AuthScope.Party("/project/members")).HandleError(ErrorHandler).Handle(controller.NewProjectMemberController())
	WhiteScope = v1.Party("/")
	WhiteScope.Get("/captcha", generateCaptcha)
	mvc.New(WhiteScope.Party("/theme")).HandleError(ErrorHandler).Handle(controller.NewThemeController())

}

func ErrorHandler(ctx context.Context, err error) {
	if err != nil {
		warp := struct {
			Msg string `json:"msg"`
		}{err.Error()}
		var result string
		switch errType := err.(type) {
		case gorm.Errors:
			errorSet := make(map[string]string)
			for _, er := range errType {
				tr := ctx.Tr(er.Error())
				if tr != "" {
					errorMsg := tr
					errorSet[er.Error()] = errorMsg
				}
			}
			for _, set := range errorSet {
				result = result + set + " "
			}
		case error:
			tr := ctx.Tr(err.Error())
			if tr != "" {
				result = tr
			} else {
				result = err.Error()
			}
		default:
			fmt.Printf("err type is %T", err)
		}
		warp.Msg = result
		bf, err := json.Marshal(&warp)
		if err != nil {
			log.Errorf("json marshal failed, %v", warp)
		}
		ctx.StatusCode(http.StatusBadRequest)
		_, _ = ctx.WriteString(string(bf))
		ctx.StopExecution()
		return
	}
}
