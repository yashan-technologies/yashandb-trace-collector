package ytcctlhandler

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"ytc/commons/std"
	"ytc/defs/bashdef"
	"ytc/defs/collecttypedef"
	"ytc/defs/errdef"
	"ytc/i18n"
	ytccollect "ytc/internal/modules/ytc/collect"
	ytccollectcommons "ytc/internal/modules/ytc/collect/commons"
	"ytc/internal/modules/ytc/collect/commons/i18nnames"
	"ytc/log"
	"ytc/utils/stringutil"
	"ytc/utils/terminalutil/barutil"

	"git.yasdb.com/go/yasutil/tabler"
)

var (
	_module_order = []string{
		collecttypedef.TYPE_BASE,
		collecttypedef.TYPE_DIAG,
		collecttypedef.TYPE_PERF,
		collecttypedef.TYPE_EXTRA,
	}
)

func (c *CollecterHandler) Collect(yasdbValidate error) error {
	noAccessMap, err := c.checkAccess(yasdbValidate)
	if err != nil {
		log.Handler.Errorf(err.Error())
		return err
	}
	moduleItems, err := c.getCollectItem(noAccessMap)
	if err != nil {
		log.Handler.Errorf(err.Error())
		return err
	}
	if err := c.printCollectItem(moduleItems); err != nil {
		log.Handler.Errorf(err.Error())
		return err
	}
	fmt.Printf("\n%s\n\n", i18n.T("collect.starting"))
	return c.collect(moduleItems)
}

func (c *CollecterHandler) checkAccess(yasdbValidateErr error) (map[string][]ytccollectcommons.NoAccessRes, error) {
	m := make(map[string][]ytccollectcommons.NoAccessRes)
	for _, c := range c.Collecters {
		noAccessList := c.CheckAccess(yasdbValidateErr)
		if len(noAccessList) != 0 {
			m[c.Type()] = noAccessList
		}
	}
	if len(m) == 0 {
		return m, nil
	}
	if err := c.printNoAccessItem(m); err != nil {
		return m, err
	}
	return m, nil
}

func (c *CollecterHandler) printNoAccessItem(m map[string][]ytccollectcommons.NoAccessRes) error {
	table := tabler.NewTable(
		"",
		tabler.NewRowTitle(i18n.T("collect.header_type"), 15),
		tabler.NewRowTitle(i18n.T("collect.header_item"), 25),
		tabler.NewRowTitle(i18n.T("collect.header_desc"), 50),
		tabler.NewRowTitle(i18n.T("collect.header_tips"), 50),
		tabler.NewRowTitle(i18n.T("collect.header_collected"), 15),
	)
	fmt.Printf("%s\n\n", bashdef.WithYellow(i18n.T("collect.tips_for_you")))
	var modules []string
	for t := range m {
		modules = append(modules, t)
	}
	sort.Strings(modules)
	for _, module := range modules {
		for i, noAccess := range m[module] {
			// Translate module name based on type
			moduleName := noAccess.ModuleItem
			switch module {
			case collecttypedef.TYPE_BASE:
				moduleName = i18nnames.GetBaseInfoItemName(moduleName)
			case collecttypedef.TYPE_DIAG:
				moduleName = i18nnames.GetDiagItemName(moduleName)
			case collecttypedef.TYPE_PERF:
				moduleName = i18nnames.GetPerfItemName(moduleName)
			}
			
			if i == 0 {
				err := table.AddColumn(strings.ToUpper(collecttypedef.GetTypeFullName(module)), moduleName, noAccess.Description, noAccess.Tips, isCollectedStr(noAccess.ForceCollect))
				if err != nil {
					log.Handler.Errorf("add column err: %s", err.Error())
					return err
				}
				continue
			}
			if err := table.AddColumn("", moduleName, noAccess.Description, noAccess.Tips, isCollectedStr(noAccess.ForceCollect)); err != nil {
				log.Module.Errorf("add column err: %s", err.Error())
				return err
			}
		}
	}
	table.Print()
	var isConfirm string
	fmt.Printf(i18n.T("collect.continue_confirm"))
	fmt.Scanln(&isConfirm)

	// record input
	std.WriteToFile(isConfirm + stringutil.STR_NEWLINE)

	isConfirm = strings.ToLower(isConfirm)
	if isConfirm != "y" {
		return fmt.Errorf(i18n.T("collect.validation_failed"))
	}
	return nil
}

