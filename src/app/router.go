package app

import (
	"github.com/teambition/gear"

	"github.com/open-trust/ot-ac/src/api"
	"github.com/open-trust/ot-ac/src/bll"
	"github.com/open-trust/ot-ac/src/conf"
	"github.com/open-trust/ot-ac/src/middleware"
	"github.com/open-trust/ot-ac/src/util"
)

func init() {
	util.DigProvide(NewRouters)
}

func getVersion(ctx *gear.Context) error {
	return ctx.OkJSON(conf.AppInfo())
}

// NewRouters ...
func NewRouters(apis *api.APIs) []*gear.Router {

	router := gear.NewRouter()
	router.Use(bll.IntoCtx(apis.Blls))

	router.Get("/", getVersion)
	router.Get("/version", getVersion)
	router.Get("/livez", apis.Healthz.Check)
	router.Get("/readyz", apis.Healthz.Check)

	router.Post("/AC/:action", middleware.VerifyTenant, apis.AC.Serve)
	router.Post("/Object/:action", middleware.VerifyTenant, apis.Object.Serve)
	router.Post("/Unit/:action", middleware.VerifyTenant, apis.Unit.Serve)

	router.Post("/Scope/Add", middleware.VerifyTenant, apis.Scope.Add)
	router.Post("/Scope/Delete", middleware.VerifyTenant, apis.Scope.Delete)
	router.Post("/Scope/DeleteAll", middleware.VerifyTenant, apis.Scope.DeleteAll)
	router.Post("/Scope/UpdateStatus", middleware.VerifyTenant, apis.Scope.UpdateStatus)
	router.Post("/Scope/List", middleware.VerifyTenant, apis.Scope.List)
	router.Post("/Scope/ListUnits", middleware.VerifyTenant, apis.Scope.ListUnits)
	router.Post("/Scope/ListObjects", middleware.VerifyTenant, apis.Scope.ListObjects)

	router.Post("/Permission/BatchAdd", middleware.VerifyTenant, apis.Permission.BatchAdd)
	router.Post("/Permission/List", middleware.VerifyTenant, apis.Permission.List)
	router.Post("/Permission/Delete", middleware.VerifyTenant, apis.Permission.Delete)

	// Admin
	router.Post("/Admin/AddTenant", middleware.VerifyAdmin, apis.Admin.AddTenant)
	router.Post("/Admin/UpdateTenantStatus", middleware.VerifyAdmin, apis.Admin.UpdateTenantStatus)
	router.Post("/Admin/DeleteTenant", middleware.VerifyAdmin, apis.Admin.DeleteTenant)
	router.Post("/Admin/ListTenants", middleware.VerifyAdmin, apis.Admin.ListTenants)
	router.Post("/Admin/BatchAddSubjects", middleware.VerifyAdmin, apis.Admin.BatchAddSubjects)
	router.Post("/Admin/UpdateSubjectStatus", middleware.VerifyAdmin, apis.Admin.UpdateSubjectStatus)
	router.Post("/Admin/ListSubjects", middleware.VerifyAdmin, apis.Admin.ListSubjects)
	// router.Post("/Admin/AddOrg", middleware.VerifyAdmin, apis.Admin.AddOrg)
	// router.Post("/Admin/UpdateOrgStatus", middleware.VerifyAdmin, apis.Admin.UpdateOrgStatus)
	// router.Post("/Admin/DeleteOrg", middleware.VerifyAdmin, apis.Admin.DeleteOrg)
	// router.Post("/Admin/AddMember", middleware.VerifyAdmin, apis.Admin.AddMember)
	// router.Post("/Admin/UpdateMemberStatus", middleware.VerifyAdmin, apis.Admin.UpdateMemberStatus)
	// router.Post("/Admin/DeleteMember", middleware.VerifyAdmin, apis.Admin.DeleteMember)
	// router.Post("/Admin/AddOU", middleware.VerifyAdmin, apis.Admin.AddOU)
	// router.Post("/Admin/UpdateOUStatus", middleware.VerifyAdmin, apis.Admin.UpdateOUStatus)
	// router.Post("/Admin/DeleteOU", middleware.VerifyAdmin, apis.Admin.DeleteOU)
	// router.Post("/Admin/AddOUMember", middleware.VerifyAdmin, apis.Admin.AddOUMember)
	// router.Post("/Admin/RemoveOUMember", middleware.VerifyAdmin, apis.Admin.RemoveOUMember)

	return []*gear.Router{router}
}
