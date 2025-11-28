package data

import (
	"fmt"
	"strings"
	"time"

	"ytc/defs/collecttypedef"
	"ytc/defs/timedef"
	"ytc/i18n"
	"ytc/internal/modules/ytc/collect/commons/datadef"
	"ytc/internal/modules/ytc/collect/commons/i18nnames"
	report "ytc/internal/modules/ytc/collect/data/reporter"
	"ytc/internal/modules/ytc/collect/data/reporter/commons"
	"ytc/internal/modules/ytc/collect/extra"
	"ytc/internal/modules/ytc/collect/resultgenner"
	"ytc/internal/modules/ytc/collect/resultgenner/reporter"
	"ytc/internal/modules/ytc/collect/resultgenner/reporter/htmldef"
	"ytc/log"
	"ytc/utils/stringutil"

	"git.yasdb.com/go/yaserr"
	"github.com/jedib0t/go-pretty/v6/table"
)

// validate interface
var _ resultgenner.Genner = (*YTCReport)(nil)

type YTCReport struct {
	CollectBeginTime time.Time                     `json:"collectBeginTime"`
	CollectEndTime   time.Time                     `json:"collectEndTime"`
	CollectParam     *collecttypedef.CollectParam  `json:"collectParam"`
	Modules          map[string]*datadef.YTCModule `json:"modules"`
	genner           resultgenner.BaseGenner
}

func NewYTCReport(param *collecttypedef.CollectParam) *YTCReport {
	return &YTCReport{
		CollectParam: param,
		Modules:      make(map[string]*datadef.YTCModule),
		genner:       resultgenner.BaseGenner{},
	}
}

// [Interface Func]
func (r *YTCReport) GenData(data interface{}, fname string) error {
	return r.genner.GenData(data, fname)
}

// [Interface Func]
func (r *YTCReport) GenReport() (content reporter.ReportContent, err error) {
	var graphs []string
	logger := log.Module.M("generate report")
	moduleNum := 0
	for _, moduleName := range _moduleOrder {
		module, ok := r.Modules[moduleName]
		if !ok {
			logger.Infof("module: %s unfound, pass", moduleName)
			continue
		}
		moduleNum++

		moduleTitlePrefix := fmt.Sprintf("%d", moduleNum)
		moduleContent := reporter.GenReportContentByTitle(fmt.Sprintf("%s %s", moduleTitlePrefix, collecttypedef.GetModuleName(moduleName)), reporter.FONT_SIZE_H1)
		content.Txt += moduleContent.Txt
		content.Markdown += moduleContent.Markdown
		content.HTML += moduleContent.HTML

		itemNum := 0
		items := module.Items()
		for _, itemName := range _itemOrder[moduleName] {
			item, ok := items[itemName]
			if !ok {
				logger.Infof("item: %s unfound, pass", itemName)
				continue
			}
			reporter, ok := report.REPORTERS[itemName]
			if !ok {
				err = fmt.Errorf("reporter of %s unfound", itemName)
				err = yaserr.Wrapf(err, "get reporter")
				return
			}
			itemNum++

			itemTitlePrefix := moduleTitlePrefix + stringutil.STR_DOT + fmt.Sprintf("%d", itemNum)
			itemContent, e := reporter.Report(*item, itemTitlePrefix)
			if e != nil {
				err = yaserr.Wrapf(e, "generete report of %s", itemName)
				return
			}
			if !stringutil.IsEmpty(itemContent.Graph) {
				graphs = append(graphs, itemContent.Graph)
			}
			content.Txt += itemContent.Txt + stringutil.STR_NEWLINE
			content.Markdown += itemContent.Markdown + stringutil.STR_NEWLINE
			content.HTML += itemContent.HTML + stringutil.STR_NEWLINE
		}
	}
	content = r.addSummary(content, r.genReportOverview(), r.genReportItems())
	content.HTML = htmldef.GenHTML(content.HTML, strings.Join(graphs, stringutil.STR_NEWLINE))
	return
}

func (r *YTCReport) GenResult(outputDir string, types map[string]struct{}) (string, error) {
	for _, m := range r.Modules {
		m.FillJSONItems()
	}
	genner := resultgenner.BaseResultGenner{
		Datas:        r.Modules,
		CollectTypes: types,
		OutputDir:    outputDir,
		Timestamp:    r.CollectBeginTime.Format(timedef.TIME_FORMAT_IN_FILE),
		PackageName:  r.CollectParam.GetPackageName(),
		Genner:       r,
	}
	return genner.GenResult()
}

func (r *YTCReport) GetPackageDir() string {
	genner := resultgenner.BaseResultGenner{
		OutputDir:   r.CollectParam.Output,
		Timestamp:   r.CollectBeginTime.Format(timedef.TIME_FORMAT_IN_FILE),
		PackageName: r.CollectParam.GetPackageName(),
	}
	return genner.GetPackageDir()
}

