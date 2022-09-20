package components

import (
	"sort"

	"github.com/google/uuid"
)

type RailroadConnectionState int8

const (
	ConnectionNotSet RailroadConnectionState = iota
	ConnectionSettingSwitches
	ConnectionSettingSignals
	ConnectionSet
	ConnectionDesolvingSignals
)

type RailroadPath struct {
	Blocks         []*SubBlock
	Switches       []*RailroadSwitch
	SwitchesStates []bool
	Score          int
	Length         int
}

type RailroadConnection struct {
	ID             string                  `json:"id"`
	StartingSignal *Signal                 `json:"starting_signal"`
	EndingSignal   *Signal                 `json:"ending_signal"`
	Blocks         []*SubBlock             `json:"blocks"`
	Switches       []*RailroadSwitch       `json:"switches"`
	State          RailroadConnectionState `json:"state"`
}

var serverRailroadConnections = []RailroadConnection{}

func (connection *RailroadConnection) Desolve() {
	SetConnectionBlocks(connection, false)
	SetConnectionSignals(connection, false)
	connection.State = ConnectionDesolvingSignals
}

func GenerateConnectionUUID() string {
	return uuid.New().String()
}

func SetConnectionBlocks(connection *RailroadConnection, state bool) {
	for index := range connection.Blocks {
		connection.Blocks[index].Occupied = state
	}

	for index := range connection.Switches {
		connection.Switches[index].Occupied = state
	}
}

func SetConnectionSwitches(connection *RailroadConnection, switchesStates []bool) {
	for index := range connection.Switches {
		connection.Switches[index].Set(switchesStates[index])
	}
}

func CheckConnectionSwitchesAcknowledged(connection *RailroadConnection, nextState RailroadConnectionState) {
	for index := range connection.Switches {
		if !connection.Switches[index].Acknowledged {
			return
		}
	}
	connection.State = nextState
}

func SetConnectionSignals(connection *RailroadConnection, state bool) {
	connection.StartingSignal.Set(state)
	for index := range connection.StartingSignal.DistantSignals {
		connection.StartingSignal.DistantSignals[index].Set(connection.EndingSignal.State)
	}
	for index := range connection.Blocks {
		for vindex := range connection.Blocks[index].DistantSignals {
			connection.Blocks[index].DistantSignals[vindex].Set(connection.EndingSignal.State)
		}
	}
}

func CheckConnectionSignalsAcknowledged(connection *RailroadConnection, nextState RailroadConnectionState) {
	if !connection.StartingSignal.Acknowledged {
		return
	}
	for _, distant_signal := range connection.StartingSignal.DistantSignals {
		if !distant_signal.Acknowledged {
			return
		}
	}
	for _, block := range connection.Blocks {
		for _, distant_signal := range block.DistantSignals {
			if !distant_signal.Acknowledged {
				return
			}
		}
	}
	connection.State = nextState
}

func GetConnectionByID(id string) *RailroadConnection {
	for index, connection := range serverRailroadConnections {
		if connection.ID == id {
			return &serverRailroadConnections[index]
		}
	}
	return nil
}

func GetConnectionBySignals(signal1 *Signal, signal2 *Signal) *RailroadConnection {
	for index, connection := range serverRailroadConnections {
		if connection.StartingSignal == signal1 && connection.EndingSignal == signal2 {
			return &serverRailroadConnections[index]
		}
	}
	return nil
}

func CheckConnections() {
	for index, connection := range serverRailroadConnections {
		if connection.State == ConnectionSettingSwitches {
			CheckConnectionSwitchesAcknowledged(&serverRailroadConnections[index], ConnectionSettingSignals)
			if serverRailroadConnections[index].State == ConnectionSettingSignals {
				SetConnectionSignals(&serverRailroadConnections[index], true)
			}
		}
		if connection.State == ConnectionSettingSignals {
			CheckConnectionSignalsAcknowledged(&serverRailroadConnections[index], ConnectionSet)
		}

		if connection.State == ConnectionDesolvingSignals {
			CheckConnectionSignalsAcknowledged(&serverRailroadConnections[index], ConnectionNotSet)
		}
	}
}

func GenerateConnection(signal1 *Signal, signal2 *Signal) *RailroadConnection {
	connectionsUnfiltered := PathFinding(nil, nil, signal1, signal2, &RailroadPath{}, false)
	var connections []*RailroadPath
	for index, connection := range connectionsUnfiltered {
		if connection.Score >= 0 {
			connections = append(connections, connectionsUnfiltered[index])
		}
	}
	sort.Slice(connections, func(i, j int) bool {
		if connections[i].Score == connections[j].Score {
			if connections[i].Length < connections[j].Length {
				return true
			} else {
				return false
			}
		} else if connections[i].Score < connections[j].Score {
			return true
		} else {
			return false
		}
	})
	if len(connections) == 0 {
		return nil
	} else {
		connectionPath := connections[0]
		connection := RailroadConnection{
			ID:             GenerateConnectionUUID(),
			StartingSignal: signal1,
			EndingSignal:   signal2,
			Blocks:         connectionPath.Blocks,
			Switches:       connectionPath.Switches,
			State:          ConnectionSettingSwitches,
		}
		SetConnectionSwitches(&connection, connectionPath.SwitchesStates)
		CheckConnectionSwitchesAcknowledged(&connection, ConnectionSettingSignals)
		if connection.State == ConnectionSettingSignals {
			SetConnectionSignals(&connection, true)
			CheckConnectionSignalsAcknowledged(&connection, ConnectionSet)
		}
		SetConnectionBlocks(&connection, true)

		serverRailroadConnections = append(serverRailroadConnections, connection)
		return &connection
	}
}

