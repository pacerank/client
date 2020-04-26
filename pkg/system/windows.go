package system

type windows struct{}

func (w *windows) Processes() ([]Process, error) {
	return []Process{}, nil
}
