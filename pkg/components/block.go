package components

type ReservedType = uint

const (
	NotReserved = iota
	Reserving
	Reserved
)

type Block struct {
	Name     string `json:"name"`
	Reserved uint   `json:"reserved"`
	Occupied bool   `json:"occupied"`
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

func OccupyBlock(
	from_block *Block,
	from_switch *RailroadSwitch,
	to_block *Block,
	to_switch *RailroadSwitch,
) {
	if from_block != nil {
		from_block.Occupied = false
	}

	if from_switch != nil {
		from_switch.Occupied = false
	}

	if to_block != nil {
		to_block.Occupied = true
	}

	if to_switch != nil {
		to_switch.Occupied = true
	}

	signal := GetSignalByFollowingBlock(from_block, from_switch, to_block, to_switch)
	if signal != nil {
		connection := GetConnectionByEndingSignal(signal)
		if connection != nil {
			connection.Desolve()
		}
	}

}
