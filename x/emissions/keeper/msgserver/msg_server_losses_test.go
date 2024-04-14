package msgserver_test

import (
	alloraMath "github.com/allora-network/allora-chain/math"
	"github.com/allora-network/allora-chain/x/emissions/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

/*
func (s *KeeperTestSuite) TestMsgInsertBulkReputerPayload() {
	ctx, msgServer := s.ctx, s.msgServer
	require := s.Require()

	// Mock setup for addresses
	reputerAddr := sdk.AccAddress(PKS[0].Address()).String()
	workerAddr := sdk.AccAddress(PKS[1].Address()).String()

	// TODO make this line work
	msgServer.keeper.stakeByReputerAndTopicId.Set(s.ctx, reputerAddr, 100)

	// Create a MsgInsertBulkReputerPayload message
	lossesMsg := &types.MsgInsertBulkReputerPayload{
		Sender: reputerAddr,
		Nonce: &types.Nonce{
			Nonce: 1,
		},
		ReputerValueBundles: []*types.ReputerValueBundle{
			{
				Reputer: reputerAddr,
				ValueBundle: &types.ValueBundle{
					TopicId:       1,
					CombinedValue: alloraMath.NewDecFromInt64(100),
					InfererValues: []*types.WorkerAttributedValue{
						{
							Worker: workerAddr,
							Value:  alloraMath.NewDecFromInt64(100),
						},
					},
					ForecasterValues: []*types.WorkerAttributedValue{
						{
							Worker: workerAddr,
							Value:  alloraMath.NewDecFromInt64(100),
						},
					},
					NaiveValue: alloraMath.NewDecFromInt64(100),
					OneOutInfererValues: []*types.WithheldWorkerAttributedValue{
						{
							Worker: workerAddr,
							Value:  alloraMath.NewDecFromInt64(100),
						},
					},
					OneOutForecasterValues: []*types.WithheldWorkerAttributedValue{
						{
							Worker: workerAddr,
							Value:  alloraMath.NewDecFromInt64(100),
						},
					},
					OneInForecasterValues: []*types.WorkerAttributedValue{
						{
							Worker: workerAddr,
							Value:  alloraMath.NewDecFromInt64(100),
						},
					},
				},
				Signature: []byte("ValueBundle Signature"),
			},
		},
		Signature:      []byte("Nonce + ReputerValueBundles Signature"),
	}

	_, err := msgServer.InsertBulkReputerPayload(ctx, lossesMsg)
	require.NoError(err, "InsertBulkReputerPayload should not return an error")
}
*/

func (s *KeeperTestSuite) TestMsgInsertBulkReputerPayloadInvalidUnauthorized() {
	ctx, msgServer := s.ctx, s.msgServer
	require := s.Require()

	// Mock setup for addresses
	reputerAddr := nonAdminAccounts[0].String()
	workerAddr := sdk.AccAddress(PKS[1].Address()).String()

	// Create a MsgInsertBulkReputerPayload message
	lossesMsg := &types.MsgInsertBulkReputerPayload{
		Sender: reputerAddr,
		ReputerValueBundles: []*types.ReputerValueBundle{
			{
				Reputer: reputerAddr,
				ValueBundle: &types.ValueBundle{
					TopicId:       1,
					CombinedValue: alloraMath.NewDecFromInt64(100),
					InfererValues: []*types.WorkerAttributedValue{
						{
							Worker: workerAddr,
							Value:  alloraMath.NewDecFromInt64(100),
						},
					},
					ForecasterValues: []*types.WorkerAttributedValue{
						{
							Worker: workerAddr,
							Value:  alloraMath.NewDecFromInt64(100),
						},
					},
					NaiveValue: alloraMath.NewDecFromInt64(100),
					OneOutInfererValues: []*types.WithheldWorkerAttributedValue{
						{
							Worker: workerAddr,
							Value:  alloraMath.NewDecFromInt64(100),
						},
					},
					OneOutForecasterValues: []*types.WithheldWorkerAttributedValue{
						{
							Worker: workerAddr,
							Value:  alloraMath.NewDecFromInt64(100),
						},
					},
					OneInForecasterValues: []*types.WorkerAttributedValue{
						{
							Worker: workerAddr,
							Value:  alloraMath.NewDecFromInt64(100),
						},
					},
				},
			},
		},
		Signature: []byte("Nonce + ReputerValueBundles Signature"),
	}

	_, err := msgServer.InsertBulkReputerPayload(ctx, lossesMsg)
	require.ErrorIs(err, types.ErrNotInReputerWhitelist, "InsertBulkReputerPayload should return an error")
}
