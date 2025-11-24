package performance

import (
	"path"
	"strings"

	"ytc/defs/confdef"
	"ytc/i18n"
	ytccollectcommons "ytc/internal/modules/ytc/collect/commons"
	"ytc/internal/modules/ytc/collect/commons/datadef"
	"ytc/internal/modules/ytc/collect/yasdb"
	"ytc/log"
	"ytc/utils/fileutil"
	"ytc/utils/yasqlutil"

	"git.yasdb.com/go/yaslog"
)

const (
	USER_SYS = "SYS"
)

func (p *PerfCollecter) checkDatabaseOpenMode(logger yaslog.YasLog) (bool, error) {
	tx := yasqlutil.GetLocalInstance(p.YasdbUser, p.YasdbPassword, p.YasdbHome, p.YasdbData)
	database, err := yasdb.QueryDatabase(tx)
	if err != nil {
		logger.Errorf("query v$database failed: %s", err)
		return false, err
	}
	return database.IsDatabaseInReadWwiteMode(), nil
}

func (p *PerfCollecter) checkAWR() *ytccollectcommons.NoAccessRes {
	noAccess := &ytccollectcommons.NoAccessRes{ModuleItem: datadef.PERF_YASDB_AWR}
	if strings.ToUpper(p.YasdbUser) != USER_SYS {
		ytccollectcommons.FillDescTips(noAccess, i18n.TWithData("perf.user_not_sys_desc", map[string]interface{}{"User": p.YasdbUser}), i18n.T("perf.user_not_sys_tips"))
		return noAccess
	}
	if p.yasdbValidateErr != nil {
		desc, tips := ytccollectcommons.YasErrDescAndTips(p.yasdbValidateErr)
		ytccollectcommons.FillDescTips(noAccess, desc, tips)
		return noAccess
	}

	logger := log.Module.M(datadef.PERF_YASDB_AWR)
	inReadWriteMode, err := p.checkDatabaseOpenMode(logger)
	if err != nil {
		desc, tips := ytccollectcommons.YasErrDescAndTips(err)
		ytccollectcommons.FillDescTips(noAccess, desc, tips)
		return noAccess
	}
	if !inReadWriteMode {
		desc, tips := i18n.T("perf.awr_skip_tips"), i18n.T("perf.awr_skip_tips")
		ytccollectcommons.FillDescTips(noAccess, desc, tips)
		return noAccess
	}
	// TODO: check v$instance boot time
	if _, _, err := p.genStartEndSnapId(logger); err != nil {
		if err == yasdb.ErrNoSatisfiedSnapshot {
			desc := i18n.T("perf.no_satisfied_snap_desc")
			tips := i18n.T("perf.no_satisfied_tips")
			ytccollectcommons.FillDescTips(noAccess, desc, tips)
			return noAccess
		}
		desc, tips := ytccollectcommons.YasErrDescAndTips(err)
		ytccollectcommons.FillDescTips(noAccess, desc, tips)
		return noAccess
	}
	collectConfig := confdef.GetStrategyConf().Collect
	timeout := collectConfig.GetAWRTimeout()
	desc := i18n.T("perf.awr_timeout_desc")
	tips := i18n.TWithData("perf.awr_timeout_tips", map[string]interface{}{"Timeout": timeout.String()})
	ytccollectcommons.FillDescTips(noAccess, desc, tips)
	noAccess.ForceCollect = true
	return noAccess
}

func (p *PerfCollecter) checkSlowSql() *ytccollectcommons.NoAccessRes {
	noAccess := &ytccollectcommons.NoAccessRes{ModuleItem: datadef.PERF_YASDB_SLOW_SQL}
	defaultSlowLog := path.Join(p.YasdbData, ytccollectcommons.LOG, ytccollectcommons.SLOW, ytccollectcommons.SLOW_LOG)
	defaultSlowLogTips := i18n.TWithData("perf.default_slowsql_tips", map[string]interface{}{"Path": defaultSlowLog})
	if p.yasdbValidateErr != nil {
		desc, tips := ytccollectcommons.YasErrDescAndTips(p.yasdbValidateErr)
		if err := fileutil.CheckAccess(defaultSlowLog); err != nil {
			ytccollectcommons.FillDescTips(noAccess, desc, tips)
			return noAccess
		}
		ytccollectcommons.FillDescTips(noAccess, desc, defaultSlowLogTips)
		noAccess.ForceCollect = true
		return noAccess
	}
	slowLogPath, err := p.getSlowLogPath()
	slowLog := path.Join(slowLogPath, ytccollectcommons.SLOW_LOG)
	if err != nil {
		desc, tips := ytccollectcommons.YasErrDescAndTips(err)
		if err := fileutil.CheckAccess(defaultSlowLog); err != nil {
			ytccollectcommons.FillDescTips(noAccess, desc, tips)
			return noAccess
		}
		ytccollectcommons.FillDescTips(noAccess, desc, tips)
		noAccess.ForceCollect = true
		return noAccess
	}
	if err := fileutil.CheckAccess(slowLog); err != nil {
		desc, tips := ytccollectcommons.PathErrDescAndTips(slowLog, err)
		ytccollectcommons.FillDescTips(noAccess, desc, tips)
		return noAccess
	}
	return nil
}
