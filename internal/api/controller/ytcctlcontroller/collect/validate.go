package collect

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"ytc/defs/collecttypedef"
	"ytc/defs/confdef"
	"ytc/defs/errdef"
	"ytc/defs/regexdef"
	"ytc/defs/runtimedef"
	"ytc/log"
	"ytc/utils/fileutil"
	"ytc/utils/jsonutil"
	"ytc/utils/stringutil"
	"ytc/utils/timeutil"
	"ytc/utils/userutil"

	"git.yasdb.com/go/yasutil/fs"
)

const (
	f_type   = "type"
	f_range  = "range"
	f_start  = "start"
	f_end    = "end"
	f_output = "output"
)

var (
	_examples_time = []string{
		"yyyy-MM-dd",
		"yyyy-MM-dd-hh",
		"yyyy-MM-dd-hh-mm",
	}

	_exmaples_range = []string{
		"1M",
		"1d",
		"1h",
		"1m",
	}

	_range_help = "you must ensure that the number before (M|d|h|m) is greater than 0"

	_type_help = "you can choose one or more of (base|diag|perf), split with ',', such as 'base,diag,perf'."
)

func (c *CollectCmd) validate() error {
	if err := c.validateType(); err != nil {
		return err
	}
	if err := c.validateRange(); err != nil {
		return err
	}
	if err := c.validateStartAndEnd(); err != nil {
		return err
	}
	if err := c.validateOutput(); err != nil {
		return err
	}
	if err := c.validateIncludePath(); err != nil {
		return err
	}
	return nil
}

func (c *CollectCmd) validateType() error {
	c.Type = trimSpace(c.Type)
	resMap := make(map[string]struct{})
	tMap := map[string]struct{}{
		collecttypedef.TYPE_BASE: {},
		collecttypedef.TYPE_DIAG: {},
		collecttypedef.TYPE_PERF: {},
	}
	types := strings.Split(c.Type, stringutil.STR_COMMA)
	for _, t := range types {
		if _, ok := tMap[t]; !ok {
			return errdef.NewErrYtcFlag(f_type, c.Type, nil, _type_help)
		}
		resMap[t] = struct{}{}
	}
	return nil
}

func (c *CollectCmd) validateIncludePath() error {
	includes := c.getExtraPath(c.Include)
	for _, include := range includes {
		if c.inExclude(include) {
			continue
		}
		if _, err := os.Stat(include); err != nil {
			return err
		}
		if err := c.checkInclude(include); err != nil {
			return err
		}
	}
	return nil
}

func (c *CollectCmd) checkInclude(include string) error {
	invalidPaths := []string{
		runtimedef.GetYTCHome(),
		c.Output,
	}
	for _, invalidPath := range invalidPaths {
		if c.inExclude(invalidPath) {
			continue
		}
		if fileutil.IsAncestorDir(include, invalidPath) {
			return fmt.Errorf("could not collect: %s, because it contails the invalid path: %s", include, invalidPath)
		}
		if fileutil.IsAncestorDir(invalidPath, include) {
			return fmt.Errorf("could not collect: %s, because the invalid path: %s contails it", include, invalidPath)
		}
	}
	return nil
}

func (c *CollectCmd) inExclude(filePath string) bool {
	excludes := c.getExtraPath(c.Exclude)
	for _, exclude := range excludes {
		if fileutil.IsAncestorDir(exclude, filePath) {
			return true
		}
	}
	return false
}

func (c *CollectCmd) validateRange() error {
	strategyConf := confdef.GetStrategyConf()
	log.Controller.Debugf("strategy: %s\n", jsonutil.ToJSONString(strategyConf))
	log.Controller.Debugf("cmd: %s", jsonutil.ToJSONString(c))
	if stringutil.IsEmpty(c.Range) {
		return nil
	}
	if !regexdef.RangeRegex.MatchString(c.Range) {
		return errdef.NewErrYtcFlag(f_range, c.Range, _exmaples_range, _range_help)
	}
	minDuration, maxDuration, err := strategyConf.Collect.GetMinAndMaxDur()
	if err != nil {
		log.Controller.Errorf("get duration err: %s", err.Error())
		return err
	}
	log.Controller.Debugf("get min %s max %s", minDuration.String(), maxDuration.String())
	r, err := timeutil.GetDuration(c.Range)
	if err != nil {
		return err
	}
	if r > maxDuration {
		return errdef.NewGreaterMaxDur(strategyConf.Collect.MaxDuration)
	}
	if r < minDuration {
		return errdef.NewLessMinDur(strategyConf.Collect.MinDuration)
	}
	return nil
}

func (c *CollectCmd) validateStartAndEnd() error {
	strategyConf := confdef.GetStrategyConf()
	var (
		startNotEmpty, endNotEmpty bool
		start, end                 time.Time
		err                        error
	)
	if !stringutil.IsEmpty(c.Start) {
		if !regexdef.TimeRegex.MatchString(c.Start) {
			return errdef.NewErrYtcFlag(f_start, c.Start, _examples_time, "")
		}
		start, err = timeutil.GetTimeDivBySepa(c.Start, stringutil.STR_HYPHEN)
		if err != nil {
			return err
		}
		now := time.Now()
		if start.After(now) {
			return errdef.ErrStartShouldLessCurr
		}
		startNotEmpty = true
	}
	if !stringutil.IsEmpty(c.End) {
		if !regexdef.TimeRegex.MatchString(c.End) {
			return errdef.NewErrYtcFlag(f_end, c.End, _examples_time, "")
		}
		end, err = timeutil.GetTimeDivBySepa(c.End, stringutil.STR_HYPHEN)
		if err != nil {
			return err
		}
		endNotEmpty = true
	}
	if startNotEmpty && endNotEmpty {
		minDuration, maxDuration, err := strategyConf.Collect.GetMinAndMaxDur()
		if err != nil {
			log.Controller.Errorf("get duration err: %s", err.Error())
			return err
		}
		if end.Before(start) {
			return errdef.ErrEndLessStart
		}
		r := end.Sub(start)
		if r > maxDuration {
			return errdef.NewGreaterMaxDur(strategyConf.Collect.MaxDuration)
		}
		if r < minDuration {
			return errdef.NewLessMinDur(strategyConf.Collect.MaxDuration)
		}
	}
	return nil
}

func (c *CollectCmd) validateOutput() error {
	output := c.Output
	if !regexdef.PathRegex.Match([]byte(output)) {
		return errdef.ErrPathFormat
	}
	if !path.IsAbs(output) {
		output = path.Join(runtimedef.GetYTCHome(), output)
	}
	_, err := os.Stat(output)
	if err != nil {
		if os.IsPermission(err) {
			return errdef.NewErrPermissionDenied(userutil.CurrentUser, output)
		}
		if !os.IsNotExist(err) {
			return err
		}
		if err := fs.Mkdir(output); err != nil {
			log.Controller.Errorf("create output err: %s", err.Error())
			if os.IsPermission(err) {
				return errdef.NewErrPermissionDenied(userutil.CurrentUser, output)
			}
			return err
		}
	}
	return fileutil.CheckUserWrite(output)

}

func (c *CollectCmd) fillDefault() {
	if stringutil.IsEmpty(c.Output) {
		c.Output = confdef.GetStrategyConf().Collect.Output
	}
	if !path.IsAbs(c.Output) {
		c.Output = path.Join(runtimedef.GetYTCHome(), c.Output)
	}
	c.Output = path.Clean(c.Output)
}