func (r *YTCReport) genReportOverview() (content reporter.ReportContent) {
	titleContent := reporter.GenReportContentByTitle(i18n.T("report.overview_title"), reporter.FONT_SIZE_H1)
	genTableRows := func(sep string) []table.Row {
		user := r.CollectParam.YasdbUser
		if stringutil.IsEmpty(user) {
			user = reporter.PLACEHOLDER
		}
		var modules []string
		for _, m := range _moduleOrder {
			if _, ok := r.Modules[m]; ok {
				modules = append(modules, collecttypedef.GetModuleName(m))
			}
		}
		separator := "，"
		if r.CollectParam.Lang == "en" {
			separator = ", "
		}
		rows := []table.Row{
			{i18n.T("report.collect_type"), strings.Join(modules, separator)},
			{i18n.T("report.collect_range_start"), r.CollectParam.StartTime.Format(timedef.TIME_FORMAT_UNTIL_MINITE)},
			{i18n.T("report.collect_range_end"), r.CollectParam.EndTime.Format(timedef.TIME_FORMAT_UNTIL_MINITE)},
			{i18n.T("report.yasdb_home"), r.CollectParam.YasdbHome},
			{i18n.T("report.yasdb_data"), r.CollectParam.YasdbData},
			{i18n.T("report.database_user"), user},
		}
		if len(r.CollectParam.Include) > 0 {
			rows = append(rows, table.Row{i18n.T("report.extra_files"), strings.Join(r.CollectParam.Include, sep)})
		}
		if len(r.CollectParam.Exclude) > 0 {
			rows = append(rows, table.Row{i18n.T("report.filtered_files"), strings.Join(r.CollectParam.Exclude, sep)})
		}
		rows = append(rows, table.Row{i18n.T("report.output_dir"), r.CollectParam.Output})
		rows = append(rows, table.Row{i18n.T("report.task_start_time"), r.CollectBeginTime.Format(timedef.TIME_FORMAT)})
		rows = append(rows, table.Row{i18n.T("report.task_end_time"), r.CollectEndTime.Format(timedef.TIME_FORMAT)})
		return rows
	}

	tw := commons.ReporterWriter.NewTableWriter()
	tw.AppendHeader(table.Row{i18n.T("report.overview_item"), i18n.T("report.overview_value")})
	baseInfoTitle := reporter.GenReportContentByTitle(i18n.T("report.basic_overview"), reporter.FONT_SIZE_H2)

	// render txt
	for _, r := range genTableRows(stringutil.STR_NEWLINE) {
		tw.AppendRow(r)
		tw.AppendSeparator()
	}
	content.Txt = strings.Join([]string{titleContent.Txt, baseInfoTitle.Txt, tw.Render()}, stringutil.STR_NEWLINE)
	tw.ResetRows()

	// render markdown and html
	tw.AppendRows(genTableRows(stringutil.STR_HTML_BR))
	content.Markdown = strings.Join([]string{titleContent.Markdown, baseInfoTitle.Markdown, tw.RenderMarkdown()}, stringutil.STR_NEWLINE)
	content.HTML = strings.Join([]string{titleContent.HTML, baseInfoTitle.HTML, tw.RenderHTML()}, stringutil.STR_NEWLINE)
	return
}

func (r *YTCReport) genModulesAndItems() (modules []string, items [][]string) {
	for _, m := range _moduleOrder {
		if module, ok := r.Modules[m]; ok {
			modules = append(modules, collecttypedef.GetModuleName(m))
			tmpItems := module.Items()
			switch module.Module {
			case collecttypedef.TYPE_BASE:
				var names []string
				for _, item := range _baseItemOrder {
					if _, ok := tmpItems[item]; ok {
						names = append(names, i18nnames.GetBaseInfoItemName(item))
					}
				}
				items = append(items, names)
			case collecttypedef.TYPE_DIAG:
				var names []string
				for _, item := range _diagItemOrder {
					if _, ok := tmpItems[item]; ok {
						names = append(names, i18nnames.GetDiagItemName(item))
					}
				}
				items = append(items, names)
			case collecttypedef.TYPE_PERF:
				var names []string
				for _, item := range _perfItemOrder {
					if _, ok := tmpItems[item]; ok {
						names = append(names, i18nnames.GetPerfItemName(item))
					}
				}
				items = append(items, names)
			case collecttypedef.TYPE_EXTRA:
				var names []string
				for _, item := range _extraItemOrder {
					if _, ok := tmpItems[item]; ok {
						names = append(names, extra.ExtraChineseName[item])
					}
				}
				items = append(items, names)
			default:
			}
		}
	}
	return

}

func (r *YTCReport) genReportItems() (content reporter.ReportContent) {
	titleContent := reporter.GenReportContentByTitle(i18n.T("report.items_overview"), reporter.FONT_SIZE_H2)
	modules, items := r.genModulesAndItems()
	var tableRow table.Row
	for _, m := range modules {
		tableRow = append(tableRow, m)
	}
	tw := commons.ReporterWriter.NewTableWriter()
	tw.AppendHeader(tableRow)
	var max int
	for _, item := range items {
		if max < len(item) {
			max = len(item)
		}
	}
	for row := 0; row < max; row++ {
		var tableRow table.Row
		for _, col := range items {
			if row < len(col) {
				tableRow = append(tableRow, col[row])
			} else {
				tableRow = append(tableRow, "")
			}
		}
		tw.AppendRow(tableRow)
		tw.AppendSeparator()
	}
	content.Txt = strings.Join([]string{titleContent.Txt, tw.Render()}, stringutil.STR_NEWLINE)
	content.Markdown = strings.Join([]string{titleContent.Markdown, tw.RenderMarkdown()}, stringutil.STR_NEWLINE)
	content.HTML = strings.Join([]string{titleContent.HTML, tw.RenderHTML()}, stringutil.STR_NEWLINE)
	return
}

func (r *YTCReport) addSummary(content, overview, items reporter.ReportContent) reporter.ReportContent {
	content.Txt = strings.Join([]string{overview.Txt, items.Txt, content.Txt}, stringutil.STR_NEWLINE)
	content.Markdown = strings.Join([]string{overview.Markdown, items.Markdown, content.Markdown}, stringutil.STR_NEWLINE)
	content.HTML = strings.Join([]string{overview.HTML, items.HTML, content.HTML}, stringutil.STR_NEWLINE)
	return content
}
