package components

type FullComponents struct {
	Signals  *[]Signal         `json:"signals"`
	Blocks   *[]Block          `json:"blocks"`
	Switches *[]RailroadSwitch `json:"switches"`
}

func GetFullComponents() *FullComponents {
	return &FullComponents{
		Signals:  &signals,
		Blocks:   &blocks,
		Switches: &switches,
	}
}
