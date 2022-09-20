package components

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type DistantSignals interface{}
type SignalConfig struct {
	DistantSignals []string `yaml:"distant_signals"`
}
type SwitchConfig interface{}
type BlockConfig interface{}
type SubblockConfig struct {
	Start          string   `yaml:"start"`
	End            string   `yaml:"end"`
	DistantSignals []string `yaml:"distant_signals"`
}

type SignalRelationsConfig struct {
	Following string `yaml:"following"`
	Previous  string `yaml:"previous"`
}

type SwitchRelationsConfig struct {
	Previous               string `yaml:"previous"`
	FollowingStraightBlade string `yaml:"following_straight_blade"`
	FollowingBendingBlade  string `yaml:"following_bending_blade"`
}

type RelationsConfig struct {
	Signals  map[string]SignalRelationsConfig `yaml:"signals"`
	Switches map[string]SwitchRelationsConfig `yaml:"switches"`
}

type TemplateConfig struct {
	DistantSignals map[string]interface{}    `yaml:"distant_signals"`
	Signals        map[string]SignalConfig   `yaml:"signals"`
	Switches       map[string]SwitchConfig   `yaml:"switches"`
	Blocks         map[string]BlockConfig    `yaml:"blocks"`
	Subblocks      map[string]SubblockConfig `yaml:"subblocks"`
	Relations      RelationsConfig           `yaml:"relations"`
}

func ParseRelationsConfig(config *RelationsConfig) {
	for name, signal := range config.Signals {
		if signal.Following != "" {
			if strings.HasPrefix(signal.Following, "W") {
				GetSignalByName(name).FollowingSwitch = GetSwitchByName(signal.Following)
			} else if strings.HasPrefix(signal.Following, "B") {
				GetSignalByName(name).FollowingBlock = GetSubBlockByName(signal.Following)
			}
		}
		if signal.Previous != "" {
			if strings.HasPrefix(signal.Previous, "W") {
				GetSignalByName(name).PreviousSwitch = GetSwitchByName(signal.Previous)
			} else if strings.HasPrefix(signal.Previous, "B") {
				GetSignalByName(name).PreviousBlock = GetSubBlockByName(signal.Previous)
			}
		}

	}

	for name, rswitch := range config.Switches {
		if rswitch.Previous != "" {
			if strings.HasPrefix(rswitch.Previous, "W") {
				GetSwitchByName(name).PreviousSwitch = GetSwitchByName(rswitch.Previous)
			} else if strings.HasPrefix(rswitch.Previous, "S") {
				GetSwitchByName(name).PreviousSignal = GetSignalByName(rswitch.Previous)
			} else if strings.HasPrefix(rswitch.Previous, "B") {
				GetSwitchByName(name).PreviousBlock = GetBlockByName(rswitch.Previous)
			}
		}
		if rswitch.FollowingStraightBlade != "" {
			if strings.HasPrefix(rswitch.FollowingStraightBlade, "W") {
				GetSwitchByName(name).FollowingSwitchStraightBlade = GetSwitchByName(rswitch.FollowingStraightBlade)
			} else if strings.HasPrefix(rswitch.FollowingStraightBlade, "B") {
				GetSwitchByName(name).FollowingBlockStraightBlade = GetBlockByName(rswitch.FollowingStraightBlade)
			}
		}
		if rswitch.FollowingBendingBlade != "" {
			if strings.HasPrefix(rswitch.FollowingBendingBlade, "W") {
				GetSwitchByName(name).FollowingSwitchBendingBlade = GetSwitchByName(rswitch.FollowingBendingBlade)
			} else if strings.HasPrefix(rswitch.FollowingBendingBlade, "B") {
				GetSwitchByName(name).FollowingBlockBendingBlade = GetBlockByName(rswitch.FollowingBendingBlade)
			}
		}
	}
}

func ParseTemplateConfig(config *TemplateConfig) {
	for name := range config.DistantSignals {
		distant_signals = append(distant_signals, DistantSignal{Name: name})
	}
	for name, signal := range config.Signals {
		ds := make([]*DistantSignal, len(signal.DistantSignals))
		for index, name := range signal.DistantSignals {
			ds[index] = GetDistantSignalByName(name)
		}
		signals = append(signals, Signal{Name: name, DistantSignals: ds})
	}

	for name := range config.Switches {
		switches = append(switches, RailroadSwitch{Name: name})
	}

	for name := range config.Blocks {
		blocks = append(blocks, Block{Name: name})
	}

	for name, block := range config.Subblocks {
		subblocks = append(subblocks, SubBlock{Name: name})
		blockName := make([]rune, 0)
		for _, element := range name {
			if strings.ToUpper(string(element)) == string(element) {
				blockName = append(blockName, element)
			}
		}
		subblocks[len(subblocks)-1].Block = GetBlockByName(string(blockName))
		if strings.HasPrefix(block.Start, "S") {
			subblocks[len(subblocks)-1].StartingSignal = GetSignalByName(block.Start)
		} else if strings.HasPrefix(block.Start, "W") {
			subblocks[len(subblocks)-1].StartingSwitch = GetSwitchByName(block.Start)
		}

		if strings.HasPrefix(block.End, "S") {
			subblocks[len(subblocks)-1].EndingSignal = GetSignalByName(block.End)
		} else if strings.HasPrefix(block.End, "W") {
			subblocks[len(subblocks)-1].EndingSwitch = GetSwitchByName(block.End)
		}
		ds := make([]*DistantSignal, len(block.DistantSignals))
		for index, distant_signal := range block.DistantSignals {
			ds[index] = GetDistantSignalByName(distant_signal)
		}
		subblocks[len(subblocks)-1].DistantSignals = ds
	}
	ParseRelationsConfig(&config.Relations)

}

func ReadTemplating() error {
	template_file_path := os.Getenv("TEMPLATE_YAML")
	if template_file_path == "" {
		template_file_path = "template.open-interlocking.yml"
	}
	fmt.Printf("%s\n", template_file_path)
	template_file, err := os.ReadFile(template_file_path)
	if err != nil {
		return err
	}
	var template_config TemplateConfig
	err = yaml.Unmarshal(template_file, &template_config)
	ParseTemplateConfig(&template_config)
	return nil
}
