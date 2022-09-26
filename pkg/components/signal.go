package components

import "github.com/adridevelopsthings/open-interlocking/pkg/config"

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
		signal.Acknowledged = config.IgnoreAcknowledgements
	}
}

func (distant_signal *DistantSignal) Set(state bool) {
	if state != distant_signal.State {
		distant_signal.State = state
		distant_signal.Acknowledged = config.IgnoreAcknowledgements
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

func GetSignalByFollowingBlock(
	from_block *Block,
	from_switch *RailroadSwitch,
	to_block *Block,
	to_switch *RailroadSwitch,
) *Signal {
	for index, element := range signals {
		if ((to_block != nil && element.FollowingBlock != nil && to_block == element.FollowingBlock.Block) ||
			(to_switch != nil && to_switch == element.FollowingSwitch)) &&
			((from_block != nil && element.PreviousBlock != nil && from_block == element.PreviousBlock.Block) ||
				(from_switch != nil && from_switch == element.PreviousSwitch)) {
			return &signals[index]
		}
	}
	return nil
}
