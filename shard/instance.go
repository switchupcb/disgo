package shard

import (
	"fmt"
	"time"

	"github.com/switchupcb/disgo"
)

// InstanceShardManager is a shard manager for a Discord Bot
// that runs on a single instance.
//
// This shard manager routes every shard to every session (1).
type InstanceShardManager struct {
	// Shards represents the number of shards this shard manager uses.
	//
	// When Shards = 0, the automatic shard manager is used.
	Shards int

	// Limit contains information about a client's sharding limits.
	Limit *disgo.ShardLimit

	// Sessions represents a list of sessions sorted by shard_id (in order of connection).
	Sessions []*disgo.Session
}

const (
	// LogCtxShardManager represents the log key for an InstanceShardManager.
	LogCtxShardManager = "shardmanager"
)

const (
	errShardManager = "shardmanager: %w"
)

func (sm *InstanceShardManager) SetNumShards(shards int) {
	sm.Shards = shards
}

func (sm *InstanceShardManager) SetLimit(bot *disgo.Client) (*disgo.GetGatewayBotResponse, error) {
	if sm.Limit == nil {
		gateway := disgo.GetGatewayBot{}
		response, err := gateway.Send(bot)
		if err != nil {
			return nil, fmt.Errorf("shardmanager: Gateway API Endpoint: %w", err)
		}

		sm.Limit = &disgo.ShardLimit{
			Reset:             time.Now().Add(time.Millisecond*time.Duration(response.SessionStartLimit.ResetAfter) + 1),
			MaxStarts:         response.SessionStartLimit.Total,
			RemainingStarts:   response.SessionStartLimit.Remaining,
			MaxConcurrency:    response.SessionStartLimit.MaxConcurrency,
			RecommendedShards: response.Shards,
		}

		return response, nil
	}

	return nil, nil
}

func (sm *InstanceShardManager) GetSessions() []*disgo.Session {
	return sm.Sessions
}

func (sm *InstanceShardManager) Ready(bot *disgo.Client, session *disgo.Session, ready *disgo.Ready) {
	if ready.Shard != nil {
		disgo.Logger.Info().Str(LogCtxShardManager, "received Ready event with nil Shard field.")

		return
	}
}

func (sm *InstanceShardManager) Connect(bot *disgo.Client) error {
	// Determine the number of shards to use.
	//
	// totalShards represents the total number of shards to use.
	totalShards := sm.Shards
	if totalShards <= 0 {
		if _, err := sm.SetLimit(bot); err != nil {
			return fmt.Errorf(errShardManager, err)
		}

		totalShards = sm.Limit.RecommendedShards
	}

	sm.Sessions = make([]*disgo.Session, totalShards)

	// Start the specified number of shards.
	//
	// shardID represents the current number of shards that have been created (minus one).
	shardID := 0
	for shardCount := 0; shardCount < totalShards; shardCount++ {
		session := disgo.NewSession()
		session.Shard = &[2]int{shardID, totalShards}

		// shards must be started in order (by bucket).
		// https://discord.com/developers/docs/topics/gateway#sharding-max-concurrency
		if err := session.Connect(bot); err != nil {
			return fmt.Errorf(errShardManager, err)
		}

		sm.Sessions[shardCount] = session

		shardID++
	}

	return nil
}

// Disconnect disconnects from the Discord Gateway using the Shard Manager.
func (sm *InstanceShardManager) Disconnect() error {
	// totalShards represents the total number of shards that are connected.
	totalShards := len(sm.Sessions)

	for sessionCount := totalShards - 1; sessionCount > -1; sessionCount-- {
		if err := sm.Sessions[sessionCount].Disconnect(); err != nil {
			return fmt.Errorf(errShardManager, err)
		}
	}

	sm.Sessions = nil

	return nil
}

// Reconnect connects to the Discord Gateway using the Shard Manager.
func (sm *InstanceShardManager) Reconnect() error {
	// totalShards represents the total number of shards that are connected.
	totalShards := len(sm.Sessions)

	for sessionCount := 0; sessionCount < totalShards; sessionCount++ {
		if err := sm.Sessions[sessionCount].Reconnect(); err != nil {
			return fmt.Errorf(errShardManager, err)
		}
	}

	return nil
}
