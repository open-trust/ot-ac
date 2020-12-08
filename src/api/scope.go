package api

import (
	"github.com/open-trust/ot-ac/src/bll"
	"github.com/teambition/gear"
)

// Scope ..
type Scope struct {
	blls *bll.Blls
}

// Add 创建范围约束
func (a *Scope) Add(ctx *gear.Context) error {
	return nil
}

// Delete 删除范围约束
func (a *Scope) Delete(ctx *gear.Context) error {
	return nil
}

// DeleteAll 删除范围约束及范围内的所有 Unit 和 Object
func (a *Scope) DeleteAll(ctx *gear.Context) error {
	return nil
}

// UpdateStatus 更新范围约束的状态，-1 表示停用
func (a *Scope) UpdateStatus(ctx *gear.Context) error {
	return nil
}

// List 列出该系统当前所有指定目标类型的范围约束
func (a *Scope) List(ctx *gear.Context) error {
	return nil
}

// ListUnits 列出范围约束下指定目标类型的直属的管理单元
func (a *Scope) ListUnits(ctx *gear.Context) error {
	return nil
}

// ListObjects 列出范围约束下指定目标类型的直属的资源对象
func (a *Scope) ListObjects(ctx *gear.Context) error {
	return nil
}
