package diagreporter

import (
	"encoding/json"
	"fmt"

	"ytc/i18n"
	"ytc/internal/modules/ytc/collect/commons/datadef"
	"ytc/internal/modules/ytc/collect/commons/i18nnames"
	"ytc/internal/modules/ytc/collect/data/reporter/commons"
	"ytc/internal/modules/ytc/collect/resultgenner/reporter"
	"ytc/utils/processutil"

	"git.yasdb.com/go/yaserr"
	"github.com/jedib0t/go-pretty/v6/table"
)

// validate interface
var _ commons.Reporter = (*YashanDBProcessStatusReporter)(nil)

type YashanDBProcessStatusReporter struct{}

func NewYashanDBProcessStatusReporter() YashanDBProcessStatusReporter {
	return YashanDBProcessStatusReporter{}
}

// [Interface Func]
func (r YashanDBProcessStatusReporter) Report(item datadef.YTCItem, titlePrefix string) (content reporter.ReportContent, err error) {
	title := fmt.Sprintf("%s %s", titlePrefix, i18nnames.GetDiagItemName(item.Name))
	fontSize := reporter.FONT_SIZE_H2

	// report error
	if len(item.Error) != 0 {
		ew := commons.ReporterWriter.NewErrorWriter(item.Error, item.Description)
		content = reporter.GenReportContentByWriterAndTitle(ew, title, fontSize)
		return
	}

	// report yasdb process status
	process, err := r.parseYashanDBProcess(item)
	if err != nil {
		err = yaserr.Wrapf(err, "parse yasdb process")
		return
	}
	writer := r.genReportContentWriter(process)
	content = reporter.GenReportContentByWriterAndTitle(writer, title, fontSize)
	return
}

func (r YashanDBProcessStatusReporter) parseYashanDBProcess(item datadef.YTCItem) (process processutil.Process, err error) {
	process, ok := item.Details.(processutil.Process)
	if !ok {
		tmp, ok := item.Details.(map[string]interface{})
		if !ok {
			err = &commons.ErrInterfaceTypeNotMatch{
				Key: item.Name,
				Targets: []interface{}{
					&processutil.Process{},
					map[string]interface{}{},
				},
				Current: item.Details,
			}
			err = yaserr.Wrapf(err, "parse process interface")
			return
		}
		data, _ := json.Marshal(tmp)
		if err = json.Unmarshal(data, &process); err != nil {
			err = yaserr.Wrapf(err, "unmarshal process info")
			return
		}
	}
	return
}

func (r YashanDBProcessStatusReporter) genReportContentWriter(p processutil.Process) reporter.Writer {
	tw := commons.ReporterWriter.NewTableWriter()
	tw.AppendHeader(table.Row{
		i18n.T("report.process_id"),
		i18n.T("report.process_command"),
		i18n.T("report.process_user"),
		i18n.T("report.process_create_time"),
		i18n.T("report.process_cpu_usage"),
		i18n.T("report.process_memory_usage"),
		i18n.T("report.process_status"),
	})

	tw.AppendRow(table.Row{
		p.Pid,
		p.ReadableCmdline,
		p.User,
		p.CreateTime,
		fmt.Sprintf("%.2f%%", p.CPUPercent),
		fmt.Sprintf("%.2f%%", p.MemoryPercent),
		p.Status,
	})
	return tw
}
