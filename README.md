# ot-ac
Open Trust Access Control service.

[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/open-trust/ot-ac/master/LICENSE)

## Features

## Concepts

[Concepts](https://github.com/open-trust/ot-ac/blob/master/doc/concepts.md)

## Documentation

[API](https://github.com/open-trust/ot-ac/blob/master/doc/openapi.md)

核心 API 设计

说明

1. 所有写操作都应该幂等，部分 API 可以通过 Header "Prefer: respond-conflict" 来声明对于写入资源冲突时响应 409 错误
2. 建议业务方创建 Permission("Default") 这一特殊类型的权限，可用于 Object 的权限中断场景或代表其它没有资源实体的默认操作权限
3. 本 API 列表尤其是 Search 相关的 API，并不要求 GBAC 系统全部实现，按照业务需求实现即可
4. List 类的 API 都支持基于游标的分页

类型

type Target {
  targetType: String
  targetId: String
}

type Permission `Resource.Operation.Constraint`

访问控制查询

// 检查请求主体到指定管理单元有没有指定权限
CheckUnit(subject: String!, unit: Target!, permission: Permission!)

// 检查请求主体到指定范围约束有没有指定权限
CheckScope(subject: String!, scope: Target!, permission: Permission!)

// 检查请求主体通过 Scope 或 Unit-Object 的连接关系到指定资源对象有没有指定权限，如果 ignoreScope 为 true，则要求必须有 Unit-Object 的连接关系
CheckObject(subject: String!, object: Target!, permission: Permission!, ignoreScope: Boolean = false)

// 列出请求主体到指定管理单元的符合 resource 的权限，如果未指定管理单元，则会查询请求主体能触达的所有管理单元，如果 resources 为空，则会列出所有触达的有效权限
ListPermissionsByUnit(subject: String!, unit: Target = null, resources: [String])

// 列出请求主体到指定范围约束的符合 resource 的权限，如果 resources 为空，则会列出所有触达的有效权限
ListPermissionsByScope(subject: String!, scope: Target!, resources: [String])

// 列出请求主体到指定资源对象的符合 resource 的权限，如果 resources 为空，则会列出所有触达的有效权限
ListPermissionsByObject(subject: String!, object: Target!, resources: [String], ignoreScope: Boolean = false)

// 列出请求主体参与的指定类型的管理单元
ListUnits(subject: String!, targetType: String!)

// 列出请求主体在指定资源对象中能触达的所有指定类型的子孙资源对象
// depth 定义对 targetType 类型资源对象的递归查询深度，而不是指定 object 到 targetType 类型资源对象的深度，默认对 targetType 类型资源对象查到底
ListObjects(subject: String!, object: Target!, permission: Permission!, targetType: String!, ignoreScope: Boolean = false, depth: Int = MaxInt)

// 根据关键词，在指定资源对象的子孙资源对象中，对请求主体能触达的所有指定类型的资源对象中进行搜索，term 为空不匹配任何资源对象
SearchObjects(subject: String!, object: Target!, permission: Permission!, targetType: String!, term: String!, ignoreScope: Boolean = false)

Scope 范围约束
// 创建范围约束
Add(scope: Target!)

// 删除范围约束
Delete(scope: Target!)

// 删除范围约束及范围内的所有 Unit 和 Object
DeleteAll(scope: Target!)

// 更新范围约束的状态，-1 表示停用
UpdateStatus(scope: Target!, status: Int!)

// 列出该系统当前所有指定目标类型的范围约束
List(targetType: String!)

// 列出范围约束下指定目标类型的直属的管理单元
ListUnits(scope: Target!, targetType: String!)

// 列出范围约束下指定目标类型的直属的资源对象
ListObjects(scope: Target!, targetType: String!)

Subject 请求主体
// 批量添加请求主体
BatchAdd(subjects: [String]!)

// 更新请求主体，-1 表示停用
UpdateStatus(subject: String!, status: Int!)

Permission 权限
// 批量添加权限
BatchAdd(permissions: [Permission]!)

// 删除权限
Delete(permission: Permission!)

// 列出该系统当前指定资源类型的权限，当 resource 为空时列出所有权限
List(resources: [String])

Unit 管理单元

// 批量添加管理单元，当检测到将形成环时会返回 400 错误
BatchAdd(units: [Target]!, parent: Target = null, scope: Target = null)

// 建立管理单元与父级管理单元的关系，当检测到将形成环时会返回 400 错误
AssignParent(unit: Target!, parent: Target!)

// 建立管理单元与范围约束的关系
AssignScope(unit: Target!, scope: Target!)

// 建立管理单元与资源对象的关系
AssignObject(unit: Target!, object: Target!)

// 清除管理单元与父级对象的关系
RemoveParent(unit: Target!, parent: Target!)

// 清除管理单元与范围约束的关系
RemoveScope(unit: Target!, scope: Target!)

// 清除管理单元与资源对象的关系
RemoveObject(unit: Target!, object: Target!)

// 删除管理单元及其所有子孙管理单元和链接关系
Delete(unit: Target!)

// 更新管理单元的状态，-1 表示停用
UpdateStatus(unit: Target!, status: Int!)

// 管理单元批量添加请求主体，当请求主体不存在时会自动创建
AddSubjects(unit: Target!, subjects: [String]!)

// 管理单元批量移除请求主体
RemoveSubjects(unit: Target!, subjects: [String]!)

// 给管理单元添加权限，权限必须预先存在
AddPermissions(unit: Target!, permissions: [Permission])

// 覆盖管理单元的权限，权限必须预先存在，当 permissions 为空时会清空权限
UpdatePermissions(unit: Target!, permissions: [Permission])

// 移除管理单元的权限
RemovePermissions(unit: Target!, permissions: [Permission])

// 列出管理单元的指定目标类型的子级管理单元
ListChildren(unit: Target!, targetType: String!)

// 列出管理单元的指定目标类型的所有子孙管理单元
// depth 定义对 targetType 类型管理单元的递归查询深度，而不是指定 unit 到 targetType 类型管理单元的深度，默认对 targetType 类型管理单元查到底
ListDescendant(unit: Target!, targetType: String!, depth: Int = MaxInt)

// 列出管理单元的直属权限
ListPermissions(unit: Target!)

// 列出管理单元的直属请求主体，不包含 status 为 -1 的请求主体
ListSubjects(unit: Target!)

// 列出管理单元及子孙管理单元下所有的请求主体，不包含 status 为 -1 的请求主体
ListDescendantSubjects(unit: Target!)

// 根据 start 和 ends 找出一个 DAG，其中 start 可以为 Subject 或 Unit，ends 为 0 到多个 Unit
GetDAG(start: Target!, ends: [Target]!)

Object 资源对象

// 批量添加资源对象，当检测到将形成环时会返回 400 错误
BatchAdd(objects: [Target]!, parent: Target = null, scope: Target = null)

// 建立资源对象与父级对象的关系，当检测到将会形成环时会返回 400 错误
AssignParent(object: Target!, parent: Target!)

// 建立资源对象与范围约束的关系
AssignScope(object: Target!, scope: Target!)

// 清除资源对象与父级对象的关系
RemoveParent(object: Target!, parent: Target!)

// 清除资源对象与范围约束的关系
RemoveScope(object: Target!, scope: Target!)

// 删除资源对象及其所有子孙资源对象和链接关系
Delete(object: Target!)

// 更新资源对象的搜索关键词
UpdateTerms(object: Target!, terms: [String]!)

// 给资源对象添加可透传的权限，权限必须预先存在
AddPermissions(object: Target!, permissions: [Permission])

// 覆盖资源对象可透传的权限，权限必须预先存在，当 permissions 为空时会清空权限
UpdatePermissions(object: Target!, permissions: [Permission])

// 移除资源对象可透传的权限
RemovePermissions(object: Target!, permissions: [Permission])

// 列出资源对象的指定目标类型的子级资源对象
ListChildren(object: Target!, targetType: String!)

// 列出资源对象的所有指定目标类型的子孙资源对象
// depth 定义对 targetType 类型资源对象的递归查询深度，而不是指定 object 到 targetType 类型资源对象的深度，默认对 targetType 类型资源对象查到底
ListDescendant(object: Target!, targetType: String!, depth: Int = MaxInt)

// 列出资源对象可透传的权限
ListPermissions(object: Target!)

// 根据 start 和 ends 找出一个 DAG，其中 start 为 Object，ends 为 0 到多个 Object
GetDAG(start: Target!, ends: [Target]!)

// 根据关键词在资源对象的所有指定类型的子孙资源对象中进行搜索，term 为空不匹配任何资源对象
Search(object: Target!, targetType: String!, term: String!)
