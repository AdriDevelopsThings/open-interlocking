package components

type DistantSignal struct {
	Name         string `json:"name"`
	State        bool   `json:"state"`
	Acknowledged bool   `json:"acknowledged"`
}

type Signal struct {
	Name         string `json:"name"`
	State        bool   `json:"state"`
	Acknowledged bool   `json:"acknowledged"`

	FollowingBlock *SubBlock `json:"-"`
	PreviousBlock  *SubBlock `json:"-"`

	FollowingSwitch *RailroadSwitch `json:"-"`
	PreviousSwitch  *RailroadSwitch `json:"-"`

	DistantSignals []*DistantSignal `json:"distant_signals"`
}

func (signal *Signal) Set(state bool) {
	if state != signal.State {
		signal.State = state
		signal.Acknowledged = false
	}
}

func (distant_signal *DistantSignal) Set(state bool) {
	if state != distant_signal.State {
		distant_signal.State = state
		distant_signal.Acknowledged = false
	}
}

func GetSignalByName(name string) *Signal {
	for index, element := range signals {
		if element.Name == name {
			return &signals[index]
		}
	}
	return nil
}

func GetDistantSignalByName(name string) *DistantSignal {
	for index, element := range distant_signals {
		if element.Name == name {
			return &distant_signals[index]
		}
	}
	return nil
}
