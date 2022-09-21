package components

type RailroadSwitch struct {
	Name string `json:"name"`

	FollowingBlockStraightBlade  *Block          `json:"-"`
	FollowingSwitchStraightBlade *RailroadSwitch `json:"-"`
	FollowingBlockBendingBlade   *Block          `json:"-"`
	FollowingSwitchBendingBlade  *RailroadSwitch `json:"-"`

	PreviousSignal *Signal         `json:"-"`
	PreviousBlock  *Block          `json:"-"`
	PreviousSwitch *RailroadSwitch `json:"-"`

	Reversed bool

	State        bool // false = straight blade; true = bending blade
	Acknowledged bool
}

func (railroad_switch *RailroadSwitch) Set(state bool) {
	if state != railroad_switch.State {
		railroad_switch.State = state
		railroad_switch.Acknowledged = false
	}
}

func GetSwitchByName(name string) *RailroadSwitch {
	for index, element := range switches {
		if element.Name == name {
			return &switches[index]
		}
	}
	return nil
}
