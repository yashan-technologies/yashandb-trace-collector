package baseinforeporter

import (
	"encoding/json"
	"fmt"

	"ytc/i18n"
	"ytc/internal/modules/ytc/collect/commons/datadef"
	"ytc/internal/modules/ytc/collect/commons/i18nnames"
	"ytc/internal/modules/ytc/collect/data/reporter/commons"
	"ytc/internal/modules/ytc/collect/resultgenner/reporter"

	"git.yasdb.com/go/yaserr"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/shirou/gopsutil/cpu"
)

// validate interface
var _ commons.Reporter = (*HostCPUReporter)(nil)

type HostCPUReporter struct{}

func NewHostCPUReporter() HostCPUReporter {
	return HostCPUReporter{}
}

// [Interface Func]
func (r HostCPUReporter) Report(item datadef.YTCItem, titlePrefix string) (content reporter.ReportContent, err error) {
	title := fmt.Sprintf("%s %s", titlePrefix, i18nnames.GetBaseInfoItemName(item.Name))
	fontSize := reporter.FONT_SIZE_H2

	// report error
	if len(item.Error) != 0 {
		ew := commons.ReporterWriter.NewErrorWriter(item.Error, item.Description)
		content = reporter.GenReportContentByWriterAndTitle(ew, title, fontSize)
		return
	}

	// report cpu info
	cpuInfos, err := r.parseCPUInfos(item)
	if err != nil {
		err = yaserr.Wrapf(err, "parse cpu infos")
		return
	}
	writer := r.genReportContentWriter(cpuInfos)
	content = reporter.GenReportContentByWriterAndTitle(writer, title, fontSize)
	return
}

func (r HostCPUReporter) parseCPUInfos(cpuInfoItem datadef.YTCItem) (cpuInfos []cpu.InfoStat, err error) {
	cpuInfos, ok := cpuInfoItem.Details.([]cpu.InfoStat)
	if !ok {
		tmp, ok := cpuInfoItem.Details.([]map[string]interface{})
		if !ok {
			err = &commons.ErrInterfaceTypeNotMatch{
				Key: cpuInfoItem.Name,
				Targets: []interface{}{
					[]cpu.InfoStat{},
					[]map[string]interface{}{},
				},
				Current: cpuInfoItem.Details,
			}
			err = yaserr.Wrapf(err, "convert cpu info interface")
			return
		}
		data, _ := json.Marshal(tmp)
		if err = json.Unmarshal(data, &cpuInfos); err != nil {
			err = yaserr.Wrapf(err, "unmarshal cpu info")
			return
		}
	}
	return
}

func (r HostCPUReporter) genReportContentWriter(cpuInfos []cpu.InfoStat) reporter.Writer {
	tw := commons.ReporterWriter.NewTableWriter()
	tw.AppendHeader(table.Row{i18n.T("report.table_check_item"), i18n.T("report.table_check_result")})

	tw.AppendRow(table.Row{i18n.T("report.field_cpu_model"), cpuInfos[0].ModelName})
	tw.AppendSeparator()

	var physicalCores, logicalCores int
	tmp := make(map[string]struct{})
	for _, c := range cpuInfos {
		tmp[c.PhysicalID] = struct{}{}
		logicalCores += int(c.Cores)
	}
	physicalCores = len(tmp)
	tw.AppendRow(table.Row{i18n.T("report.field_cpu_cores"), fmt.Sprintf("%s: %d, %s: %d", i18n.T("report.field_physical_cores"), physicalCores, i18n.T("report.field_logical_cores"), logicalCores)})
	tw.AppendSeparator()

	tw.AppendRow(table.Row{i18n.T("report.field_cpu_frequency"), fmt.Sprintf("@%.2fGHz", cpuInfos[0].Mhz/1000)})
	tw.AppendSeparator()

	tw.AppendRow(table.Row{i18n.T("report.field_vendor_id"), cpuInfos[0].VendorID})

	return tw
}
