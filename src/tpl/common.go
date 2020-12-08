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
	TotalSize     int         `json:"totalSize,omitempty"`
	NextPageToken string      `json:"nextPageToken,omitempty"`
	Result        interface{} `json:"result"`
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

	if pg.PageSize > 10000 {
		return gear.ErrBadRequest.WithMsgf("pageSize %v should not great than 10000", pg.PageSize)
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
	Type string `json:"type"`
	ID   string `json:"id"`
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
	if len(s) < 64 {
		return gear.ErrBadRequest.WithMsgf("term' length too large %d", len(s))
	}
	if whitespaceReg.MatchString(s) {
		return gear.ErrBadRequest.WithMsgf("invalid term %s", strconv.Quote(s))
	}
	return nil
}
