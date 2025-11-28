package ytccollectcommons

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ytc/defs/errdef"
	"ytc/i18n"
	"ytc/log"
	"ytc/utils/stringutil"
	"ytc/utils/userutil"
	"ytc/utils/yasqlutil"

	"git.yasdb.com/go/yaslog"
	"git.yasdb.com/go/yasutil/fs"
)

// Deprecated constants - kept for backward compatibility
// Use i18n functions instead

// yasdb home
const (
	BIN   = "bin"
	YASQL = "yasql"
	YASDB = "yasdb"
)

// yasdb path
const (
	RUN    = "run"
	LOG    = "log"
	CONFIG = "config"
	SLOW   = "slow"
	DIAG   = "diag"
	ALERT  = "alert"
)

// yasdb node file
const (
	RUN_LOG   = "run.log"
	ALERT_LOG = "alert.log"
	SLOW_LOG  = "slow.log"
	YASDB_INI = "yasdb.ini"
)

const (
	HOST_DIR_NAME  = "host"
	YASDB_DIR_NAME = "yasdb"
)

type NoAccessRes struct {
	ModuleItem   string
	Description  string
	Tips         string
	ForceCollect bool // default false
}

type TimeRange struct {
	Start time.Time
	End   time.Time
}

func PathErrDescAndTips(path string, e error) (desc, tips string) {
	if os.IsNotExist(e) {
		desc = i18n.TWithData("common.file_not_exist_desc", map[string]interface{}{"Path": path})
		tips = i18n.TWithData("common.file_not_exist_tips", map[string]interface{}{"Path": path})
		return
	}
	if os.IsPermission(e) {
		desc = i18n.TWithData("common.file_permission_denied_desc", map[string]interface{}{
			"User": userutil.CurrentUser,
			"Path": path,
		})
		if err := userutil.CheckSudovn(log.Module); err != nil {
			if err == userutil.ErrSudoNeedPwd {
				tips = i18n.T("common.run_with_sudo_tips")
				return
			}
			tips = i18n.T("common.run_with_root_tips")
			return
		}
		tips = i18n.T("common.run_with_sudo_tips")
		return
	}
	desc = e.Error()
	tips = " "
	return
}

func FillDescTips(no *NoAccessRes, desc, tips string) {
	if no == nil {
		return
	}
	no.Description = desc
	no.Tips = tips
}

func CheckSudoTips(err error) string {
	if err == nil {
		return ""
	}
	if err == userutil.ErrSudoNeedPwd {
		return i18n.T("common.run_with_sudo_tips")
	}
	return i18n.T("common.run_with_root_tips")
}

func YasErrDescAndTips(err error) (desc string, tips string) {
	if err == nil {
		return
	}
	desc = i18n.TWithData("common.login_error", map[string]interface{}{"Error": err.Error()})
	switch e := err.(type) {
	case *yasqlutil.YasErr:
		switch e.Prefix {
		case yasqlutil.YAS_DB_NOT_OPEN:
			tips = i18n.T("common.database_not_open_tips")
		case yasqlutil.YAS_NO_DBUSER, yasqlutil.YAS_INVALID_USER_OR_PASSWORD:
			tips = i18n.T("common.invalid_user_password_tips")
		case yasqlutil.YAS_USER_LACK_LOGIN_AUTH:
			tips = i18n.T("common.lack_login_permission_tips")
		case yasqlutil.YAS_USER_LACK_AUTH:
			tips = i18n.T("common.lack_necessary_permission_tips")
		case yasqlutil.YAS_TABLE_OR_VIEW_DOES_NOT_EXIST:
			tips = i18n.T("common.object_not_exist_tips")
		case yasqlutil.YAS_FAILED_CONNECT_SOCKET:
			tips = i18n.T("common.connect_failed_tips")
		default:
			tips = i18n.T("common.yas_other_tips")
		}
	case *errdef.ItemEmpty:
		tips = i18n.T("common.item_empty_tips")
	default:
		tips = " "
	}
	return
}

func NotAccessItemToMap(noAccess []NoAccessRes) (res map[string]struct{}) {
	res = make(map[string]struct{})
	for _, noAccessRes := range noAccess {
		if noAccessRes.ForceCollect {
			continue
		}
		res[noAccessRes.ModuleItem] = struct{}{}
	}
	return
}

func FilesErrDescAndTips(res map[string]error) (desc string, tips string) {
	i := 1
	var files []string
	var buf strings.Builder
	for path, err := range res {
		buf.WriteString(fmt.Sprintf("%d.%s:%s; ", i, path, err.Error()))
		files = append(files, path)
		i++
	}
	desc = buf.String()
	tips = fmt.Sprintf("these files [%s] will not be collected", strings.Join(files, stringutil.STR_COMMA))
	return
}

func CopyDir(log yaslog.YasLog, src, dest string, excludeMap map[string]struct{}) (err error) {
	if strings.TrimSpace(src) == strings.TrimSpace(dest) {
		log.Infof("src path: %s is equal to dest path: %s, skip", src, dest)
		return
	}
	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if _, ok := excludeMap[path]; ok {
			log.Infof("skip exclude path: %s", path)
			return nil
		}
		if err != nil {
			log.Errorf("failed to copy dir, err: %s", err.Error())
			return nil
		}
		destNewPath := strings.Replace(path, src, dest, -1)
		if info.IsDir() {
			if err = os.MkdirAll(destNewPath, info.Mode()); err != nil {
				log.Errorf("failed to mkdir: %s, err: %s", destNewPath, err.Error())
				return nil
			}
		} else {
			if err = fs.CopyFile(path, destNewPath); err != nil {
				log.Infof("skip path: %s, because of err: %s", path, err.Error())
				return nil
			}
		}
		return nil
	})
	return
}
