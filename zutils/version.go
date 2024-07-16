package zutils

import (
	"fmt"
	"strings"
)

/*
 *  ┏┓      ┏┓
 *┏━┛┻━━━━━━┛┻┓
 *┃　　　━　　  ┃
 *┃   ┳┛ ┗┳   ┃
 *┃           ┃
 *┃     ┻     ┃
 *┗━━━┓     ┏━┛
 *　　 ┃　　　┃神兽保佑
 *　　 ┃　　　┃代码无BUG！
 *　　 ┃　　　┗━━━┓
 *　　 ┃         ┣┓
 *　　 ┃         ┏┛
 *　　 ┗━┓┓┏━━┳┓┏┛
 *　　   ┃┫┫  ┃┫┫
 *      ┗┻┛　 ┗┻┛
 @Time    : 2024/7/13 -- 15:49
 @Author  : bishop ❤️ MONEY
 @Description: version.go
*/

type VersionCmp struct {
	ver string
}

func NewVersionCmp(ver string) *VersionCmp {
	v := &VersionCmp{}

	v.ver = v.fmtver(ver)
	return v
}

func (m *VersionCmp) fmtver(ver string) string {
	pvs := strings.Split(ver, ".")

	rv := ""
	for _, pv := range pvs {
		rv += fmt.Sprintf("%020s", pv)
	}

	return rv

}

func (m *VersionCmp) Min() string {
	return m.fmtver("0")
}

func (m *VersionCmp) Max() string {
	return m.fmtver("99999999999999999999")
}

func (m *VersionCmp) Lt(ver string) bool {
	return m.ver < m.fmtver(ver)
}

func (m *VersionCmp) Lte(ver string) bool {
	return m.ver <= m.fmtver(ver)
}

func (m *VersionCmp) Gt(ver string) bool {
	return m.ver > m.fmtver(ver)
}

func (m *VersionCmp) Gte(ver string) bool {
	return m.ver >= m.fmtver(ver)
}

func (m *VersionCmp) Eq(ver string) bool {
	return m.ver == m.fmtver(ver)
}

func (m *VersionCmp) Ne(ver string) bool {
	return m.ver != m.fmtver(ver)
}

func (m *VersionCmp) GetFormatVersion() string {
	return m.ver
}
