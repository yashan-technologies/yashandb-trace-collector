package baseinforeporter

import (
	"encoding/json"
	"fmt"
	"time"

	"ytc/defs/timedef"
	"ytc/i18n"
	"ytc/internal/modules/ytc/collect/commons/datadef"
	"ytc/internal/modules/ytc/collect/commons/i18nnames"
	"ytc/internal/modules/ytc/collect/data/reporter/commons"
	"ytc/internal/modules/ytc/collect/resultgenner/reporter"

	"git.yasdb.com/go/yaserr"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/shirou/gopsutil/host"
)

// validate interface
var _ commons.Reporter = (*HostOSInfoReporter)(nil)

type HostOSInfoReporter struct{}

func NewHostOSInfoReporter() HostOSInfoReporter {
	return HostOSInfoReporter{}
}

// [Interface Func]
func (r HostOSInfoReporter) Report(item datadef.YTCItem, titlePrefix string) (content reporter.ReportContent, err error) {
	title := fmt.Sprintf("%s %s", titlePrefix, i18nnames.GetBaseInfoItemName(item.Name))
	fontSize := reporter.FONT_SIZE_H2

	// report error
	if len(item.Error) != 0 {
		ew := commons.ReporterWriter.NewErrorWriter(item.Error, item.Description)
		content = reporter.GenReportContentByWriterAndTitle(ew, title, fontSize)
		return
	}

	// report host info
	hostInfo, err := r.parseHostInfo(item)
	if err != nil {
		err = yaserr.Wrapf(err, "parse network info")
		return
	}
	writer := r.genReportContentWriter(hostInfo)
	content = reporter.GenReportContentByWriterAndTitle(writer, title, fontSize)
	return
}

func (r HostOSInfoReporter) parseHostInfo(item datadef.YTCItem) (hostInfo *host.InfoStat, err error) {
	hostInfo, ok := item.Details.(*host.InfoStat)
	if !ok {
		tmp, ok := item.Details.(map[string]interface{})
		if !ok {
			err = &commons.ErrInterfaceTypeNotMatch{
				Key: item.Name,
				Targets: []interface{}{
					&host.InfoStat{},
					map[string]interface{}{},
				},
				Current: item.Details,
			}
			err = yaserr.Wrapf(err, "parse host info interface")
			return
		}
		data, _ := json.Marshal(tmp)
		if err = json.Unmarshal(data, &hostInfo); err != nil {
			err = yaserr.Wrapf(err, "unmarshal host info")
			return
		}
	}
	return
}

func (r HostOSInfoReporter) genReportContentWriter(hostInfo *host.InfoStat) reporter.Writer {
	tw := commons.ReporterWriter.NewTableWriter()
	tw.AppendHeader(table.Row{i18n.T("report.table_check_item"), i18n.T("report.table_check_result")})

	tw.AppendRow(table.Row{i18n.T("report.field_hostname"), hostInfo.Hostname})
	tw.AppendSeparator()

	tw.AppendRow(table.Row{i18n.T("report.field_boot_time"), time.Unix(int64(hostInfo.BootTime), 0).Format(timedef.TIME_FORMAT)})
	tw.AppendSeparator()

	tw.AppendRow(table.Row{i18n.T("report.field_os"), hostInfo.OS})
	tw.AppendSeparator()

	tw.AppendRow(table.Row{i18n.T("report.field_platform"), fmt.Sprintf("%s %s (%s%s)", hostInfo.Platform, hostInfo.PlatformVersion, hostInfo.PlatformFamily, i18n.T("report.field_platform_family"))})
	tw.AppendSeparator()

	tw.AppendRow(table.Row{i18n.T("report.field_kernel_version"), hostInfo.KernelVersion})
	tw.AppendSeparator()

	tw.AppendRow(table.Row{i18n.T("report.field_kernel_arch"), hostInfo.KernelArch})
	return tw
}