func (c *CollecterHandler) printCollectItem(typeItem map[string][]string) error {
	var (
		itemTitle   []*tabler.RowTitle
		moduleItems = make([][]string, 0)
		moduleNames = make([]string, 0)
	)
	for module := range typeItem {
		if len(typeItem[module]) == 0 {
			continue
		}
		moduleNames = append(moduleNames, module)
	}
	sort.Strings(moduleNames)
	if len(moduleNames) == 0 {
		return errdef.ErrNoneCollectTtem
	}
	for _, t := range moduleNames {
		itemTitle = append(itemTitle, tabler.NewRowTitle(strings.ToUpper(collecttypedef.GetTypeFullName(t)), 30))
		moduleItems = append(moduleItems, typeItem[t])
	}
	table := tabler.NewTable("", itemTitle...)
	maxCol := maxCol(moduleItems)
	for i := 0; i < maxCol; i++ {
		row := make([]interface{}, len(moduleNames))
		for j, item := range moduleItems {
			if i < len(item) {
				// Translate module name based on type
				moduleName := item[i]
				moduleType := moduleNames[j]
				switch moduleType {
				case collecttypedef.TYPE_BASE:
					moduleName = i18nnames.GetBaseInfoItemName(moduleName)
				case collecttypedef.TYPE_DIAG:
					moduleName = i18nnames.GetDiagItemName(moduleName)
				case collecttypedef.TYPE_PERF:
					moduleName = i18nnames.GetPerfItemName(moduleName)
				}
				row[j] = moduleName
				continue
			}
			row[j] = " "
		}
		if err := table.AddColumn(row...); err != nil {
			return err
		}
	}
	fmt.Printf("%s\n\n", bashdef.WithBlue(i18n.T("collect.following_modules")))
	table.Print()
	return nil
}

func (c *CollecterHandler) getCollectItem(noAccessMap map[string][]ytccollectcommons.NoAccessRes) (typeItem map[string][]string, err error) {
	typeItem = make(map[string][]string)
	for _, collect := range c.Collecters {
		t := collect.Type()
		noAccess, ok := noAccessMap[t]
		if !ok {
			noAccess = make([]ytccollectcommons.NoAccessRes, 0)
		}
		typeItem[t] = collect.ItemsToCollect(noAccess)
	}
	return
}

func (c *CollecterHandler) collect(moduleItems map[string][]string) error {
	progress := barutil.NewProgress(barutil.WithWidth(100))
	if e := c.PreCollect(); e != nil {
		return e
	}
	collMap := c.collecterMap()
	moduleFuncs := make(map[string]map[string]func() error)
	for module, items := range moduleItems {
		_, ok := collMap[module]
		if !ok {
			log.Handler.Errorf("collect type: %s not exist", module)
			continue
		}
		moduleFuncs[module] = collMap[module].CollectFunc(items)
	}
	for _, module := range _module_order {
		if funcs, ok := moduleFuncs[module]; ok {
			progress.AddBar(module, funcs)
		}
	}
	progress.Start()
	return c.CollectOK()
}

func (c *CollecterHandler) PreCollect() error {
	c.CollectResult.CollectBeginTime = c.CollectResult.CollectParam.BeginTime
	packageDir := c.CollectResult.GetPackageDir()
	for _, collecter := range c.Collecters {
		if err := collecter.PreCollect(packageDir); err != nil {
			return err
		}
	}
	return nil
}

func (c *CollecterHandler) collecterMap() (res map[string]ytccollect.TypedCollecter) {
	res = make(map[string]ytccollect.TypedCollecter)
	for _, c := range c.Collecters {
		res[c.Type()] = c
	}
	return
}

func (c *CollecterHandler) CollectOK() error {
	c.CollectResult.CollectEndTime = time.Now()
	for _, collecter := range c.Collecters {
		c.CollectResult.Modules[collecter.Type()] = collecter.CollectOK()
	}
	fmt.Printf("%s\n\n", i18n.T("collect.packing_results"))
	path, err := c.CollectResult.GenResult(c.CollectResult.CollectParam.Output, c.Types)
	if err != nil {
		err = fmt.Errorf(i18n.TWithData("collect.gen_result_failed", map[string]interface{}{"Err": err}))
		log.Handler.Error(err)
		fmt.Println(err.Error())
		return err
	}
	status := bashdef.WithGreen(i18n.T("collect.completed"))
	pathStr := bashdef.WithBlue(path)
	fmt.Println(i18n.TWithData("collect.completed_msg", map[string]interface{}{"Status": status, "Path": pathStr}))
	return nil
}

func isCollectedStr(f bool) string {
	flag := strconv.FormatBool(f)
	if f {
		return bashdef.WithGreen(flag)
	}
	return bashdef.WithRed(flag)
}

func maxCol(rows [][]string) int {
	max := -1
	for _, row := range rows {
		if len(row) > max {
			max = len(row)
		}
	}
	return max
}
