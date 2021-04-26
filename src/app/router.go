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
	router.Get("/healthz", apis.Healthz.Check)

	router.Post("/AC/CheckUnit", middleware.VerifyTenant, apis.AC.CheckUnit)
	router.Post("/AC/CheckScope", middleware.VerifyTenant, apis.AC.CheckScope)
	router.Post("/AC/CheckObject", middleware.VerifyTenant, apis.AC.CheckObject)
	router.Post("/AC/ListPermissionsByUnit", middleware.VerifyTenant, apis.AC.ListPermissionsByUnit)
	router.Post("/AC/ListPermissionsByScope", middleware.VerifyTenant, apis.AC.ListPermissionsByScope)
	router.Post("/AC/ListPermissionsByObject", middleware.VerifyTenant, apis.AC.ListPermissionsByObject)
	router.Post("/AC/ListObject", middleware.VerifyTenant, apis.AC.ListObject)
	router.Post("/AC/SearchObject", middleware.VerifyTenant, apis.AC.SearchObject)

	router.Post("/Object/BatchAdd", middleware.VerifyTenant, apis.Object.BatchAdd)
	router.Post("/Object/AssignParent", middleware.VerifyTenant, apis.Object.AssignParent)
	router.Post("/Object/AssignScope", middleware.VerifyTenant, apis.Object.AssignScope)
	router.Post("/Object/RemoveParent", middleware.VerifyTenant, apis.Object.RemoveParent)
	router.Post("/Object/RemoveScope", middleware.VerifyTenant, apis.Object.RemoveScope)
	router.Post("/Object/Delete", middleware.VerifyTenant, apis.Object.Delete)
	router.Post("/Object/UpdateTerms", middleware.VerifyTenant, apis.Object.UpdateTerms)
	router.Post("/Object/AddPermissions", middleware.VerifyTenant, apis.Object.AddPermissions)
	router.Post("/Object/UpdatePermissions", middleware.VerifyTenant, apis.Object.UpdatePermissions)
	router.Post("/Object/RemovePermissions", middleware.VerifyTenant, apis.Object.RemovePermissions)
	router.Post("/Object/ListChildren", middleware.VerifyTenant, apis.Object.ListChildren)
	router.Post("/Object/ListDescendant", middleware.VerifyTenant, apis.Object.ListDescendant)
	router.Post("/Object/ListPermissions", middleware.VerifyTenant, apis.Object.ListPermissions)
	router.Post("/Object/GetDAG", middleware.VerifyTenant, apis.Object.GetDAG)
	router.Post("/Object/Search", middleware.VerifyTenant, apis.Object.Search)

	router.Post("/Unit/BatchAdd", middleware.VerifyTenant, apis.Unit.BatchAdd)
	router.Post("/Unit/AddFromOrg", middleware.VerifyTenant, apis.Unit.AddFromOrg)
	router.Post("/Unit/AddFromOU", middleware.VerifyTenant, apis.Unit.AddFromOU)
	router.Post("/Unit/AddFromMembers", middleware.VerifyTenant, apis.Unit.AddFromMembers)
	router.Post("/Unit/AssignParent", middleware.VerifyTenant, apis.Unit.AssignParent)
	router.Post("/Unit/AssignScope", middleware.VerifyTenant, apis.Unit.AssignScope)
	router.Post("/Unit/AssignObject", middleware.VerifyTenant, apis.Unit.AssignObject)
	router.Post("/Unit/RemoveParent", middleware.VerifyTenant, apis.Unit.RemoveParent)
	router.Post("/Unit/RemoveScope", middleware.VerifyTenant, apis.Unit.RemoveScope)
	router.Post("/Unit/RemoveObject", middleware.VerifyTenant, apis.Unit.RemoveObject)
	router.Post("/Unit/Delete", middleware.VerifyTenant, apis.Unit.Delete)
	router.Post("/Unit/UpdateStatus", middleware.VerifyTenant, apis.Unit.UpdateStatus)
	router.Post("/Unit/AddSubjects", middleware.VerifyTenant, apis.Unit.AddSubjects)
	router.Post("/Unit/RemoveSubjects", middleware.VerifyTenant, apis.Unit.RemoveSubjects)
	router.Post("/Unit/AddPermissions", middleware.VerifyTenant, apis.Unit.AddPermissions)
	router.Post("/Unit/UpdatePermissions", middleware.VerifyTenant, apis.Unit.UpdatePermissions)
	router.Post("/Unit/RemovePermissions", middleware.VerifyTenant, apis.Unit.RemovePermissions)
	router.Post("/Unit/ListChildren", middleware.VerifyTenant, apis.Unit.ListChildren)
	router.Post("/Unit/ListDescendant", middleware.VerifyTenant, apis.Unit.ListDescendant)
	router.Post("/Unit/ListPermissions", middleware.VerifyTenant, apis.Unit.ListPermissions)
	router.Post("/Unit/ListSubjects", middleware.VerifyTenant, apis.Unit.ListSubjects)
	router.Post("/Unit/ListDescendantSubjects", middleware.VerifyTenant, apis.Unit.ListDescendantSubjects)
	router.Post("/Unit/GetDAG", middleware.VerifyTenant, apis.Unit.GetDAG)

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

	// Organization
	router.Post("/Organization/AddOrg", middleware.VerifyAdmin, apis.Organization.AddOrg)
	router.Post("/Organization/UpdateOrgStatus", middleware.VerifyAdmin, apis.Organization.UpdateOrgStatus)
	router.Post("/Organization/DeleteOrg", middleware.VerifyAdmin, apis.Organization.DeleteOrg)
	router.Post("/Organization/ListOrgs", middleware.VerifyAdmin, apis.Organization.ListOrgs)
	router.Post("/Organization/ListSubjectOrgs", middleware.VerifyAdmin, apis.Organization.ListSubjectOrgs)

	router.Post("/Organization/AddOU", middleware.VerifyAdmin, apis.Organization.AddOU)
	router.Post("/Organization/UpdateOUParent", middleware.VerifyAdmin, apis.Organization.UpdateOUParent)
	router.Post("/Organization/UpdateOUStatus", middleware.VerifyAdmin, apis.Organization.UpdateOUStatus)
	router.Post("/Organization/UpdateOUTerms", middleware.VerifyAdmin, apis.Organization.UpdateOUTerms)
	router.Post("/Organization/DeleteOU", middleware.VerifyAdmin, apis.Organization.DeleteOU)
	router.Post("/Organization/ListOUs", middleware.VerifyAdmin, apis.Organization.ListOUs)
	router.Post("/Organization/ListSubjectOUs", middleware.VerifyAdmin, apis.Organization.ListSubjectOUs)
	router.Post("/Organization/SearchOUs", middleware.VerifyAdmin, apis.Organization.SearchOUs)

	router.Post("/Organization/BatchAddMember", middleware.VerifyAdmin, apis.Organization.BatchAddMember)
	router.Post("/Organization/UpdateMemberStatus", middleware.VerifyAdmin, apis.Organization.UpdateMemberStatus)
	router.Post("/Organization/UpdateMemberTerms", middleware.VerifyAdmin, apis.Organization.UpdateMemberTerms)
	router.Post("/Organization/RemoveMember", middleware.VerifyAdmin, apis.Organization.RemoveMember)
	router.Post("/Organization/ListMembers", middleware.VerifyAdmin, apis.Organization.ListMembers)
	router.Post("/Organization/SearchMember", middleware.VerifyAdmin, apis.Organization.SearchMember)

	router.Post("/Organization/BatchAddOUMember", middleware.VerifyAdmin, apis.Organization.BatchAddOUMember)
	router.Post("/Organization/RemoveOUMember", middleware.VerifyAdmin, apis.Organization.RemoveOUMember)
	router.Post("/Organization/ListOUMembers", middleware.VerifyAdmin, apis.Organization.ListOUMembers)
	router.Post("/Organization/ListOUDescendantMembers", middleware.VerifyAdmin, apis.Organization.ListOUDescendantMembers)

	return []*gear.Router{router}
}
