package i18nnames

import (
	"ytc/i18n"
	"ytc/internal/modules/ytc/collect/commons/datadef"
)

// GetBaseInfoItemName returns the localized name for a baseinfo item
func GetBaseInfoItemName(itemName string) string {
	keyMap := map[string]string{
		datadef.BASE_YASDB_VERION:      "report.base_yasdb_version",
		datadef.BASE_YASDB_PARAMETER:   "report.base_yasdb_parameter",
		datadef.BASE_HOST_OS_INFO:      "report.base_host_os_info",
		datadef.BASE_HOST_FIREWALLD:    "report.base_host_firewalld",
		datadef.BASE_HOST_CPU:          "report.base_host_cpu",
		datadef.BASE_HOST_DISK:         "report.base_host_disk",
		datadef.BASE_HOST_NETWORK:      "report.base_host_network",
		datadef.BASE_HOST_MEMORY:       "report.base_host_memory",
		datadef.BASE_HOST_NETWORK_IO:   "report.base_host_network_io",
		datadef.BASE_HOST_CPU_USAGE:    "report.base_host_cpu_usage",
		datadef.BASE_HOST_DISK_IO:      "report.base_host_disk_io",
		datadef.BASE_HOST_MEMORY_USAGE: "report.base_host_memory_usage",
	}
	if key, ok := keyMap[itemName]; ok {
		return i18n.T(key)
	}
	return itemName
}

// GetBaseInfoChildItemName returns the localized name for a baseinfo child item
func GetBaseInfoChildItemName(childName string) string {
	keyMap := map[string]string{
		"yasdb.ini":   "report.child_yasdb_ini",
		"v$parameter": "report.child_yasdb_parameter",
		"history":     "report.child_history",
		"current":     "report.child_current",
	}
	if key, ok := keyMap[childName]; ok {
		return i18n.T(key)
	}
	return childName
}

// GetDiagItemName returns the localized name for a diagnosis item
func GetDiagItemName(itemName string) string {
	keyMap := map[string]string{
		datadef.DIAG_YASDB_ADR:             "report.diag_yasdb_adr",
		datadef.DIAG_YASDB_RUNLOG:          "report.diag_yasdb_runlog",
		datadef.DIAG_YASDB_ALERTLOG:        "report.diag_yasdb_alertlog",
		datadef.DIAG_YASDB_PROCESS_STATUS:  "report.diag_yasdb_process_status",
		datadef.DIAG_YASDB_INSTANCE_STATUS: "report.diag_yasdb_instance_status",
		datadef.DIAG_YASDB_DATABASE_STATUS: "report.diag_yasdb_database_status",
		datadef.DIAG_HOST_SYSTEMLOG:        "report.diag_host_systemlog",
		datadef.DIAG_HOST_KERNELLOG:        "report.diag_host_kernellog",
		datadef.DIAG_YASDB_COREDUMP:        "report.diag_yasdb_coredump",
		datadef.DIAG_HOST_BASH_HISTORY:     "report.diag_host_bash_history",
	}
	if key, ok := keyMap[itemName]; ok {
		return i18n.T(key)
	}
	return itemName
}

// GetPerfItemName returns the localized name for a performance item
func GetPerfItemName(itemName string) string {
	keyMap := map[string]string{
		datadef.PERF_YASDB_AWR:      "report.perf_yasdb_awr",
		datadef.PERF_YASDB_SLOW_SQL: "report.perf_yasdb_slowsql",
	}
	if key, ok := keyMap[itemName]; ok {
		return i18n.T(key)
	}
	return itemName
}

// GetPerfChildItemName returns the localized name for a performance child item
func GetPerfChildItemName(childName string) string {
	keyMap := map[string]string{
		"slowParameter":     "report.perf_slowsql_parameter",
		"slowLogsInFile":    "report.perf_slowsql_table",
		"slowCutFile":       "report.perf_slowsql_file",
	}
	if key, ok := keyMap[childName]; ok {
		return i18n.T(key)
	}
	return childName
}
