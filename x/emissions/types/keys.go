package types

import "cosmossdk.io/collections"

const (
	ModuleName                                 = "emissions"
	StoreKey                                   = "emissions"
	AlloraStakingAccountName                   = "allorastaking"
	AlloraRequestsAccountName                  = "allorarequests"
	AlloraRewardsAccountName                   = "allorarewards"
	AlloraPendingRewardForDelegatorAccountName = "allorapendingrewards"
)

var (
	ParamsKey                                   = collections.NewPrefix(0)
	TotalStakeKey                               = collections.NewPrefix(1)
	TopicStakeKey                               = collections.NewPrefix(2)
	RewardsKey                                  = collections.NewPrefix(3)
	NextTopicIdKey                              = collections.NewPrefix(4)
	TopicsKey                                   = collections.NewPrefix(5)
	TopicWorkersKey                             = collections.NewPrefix(6)
	TopicReputersKey                            = collections.NewPrefix(7)
	DelegatorStakeKey                           = collections.NewPrefix(8)
	DelegateStakePlacementKey                   = collections.NewPrefix(9)
	TargetStakeKey                              = collections.NewPrefix(10)
	InferencesKey                               = collections.NewPrefix(11)
	ForecastsKey                                = collections.NewPrefix(12)
	WorkerNodesKey                              = collections.NewPrefix(13)
	ReputerNodesKey                             = collections.NewPrefix(14)
	LatestInferencesTsKey                       = collections.NewPrefix(15)
	ActiveTopicsKey                             = collections.NewPrefix(16)
	AllInferencesKey                            = collections.NewPrefix(17)
	AllForecastsKey                             = collections.NewPrefix(18)
	AllLossBundlesKey                           = collections.NewPrefix(19)
	StakeRemovalKey                             = collections.NewPrefix(20)
	StakeByReputerAndTopicId                    = collections.NewPrefix(21)
	DelegateStakeRemovalKey                     = collections.NewPrefix(22)
	AllTopicStakeSumKey                         = collections.NewPrefix(23)
	AddressTopicsKey                            = collections.NewPrefix(24)
	WhitelistAdminsKey                          = collections.NewPrefix(24)
	ChurnReadyTopicsKey                         = collections.NewPrefix(25)
	NetworkLossBundlesKey                       = collections.NewPrefix(26)
	NetworkRegretsKey                           = collections.NewPrefix(27)
	StakeByReputerAndTopicIdKey                 = collections.NewPrefix(28)
	ReputerScoresKey                            = collections.NewPrefix(29)
	InferenceScoresKey                          = collections.NewPrefix(30)
	ForecastScoresKey                           = collections.NewPrefix(31)
	ReputerListeningCoefficientKey              = collections.NewPrefix(32)
	InfererNetworkRegretsKey                    = collections.NewPrefix(33)
	ForecasterNetworkRegretsKey                 = collections.NewPrefix(34)
	OneInForecasterNetworkRegretsKey            = collections.NewPrefix(35)
	UnfulfilledWorkerNoncesKey                  = collections.NewPrefix(36)
	UnfulfilledReputerNoncesKey                 = collections.NewPrefix(37)
	FeeRevenueEpochKey                          = collections.NewPrefix(38)
	TopicFeeRevenueKey                          = collections.NewPrefix(39)
	PreviousTopicWeightKey                      = collections.NewPrefix(40)
	PreviousReputerRewardFractionKey            = collections.NewPrefix(41)
	PreviousInferenceRewardFractionKey          = collections.NewPrefix(42)
	PreviousForecastRewardFractionKey           = collections.NewPrefix(43)
	LatestInfererScoresByWorkerKey              = collections.NewPrefix(44)
	LatestForecasterScoresByWorkerKey           = collections.NewPrefix(45)
	LatestReputerScoresByReputerKey             = collections.NewPrefix(46)
	TopicRewardNonceKey                         = collections.NewPrefix(47)
	DelegateRewardPerShare                      = collections.NewPrefix(48)
	PreviousPercentageRewardToStakedReputersKey = collections.NewPrefix(49)
)
