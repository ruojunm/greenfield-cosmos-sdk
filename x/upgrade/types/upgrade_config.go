package types

import (
	"math"

	"github.com/cosmos/cosmos-sdk/types"
)

const (
	// EnablePublicDelegationUpgrade is the upgrade name for enabling public delegation
	EnablePublicDelegationUpgrade = types.EnablePublicDelegationUpgrade

	// Nagqu is the upgrade name for Nagqu upgrade
	Nagqu = types.Nagqu

	// Pampas is the upgrade name for Pampas upgrade
	Pampas = types.Pampas
)

// The default upgrade config for networks
var (
	MainnetChainID = "greenfield_1017-1"
	MainnetConfig  = NewUpgradeConfig().SetPlan(&Plan{
		Name:   Nagqu,
		Height: 1,
		Info:   "Nagqu hardfork",
	})

	TestnetChainID = "greenfield_5600-1"
	TestnetConfig  = NewUpgradeConfig().SetPlan(&Plan{
		Name:   Nagqu,
		Height: 471350,
		Info:   "Nagqu hardfork",
	}).SetPlan(&Plan{
		Name:   Pampas,
		Height: 2427233,
		Info:   "Pampas hardfork",
	})
)

func NewUpgradeConfig() *UpgradeConfig {
	return &UpgradeConfig{
		keys:     make(map[string]*key),
		elements: make(map[int64][]*Plan),
	}
}

type key struct {
	index  int
	height int64
}

// UpgradeConfig is a list of upgrade plans
type UpgradeConfig struct {
	keys     map[string]*key
	elements map[int64][]*Plan
}

// SetPlan sets a new upgrade plan
func (c *UpgradeConfig) SetPlan(plan *Plan) *UpgradeConfig {
	if key, ok := c.keys[plan.Name]; ok {
		if c.elements[key.height][key.index].Height == plan.Height {
			*c.elements[key.height][key.index] = *plan
			return c
		}

		c.elements[key.height] = append(c.elements[key.height][:key.index], c.elements[key.height][key.index+1:]...)
	}

	c.elements[plan.Height] = append(c.elements[plan.Height], plan)
	c.keys[plan.Name] = &key{height: plan.Height, index: len(c.elements[plan.Height]) - 1}

	return c
}

// Clear removes all upgrade plans at a given height
func (c *UpgradeConfig) Clear(height int64) {
	for _, plan := range c.elements[height] {
		delete(c.keys, plan.Name)
	}
	c.elements[height] = nil
}

// GetPlan returns the upgrade plan at a given height
func (c *UpgradeConfig) GetPlan(height int64) []*Plan {
	plans, exist := c.elements[height]
	if exist && len(plans) != 0 {
		return plans
	}

	// get recent upgrade plan
	recentHeight := int64(math.MaxInt64)
	for vHeight, vPlans := range c.elements {
		if vHeight > height && vHeight < recentHeight {
			plans = vPlans
			recentHeight = vHeight
		}
	}
	return plans
}
