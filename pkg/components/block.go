package components

type Block struct {
	Name     string `json:"name"`
	Reversed bool   `json:"reversed"`
	Length   int    `json:"length"`
}

type SubBlock struct {
	Name           string  `json:"name"`
	StartingSignal *Signal `json:"starting_signal"`
	EndingSignal   *Signal `json:"ending_signal"`

	StartingSwitch *RailroadSwitch  `json:"starting_switch"`
	EndingSwitch   *RailroadSwitch  `json:"ending_switch"`
	DistantSignals []*DistantSignal `json:"distant_signals"`
	*Block
}

func GetBlockByName(name string) *Block {
	for index, element := range blocks {
		if element.Name == name {
			return &blocks[index]
		}
	}
	return nil
}

func GetSubBlockByName(name string) *SubBlock {
	for index, element := range subblocks {
		if element.Name == name {
			return &subblocks[index]
		}
	}
	return nil
}

func GetSubBlockFromBlock(block *Block, starting_switch *RailroadSwitch, starting_signal *Signal) *SubBlock {
	if block == nil {
		return nil
	}
	for index, element := range subblocks {
		if element.Block == block && ((starting_switch != nil && element.StartingSwitch == starting_switch) || (starting_signal != nil && element.StartingSignal == starting_signal)) {
			return &subblocks[index]
		}
	}
	return nil
}
