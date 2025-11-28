package errdef

import (
	"errors"

	"ytc/i18n"
)

var (
	ErrPathFormat = errors.New("path format error, please check")
)

type ErrFileNotFound struct {
	Fname string
}

type ErrFileParseFailed struct {
	Fname string
	Err   error
}

type ErrCmdNotExist struct {
	Cmd string
}

type ErrCmdNeedRoot struct {
	Cmd string
}

type ErrPermissionDenied struct {
	User     string
	FileName string
}

func NewErrCmdNotExist(cmd string) *ErrCmdNotExist {
	return &ErrCmdNotExist{
		Cmd: cmd,
	}
}

func NewErrCmdNeedRoot(cmd string) *ErrCmdNeedRoot {
	return &ErrCmdNeedRoot{
		Cmd: cmd,
	}
}

func NewErrPermissionDenied(user string, path string) *ErrPermissionDenied {
	return &ErrPermissionDenied{
		User:     user,
		FileName: path,
	}
}

func (e *ErrFileNotFound) Error() string {
	return i18n.TWithData("err.file_not_found", map[string]interface{}{"Fname": e.Fname})
}

func (e *ErrFileParseFailed) Error() string {
	return i18n.TWithData("err.file_parse_failed", map[string]interface{}{"Fname": e.Fname, "Err": e.Err})
}

func (e *ErrCmdNotExist) Error() string {
	return i18n.TWithData("err.cmd_not_exist", map[string]interface{}{"Cmd": e.Cmd})
}

func (e *ErrCmdNeedRoot) Error() string {
	return i18n.TWithData("err.cmd_need_root", map[string]interface{}{"Cmd": e.Cmd})
}

func (e *ErrPermissionDenied) Error() string {
	return i18n.TWithData("err.permission_denied", map[string]interface{}{"User": e.User, "FileName": e.FileName})
}
