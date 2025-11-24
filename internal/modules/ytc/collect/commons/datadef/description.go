package datadef

import (
	"strings"

	"ytc/i18n"
)

// Deprecated constants - kept for backward compatibility
// Use i18n functions instead

const (
	key_command_not_found = "command not found"
)

func GenDefaultDesc() string {
	return i18n.T("desc.default")
}

func GenNoPermissionDesc(str string) string {
	return i18n.TWithData("desc.no_permission", map[string]interface{}{"Path": str})
}

func GenHostWorkloadDesc(e error) string {
	if strings.Contains(e.Error(), key_command_not_found) {
		return i18n.T("desc.no_sar_command")
	}
	return i18n.T("desc.default")
}

func GenUbuntuFirewalldDesc() string {
	return i18n.T("desc.ubuntu_firewalld")
}

func GenKylinDmesgDesc() string {
	return i18n.T("desc.kylin_dmesg")
}

func GenGetCoreDumpPathDesc() string {
	return i18n.T("desc.get_coredump_path")
}

func GenReadCoreDumpPathDesc(path string) string {
	return i18n.TWithData("desc.read_coredump_path", map[string]interface{}{"Path": path})
}

func GenGetDatabaseParameterDesc(parameter string) string {
	return i18n.TWithData("desc.get_database_parameter", map[string]interface{}{"Parameter": parameter})
}

func GenSkipCollectDatabaseInfoDesc() string {
	return i18n.T("desc.skip_collect_database_info")
}

func GenYasdbProcessStatusDesc() string {
	return i18n.T("desc.yasdb_process_status")
}

func GenGetDatabaseViewDesc(view string) string {
	return i18n.TWithData("desc.get_database_view", map[string]interface{}{"View": view})
}

func GenNoPermissionSyslogDesc() string {
	return i18n.T("desc.no_permission_syslog")
}
