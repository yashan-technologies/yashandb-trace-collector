package flags

import (
	"fmt"

	"ytc/defs/compiledef"
	"ytc/i18n"

	"git.yasdb.com/go/yasutil/tabler"
	"github.com/alecthomas/kong"
)

type showFlag bool

// [Interface Func]
// BeforeReset shows software compilation information and terminates with a 0 exit status.
func (s showFlag) BeforeReset(app *kong.Kong, vars kong.Vars) error {
	fmt.Fprint(app.Stdout, s.genContent())
	app.Exit(0)
	return nil
}

// genContent generates data of software compilation information.
func (s showFlag) genContent() string {
	titles := []*tabler.RowTitle{
		{Name: i18n.T("app.header_key")},
		{Name: i18n.T("app.header_value")},
	}
	table := tabler.NewTable(i18n.T("app.info"), titles...)
	_ = table.AddColumn(i18n.T("app.version"), compiledef.GetAPPVersion())
	_ = table.AddColumn(i18n.T("app.go_version"), compiledef.GetGoVersion())
	_ = table.AddColumn(i18n.T("app.git_commit"), compiledef.GetGitCommitID())
	_ = table.AddColumn(i18n.T("app.git_describe"), compiledef.GetGitDescribe())
	return table.String()
}
