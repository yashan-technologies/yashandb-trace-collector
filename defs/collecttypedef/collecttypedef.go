package collecttypedef

import (
	"errors"
	"fmt"
	"path"
	"time"

	"ytc/defs/timedef"
	"ytc/i18n"
)

const (
	TYPE_BASE  = "base"
	TYPE_DIAG  = "diag"
	TYPE_PERF  = "perf"
	TYPE_EXTRA = "extra"
)

const (
	WT_CPU     WorkloadType = "cpu"
	WT_NETWORK WorkloadType = "network"
	WT_MEMORY  WorkloadType = "memory"
	WT_DISK    WorkloadType = "disk"
)

const PACKAGE_NAME_PREFIX = "ytc"

var (
	ErrKnownType = errors.New("unknow collect type")
)

var (
	typeFullName = map[string]string{
		TYPE_BASE: "baseinfo",
		TYPE_DIAG: "diagnosis",
		TYPE_PERF: "performance",
	}
)

var (
	CollectTypeChineseName = map[string]string{
		TYPE_BASE:  "基础信息",
		TYPE_DIAG:  "诊断信息",
		TYPE_PERF:  "性能调优信息",
		TYPE_EXTRA: "额外收集项",
	}
)

type CollectParam struct {
	StartTime       time.Time `json:"startTime"`
	EndTime         time.Time `json:"endTime"`
	Output          string    `json:"output"`
	YasdbHome       string    `json:"yasdbHome"`
	YasdbData       string    `json:"yasdbData"`
	YasdbUser       string    `json:"yasdbUser"`
	YasdbPassword   string    `json:"yasdbPassword"`
	Include         []string  `json:"include"`
	Exclude         []string  `json:"exclude"`
	Lang            string    `json:"lang"`
	BeginTime       time.Time `json:"-"`
	YasdbHomeOSUser string    `json:"-"`
}

type WorkloadItem map[string]interface{}

type WorkloadOutput map[int64]WorkloadItem

type WorkloadType string

func GetTypeFullName(s string) string {
	// Use i18n for type names
	keyMap := map[string]string{
		TYPE_BASE:  "collect.type_baseinfo",
		TYPE_DIAG:  "collect.type_diagnosis",
		TYPE_PERF:  "collect.type_performance",
		TYPE_EXTRA: "collect.type_extra",
	}
	if key, ok := keyMap[s]; ok {
		return i18n.T(key)
	}
	// Fallback to old behavior
	full, ok := typeFullName[s]
	if !ok {
		full = s
	}
	return full
}

func (c *CollectParam) GetPackageTimestamp() string {
	// use begin time as timestamp
	return c.BeginTime.Format(timedef.TIME_FORMAT_IN_FILE)
}

func (c *CollectParam) GetPackageName() string {
	return fmt.Sprintf("%s-%s", PACKAGE_NAME_PREFIX, c.GetPackageTimestamp())
}

func (c *CollectParam) GenPackageRelativePath(p string) string {
	return path.Join(c.GetPackageName(), p)
}

// GetModuleName returns the localized module name
func GetModuleName(moduleType string) string {
	keyMap := map[string]string{
		TYPE_BASE:  "report.module_baseinfo",
		TYPE_DIAG:  "report.module_diagnosis",
		TYPE_PERF:  "report.module_performance",
		TYPE_EXTRA: "report.module_extra",
	}
	if key, ok := keyMap[moduleType]; ok {
		return i18n.T(key)
	}
	return moduleType
}
