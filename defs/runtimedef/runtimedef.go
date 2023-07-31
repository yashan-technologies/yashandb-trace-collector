package runtimedef

func InitRuntime() error {
	if err := initYTCHome(); err != nil {
		return err
	}
	if err := initExecuteable(); err != nil {
		return err
	}
	return nil
}
