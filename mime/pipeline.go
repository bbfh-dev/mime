package mime

func Pipeline(calls ...func() error) error {
	for _, call := range calls {
		if call == nil {
			continue
		}
		if err := call(); err != nil {
			return err
		}
	}
	return nil
}
