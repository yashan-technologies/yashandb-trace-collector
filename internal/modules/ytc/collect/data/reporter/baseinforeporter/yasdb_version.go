package baseinforeporter

import (
	"fmt"

	"ytc/i18n"
	"ytc/internal/modules/ytc/collect/commons/datadef"
	"ytc/internal/modules/ytc/collect/commons/i18nnames"
	"ytc/internal/modules/ytc/collect/data/reporter/commons"
	"ytc/internal/modules/ytc/collect/resultgenner/reporter"
)

// validate interface
var _ commons.Reporter = (*YashanDBVersionReporter)(nil)

type YashanDBVersionReporter struct{}

func NewYashanDBVersionReporter() YashanDBVersionReporter {
	return YashanDBVersionReporter{}
}

// [Interface Func]
func (r YashanDBVersionReporter) Report(item datadef.YTCItem, titlePrefix string) (content reporter.ReportContent, err error) {
	title := fmt.Sprintf("%s %s", titlePrefix, i18nnames.GetBaseInfoItemName(item.Name))
	fontSize := reporter.FONT_SIZE_H2

	// report error
	if len(item.Error) != 0 {
		ew := commons.ReporterWriter.NewErrorWriter(item.Error, item.Description)
		content = reporter.GenReportContentByWriterAndTitle(ew, title, fontSize)
		return
	}

	// report yasdb version
	version, err := commons.ParseString(item.Name, item.Details, "parse yasdb version")
	if err != nil {
		return
	}
	writer := r.genReportContentWriter(version)
	content = reporter.GenReportContentByWriterAndTitle(writer, title, fontSize)
	return
}

func (r YashanDBVersionReporter) genReportContentWriter(version string) reporter.Writer {
	return commons.GenStringWriter(i18n.T("report.table_version_info"), version)
}