func returnPathNotExist(fahrstrasse *RailroadPath) []*RailroadPath {
	fahrstrasse.Score = -1
	return []*RailroadPath{}
}

func PathFinding(
	block *SubBlock,
	rswitch *RailroadSwitch,
	signal *Signal,
	find *Signal,
	connection *RailroadPath,
	direction bool,
) []*RailroadPath {

	if (block != nil && block.Occupied) || (rswitch != nil && rswitch.Occupied) {
		return returnPathNotExist(connection)
	}

	if block != nil {
		connection.Length += block.Length
		connection.Blocks = append(connection.Blocks, block)
	}
	connection.Score += 1

	var next_signal *Signal
	var next_switch *RailroadSwitch

	if block != nil {
		next_signal = block.EndingSignal
		next_switch = block.EndingSwitch
	}

	if signal != nil {
		if signal.FollowingSwitch != nil {
			next_switch = signal.FollowingSwitch
			// next_switch direction MUST BE false
		} else if signal.FollowingBlock != nil {
			d1 := false
			if signal.FollowingBlock.StartingSignal != signal {
				d1 = true
			}
			return PathFinding(signal.FollowingBlock, nil, nil, find, connection, d1)
		}
	}

	if rswitch != nil {
		connection.Switches = append(connection.Switches, rswitch)
		if direction {
			next_signal = rswitch.PreviousSignal
			next_switch = rswitch.PreviousSwitch
		} else {
			straightBlade := connection
			bendingBlade := *connection
			straightBlade.SwitchesStates = append(straightBlade.SwitchesStates, false)
			d1 := false
			if rswitch.FollowingSwitchStraightBlade != nil && rswitch.FollowingSwitchStraightBlade.PreviousSwitch != rswitch {
				d1 = true
			}
			newC1 := PathFinding(GetSubBlockFromBlock(rswitch.FollowingBlockStraightBlade, rswitch, nil), rswitch.FollowingSwitchStraightBlade, nil, find, straightBlade, d1)
			bendingBlade.SwitchesStates = append(bendingBlade.SwitchesStates, true)
			d2 := false
			if rswitch.FollowingSwitchBendingBlade != nil && rswitch.FollowingSwitchBendingBlade.PreviousSwitch != rswitch {
				d2 = true
			}
			newC2 := PathFinding(GetSubBlockFromBlock(rswitch.FollowingBlockBendingBlade, rswitch, nil), rswitch.FollowingSwitchBendingBlade, nil, find, &bendingBlade, d2)
			newC1 = append(newC1, newC2...)
			return newC1
		}
	}

	if next_signal != nil {
		if direction {
			connection.Score += 1
			d := false
			if next_signal.PreviousBlock.EndingSignal == next_signal {
				d = true
			}

			return PathFinding(next_signal.PreviousBlock, nil, nil, find, connection, d)
		}
		if next_signal == find {
			return []*RailroadPath{connection}
		} else {
			return returnPathNotExist(connection)
		}
	} else if next_switch != nil {
		if next_switch.PreviousSwitch == rswitch || next_switch.PreviousSignal == signal || next_switch.PreviousBlock == block.Block {
			// direction = false
			return PathFinding(nil, next_switch, nil, find, connection, false)
		} else {
			connection.Switches = append(connection.Switches, next_switch)
			// direction = true
			switchState := false
			if block != nil && GetSubBlockFromBlock(next_switch.FollowingBlockStraightBlade, next_switch, nil) != block {
				switchState = true
			} else if rswitch != nil && next_switch.FollowingSwitchStraightBlade != rswitch {
				switchState = true
			}
			connection.SwitchesStates = append(connection.SwitchesStates, switchState)
			followingBlock := GetSubBlockFromBlock(next_switch.FollowingBlockStraightBlade, next_switch, nil)
			followingSwitch := next_switch.FollowingSwitchStraightBlade
			if switchState {
				followingBlock = GetSubBlockFromBlock(next_switch.FollowingBlockBendingBlade, next_switch, nil)
				followingSwitch = next_switch.FollowingSwitchBendingBlade
			}
			return PathFinding(followingBlock, followingSwitch, nil, find, connection, true)
		}
	}
	return returnPathNotExist(connection)
}
