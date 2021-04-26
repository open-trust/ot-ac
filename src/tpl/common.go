package tpl

import (
	"regexp"
	"strconv"
	"strings"

	otgo "github.com/open-trust/ot-go-lib"
	"github.com/teambition/gear"
)

// ErrorResponseType 定义了标准的 API 接口错误时返回数据模型
type ErrorResponseType = gear.ErrorResponse

// SuccessResponseType 定义了标准的 API 接口成功时返回数据模型
type SuccessResponseType struct {
	TotalCount int         `json:"totalCount,omitempty"`
	NextToken  string      `json:"nextToken,omitempty"`
	Result     interface{} `json:"result"`
}

// ResponseType ...
type ResponseType struct {
	ErrorResponseType
	SuccessResponseType
}

// Pagination ...
type Pagination struct {
	PageToken string `json:"pageToken" query:"pageToken"`
	PageSize  int    `json:"pageSize" query:"pageSize"`
	Skip      int    `json:"skip" query:"skip"`
}

var uidReg = regexp.MustCompile(`^0x[0-9a-f]{1,16}$`)

// Validate ...
func (pg *Pagination) Validate() error {
	if pg.Skip < 0 {
		pg.Skip = 0
	}

	if pg.PageSize > 1000 {
		return gear.ErrBadRequest.WithMsgf("pageSize %v should not great than 1000", pg.PageSize)
	}

	if pg.PageSize <= 0 {
		pg.PageSize = 10
	}
	if pg.PageToken == "" {
		pg.PageToken = "0x0"
	}
	if !uidReg.MatchString(pg.PageToken) {
		return gear.ErrBadRequest.WithMsgf("invalid PageToken %s", strconv.Quote(pg.PageToken))
	}
	return nil
}

// Target ...
type Target struct {
	Type string `json:"targetType"`
	ID   string `json:"targetId"`
}

// Validate 实现 gear.BodyTemplate
func (t *Target) Validate() error {
	if err := CheckResource(t.Type); err != nil {
		return err
	}
	_, err := otgo.NewOTID("ac", t.Type, t.ID)
	if err != nil {
		return gear.ErrBadRequest.From(err)
	}
	return err
}

// ResourcesInput ...
type ResourcesInput struct {
	Resources []string `json:"resources"`
}

// Extensions ...
type Extensions map[string]interface{}

// Validate 实现 gear.BodyTemplate
func (t Extensions) Validate() error {
	for _, v := range t {
		switch v.(type) {
		case string, bool, int, int32, float64:
			continue
		default:
			return gear.ErrBadRequest.WithMsgf("unsupported extension value type: %T", v)
		}
	}
	return nil
}

// Validate 实现 gear.BodyTemplate
func (t *ResourcesInput) Validate() error {
	if t.Resources == nil {
		t.Resources = make([]string, 0)
	} else {
		for _, resource := range t.Resources {
			if err := CheckResource(resource); err != nil {
				return err
			}
		}
	}
	return nil
}

var resourceReg = regexp.MustCompile(`^[A-Za-z]{2,32}$`)

// CheckResource ...
func CheckResource(s string) error {
	if s == "" {
		return gear.ErrBadRequest.WithMsgf("empty resource type")
	}
	if !resourceReg.MatchString(s) {
		return gear.ErrBadRequest.WithMsgf("invalid resource type %s", strconv.Quote(s))
	}
	return nil
}

// CheckSubject ...
func CheckSubject(s string) error {
	if s == "" {
		return gear.ErrBadRequest.WithMsgf("empty subject")
	}
	_, err := otgo.ParseOTID("otid:ac:" + s)
	if err != nil {
		return gear.ErrBadRequest.From(err)
	}
	return nil
}

var permissionReg = regexp.MustCompile(`^[0-9A-Za-z]{2,32}$`)

// CheckPermission ...
func CheckPermission(s string) error {
	if s == "" {
		return gear.ErrBadRequest.WithMsgf("empty permission")
	}
	ss := strings.Split(s, ".")
	if err := CheckResource(ss[0]); err != nil {
		return err
	}
	for i := 1; i < len(ss); i++ {
		if !permissionReg.MatchString(ss[i]) {
			return gear.ErrBadRequest.WithMsgf("invalid permission %s", strconv.Quote(s))
		}
	}

	return nil
}

var whitespaceReg = regexp.MustCompile(`\s`)

// CheckTerm ...
func CheckTerm(s string) error {
	if len(s) < 3 {
		return gear.ErrBadRequest.WithMsgf("term' length too small %d", len(s))
	}
	if len(s) > 1024 {
		return gear.ErrBadRequest.WithMsgf("term' length too large %d", len(s))
	}
	if whitespaceReg.MatchString(s) {
		return gear.ErrBadRequest.WithMsgf("invalid term %s", strconv.Quote(s))
	}
	return nil
}

// TargetBatchAddInput ...
type TargetBatchAddInput struct {
	Targets []Target `json:"targets"`
	Parent  *Target  `json:"parent"`
	Scope   *Target  `json:"scope"`
}

// Validate 实现 gear.BodyTemplate
func (t *TargetBatchAddInput) Validate() error {
	if len(t.Targets) == 0 {
		return gear.ErrBadRequest.WithMsgf("targets empty")
	}
	if len(t.Targets) > 1000 {
		return gear.ErrBadRequest.WithMsgf("too many targets: %d", len(t.Targets))
	}
	cr := make(checkRepetitive)
	for _, target := range t.Targets {
		if err := cr.Check(target.Type + target.ID); err != nil {
			return err
		}
		if err := target.Validate(); err != nil {
			return err
		}
	}
	if t.Parent != nil {
		if err := t.Parent.Validate(); err != nil {
			return err
		}
	}
	if t.Scope != nil {
		if err := t.Scope.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// SubjectsInput ...
type SubjectsInput struct {
	Subjects []string `json:"subjects"`
}

// Validate 实现 gear.BodyTemplate
func (t *SubjectsInput) Validate() error {
	if len(t.Subjects) == 0 {
		return gear.ErrBadRequest.WithMsg("empty subjects")
	}
	if len(t.Subjects) > 1000 {
		return gear.ErrBadRequest.WithMsgf("too many subjects: %d", len(t.Subjects))
	}
	cr := make(checkRepetitive)
	for _, v := range t.Subjects {
		if err := cr.Check(v); err != nil {
			return err
		}
		if err := CheckSubject(v); err != nil {
			return err
		}
	}
	return nil
}

type checkRepetitive map[string]interface{}

func (c checkRepetitive) Check(s string) error {
	if _, ok := c[s]; ok {
		return gear.ErrBadRequest.WithMsgf("%s is repeated", s)
	}
	c[s] = struct{}{}
	return nil
}
