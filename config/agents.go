package config

import (
	"fmt"
	"github.com/forta-protocol/forta-node/protocol"
	"github.com/forta-protocol/forta-node/utils"
)

type AgentConfig struct {
	ID         string  `yaml:"id" json:"id"`
	Image      string  `yaml:"image" json:"image"`
	Manifest   string  `yaml:"manifest" json:"manifest"`
	IsLocal    bool    `yaml:"isLocal" json:"isLocal"`
	StartBlock *uint64 `yaml:"startBlock" json:"startBlock,omitempty"`
	StopBlock  *uint64 `yaml:"stopBlock" json:"stopBlock,omitempty"`
}

// ToAgentInfo transforms the agent config to the agent info.
func (ac AgentConfig) ToAgentInfo() *protocol.AgentInfo {
	return &protocol.AgentInfo{
		Id:        ac.ID,
		Image:     ac.Image,
		ImageHash: ac.ImageHash(),
		IsTest:    ac.IsLocal,
		Manifest:  ac.Manifest,
	}
}

func (ac AgentConfig) ImageHash() string {
	_, digest := utils.SplitImageRef(ac.Image)
	return digest
}

func (ac AgentConfig) ContainerName() string {
	_, digest := utils.SplitImageRef(ac.Image)
	return fmt.Sprintf("%s-agent-%s-%s", ContainerNamePrefix, utils.ShortenString(ac.ID, 8), utils.ShortenString(digest, 4))
}

func (ac AgentConfig) GrpcPort() string {
	return "50051"
}