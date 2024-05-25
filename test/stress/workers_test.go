package stress_test

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"

	cosmosMath "cosmossdk.io/math"
	"github.com/allora-network/allora-chain/app/params"
	alloraMath "github.com/allora-network/allora-chain/math"
	testCommon "github.com/allora-network/allora-chain/test/common"
	emissionstypes "github.com/allora-network/allora-chain/x/emissions/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/stretchr/testify/require"
)

// creates the worker addresses in the account registry
func createWorkerAddresses(
	m testCommon.TestConfig,
	topicId uint64,
	workersMax int,
) (workers NameToAccountMap) {
	workers = make(map[string]AccountAndAddress)

	for workerIndex := 0; workerIndex < workersMax; workerIndex++ {
		workerAccountName := getWorkerAccountName(workerIndex, topicId)
		workerAccount, _, err := m.Client.AccountRegistryCreate(workerAccountName)
		if err != nil {
			fmt.Println("Error creating funder address: ", workerAccountName, " - ", err)
			continue
		}
		workerAddressToFund, err := workerAccount.Address(params.HumanCoinUnit)
		if err != nil {
			fmt.Println("Error creating funder address: ", workerAccountName, " - ", err)
			continue
		}
		workers[workerAccountName] = AccountAndAddress{
			acc:  workerAccount,
			addr: workerAddressToFund,
		}
	}
	return workers
}

func initializeNewWorkerAccount(
	m testCommon.TestConfig,
	topicId uint64,
	makeReport bool,
	workerAddresses *map[string]string, // pointer mutate map itself not a copy
) {
	// Generate new worker accounts
	workerAccountName := getWorkerAccountName(len(*workerAddresses), topicId)
	topicLog(topicId, "Initializing worker address: ", workerAccountName)
	workerAccount, err := m.Client.AccountRegistryGetByName(workerAccountName)
	if err != nil {
		topicLog(topicId, "Error getting worker address: ", workerAccountName, " - ", err)
		if makeReport {
			saveWorkerError(topicId, workerAccountName, err)
			saveTopicError(topicId, err)
		}
		return
	}
	workerAddress, err := workerAccount.Address(params.HumanCoinUnit)
	if err != nil {
		topicLog(topicId, "Error getting worker address: ", workerAccountName, " - ", err)
		if makeReport {
			saveWorkerError(topicId, workerAccountName, err)
			saveTopicError(topicId, err)
		}
		return
	}

	(*workerAddresses)[workerAccountName] = workerAddress
}

// register all the created workers for this iteration
func registerWorkersForIteration(
	m testCommon.TestConfig,
	topicId uint64,
	iteration int,
	workersPerIteration int,
	countWorkers int,
	maxWorkersPerTopic int,
	workers NameToAccountMap,
	makeReport bool,
) int {
	for j := 0; j < workersPerIteration && countWorkers < maxWorkersPerTopic; j++ {
		workerName := getWorkerAccountName(iteration*j, topicId)
		worker := workers[workerName]
		err := RegisterWorkerForTopic(m, worker.addr, worker.acc, topicId)
		if err != nil {
			topicLog(topicId, "Error registering worker address: ", worker.addr, " - ", err)
			if makeReport {
				saveWorkerError(topicId, workerName, err)
				saveTopicError(topicId, err)
			}
			return countWorkers
		}
		countWorkers++
	}
	return countWorkers
}

// pick a worker to upload a bundle, then try to insert the bundle
// if the bundle nonce is already fulfilled, realign the blockHeights and retry
// up to retry times
func generateInsertWorkerBundle(
	m testCommon.TestConfig,
	topic *emissionstypes.Topic,
	workers NameToAccountMap,
	blockHeightCurrent int64,
	retryTimes int,
	makeReport bool,
) (blockHeightEval int64, err error) {
	leaderWorkerAccountName, err := pickRandomKeyFromMap(workers)
	if err != nil {
		topicLog(topic.Id, "Error getting random worker address: ", err)
		return blockHeightCurrent, err
	}
	startWorker := time.Now()
	err = insertWorkerBulk(m, topic, leaderWorkerAccountName, workers, blockHeightCurrent)
	if err != nil {
		if strings.Contains(err.Error(), "nonce already fulfilled") {
			// realign blockHeights before retrying
			topic, err = getLastTopic(m.Ctx, m.Client.QueryEmissions(), topic.Id)
			if err == nil {
				blockHeightCurrent = topic.EpochLastEnded + topic.EpochLength
				blockHeightEval = blockHeightCurrent - topic.EpochLength
				topicLog(topic.Id, "Reset blockHeights to (", blockHeightCurrent, ",", blockHeightEval, ")")
			} else {
				topicLog(topic.Id, "Error getting topic!")
				if makeReport {
					saveTopicError(topic.Id, err)
				}
			}
		}
		return blockHeightEval, err
	} else {
		topicLog(topic.Id, "Inserted worker bulk, blockHeight: ", blockHeightCurrent, " with ", len(workers), " workers")
		elapsedBulk := time.Since(startWorker)
		topicLog(topic.Id, "Insert Worker ", blockHeightCurrent, " Elapsed time:", elapsedBulk)
	}
	return blockHeightCurrent, nil
}

// Inserts bulk inference and forecast data for a worker
func insertWorkerBulk(
	m testCommon.TestConfig,
	topic *emissionstypes.Topic,
	leaderWorkerAccountName string,
	workers map[string]AccountAndAddress,
	blockHeight int64,
) error {
	// Get Bundles
	workerDataBundles := make([]*emissionstypes.WorkerDataBundle, 0)
	for key := range workers {
		workerDataBundles = append(workerDataBundles,
			generateSingleWorkerBundle(m, topic.Id, blockHeight, key, workers))
	}
	leaderWorker := workers[leaderWorkerAccountName]
	return insertLeaderWorkerBulk(m, topic.Id, blockHeight, leaderWorkerAccountName, leaderWorker.addr, workerDataBundles)
}

// create inferences and forecasts for a worker
func generateSingleWorkerBundle(
	m testCommon.TestConfig,
	topicId uint64,
	blockHeight int64,
	workerAddressName string,
	workers map[string]AccountAndAddress,
) *emissionstypes.WorkerDataBundle {
	// Iterate workerAddresses to get the worker address, and generate as many forecasts as there are workers
	forecastElements := make([]*emissionstypes.ForecastElement, 0)
	for key := range workers {
		forecastElements = append(forecastElements, &emissionstypes.ForecastElement{
			Inferer: workers[key].addr,
			Value:   alloraMath.NewDecFromInt64(int64(rand.Intn(51) + 50)),
		})
	}
	infererAddress := workers[workerAddressName].addr
	infererValue := alloraMath.NewDecFromInt64(int64(rand.Intn(300) + 3000))

	// Create a MsgInsertBulkReputerPayload message
	workerDataBundle := &emissionstypes.WorkerDataBundle{
		Worker: infererAddress,
		InferenceForecastsBundle: &emissionstypes.InferenceForecastBundle{
			Inference: &emissionstypes.Inference{
				TopicId:     topicId,
				BlockHeight: blockHeight,
				Inferer:     infererAddress,
				Value:       infererValue,
			},
			Forecast: &emissionstypes.Forecast{
				TopicId:          topicId,
				BlockHeight:      blockHeight,
				Forecaster:       infererAddress,
				ForecastElements: forecastElements,
			},
		},
	}

	// Sign
	src := make([]byte, 0)
	src, err := workerDataBundle.InferenceForecastsBundle.XXX_Marshal(src, true)
	require.NoError(m.T, err, "Marshall reputer value bundle should not return an error")

	sig, pubKey, err := m.Client.Context().Keyring.Sign(workerAddressName, src, signing.SignMode_SIGN_MODE_DIRECT)
	require.NoError(m.T, err, "Sign should not return an error")
	workerPublicKeyBytes := pubKey.Bytes()
	workerDataBundle.InferencesForecastsBundleSignature = sig
	workerDataBundle.Pubkey = hex.EncodeToString(workerPublicKeyBytes)

	return workerDataBundle
}

// Inserts worker bulk, given a topic, blockHeight, and leader worker address (which should exist in the keyring)
func insertLeaderWorkerBulk(
	m testCommon.TestConfig,
	topicId uint64,
	blockHeight int64,
	leaderWorkerAccountName, leaderWorkerAddress string,
	WorkerDataBundles []*emissionstypes.WorkerDataBundle) error {

	nonce := emissionstypes.Nonce{BlockHeight: blockHeight}

	// Create a MsgInsertBulkReputerPayload message
	workerMsg := &emissionstypes.MsgInsertBulkWorkerPayload{
		Sender:            leaderWorkerAddress,
		Nonce:             &nonce,
		TopicId:           topicId,
		WorkerDataBundles: WorkerDataBundles,
	}
	// serialize workerMsg to json and print
	LeaderAcc, err := m.Client.AccountRegistryGetByName(leaderWorkerAccountName)
	if err != nil {
		fmt.Println("Error getting leader worker account: ", leaderWorkerAccountName, " - ", err)
		return err
	}
	txResp, err := m.Client.BroadcastTx(m.Ctx, LeaderAcc, workerMsg)
	if err != nil {
		fmt.Println("Error broadcasting worker bulk: ", err)
		return err
	}
	_, err = m.Client.WaitForTx(m.Ctx, txResp.TxHash)
	if err != nil {
		fmt.Println("Error waiting for worker bulk: ", err)
		return err
	}
	return nil
}

func checkWorkersReceivedRewards(
	m testCommon.TestConfig,
	topicId uint64,
	workers NameToAccountMap,
	countWorkers int,
	maxIterations int,
	makeReport bool,
) (rewardedWorkersCount int, err error) {
	rewardedWorkersCount = 0
	err = nil
	for workerIndex := 0; workerIndex < countWorkers; workerIndex++ {
		workerName := getWorkerAccountName(workerIndex, topicId)
		balance, err := getAccountBalance(
			m.Ctx,
			m.Client.QueryBank(),
			workers[workerName].addr,
		)
		if err != nil {
			topicLog(topicId, "Error getting worker balance for worker: ", workerName, err)
			if maxIterations > 20 && workerIndex < 10 {
				topicLog(topicId, "ERROR: Worker", workerName, "has insufficient stake:", balance)
			}
			if makeReport {
				saveWorkerError(topicId, workerName, err)
				saveTopicError(topicId, err)
			}
		} else {
			if balance.Amount.LTE(cosmosMath.NewInt(initialWorkerReputerFundAmount)) {
				topicLog(topicId, "Worker ", workerName, " balance is not greater than initial amount: ", balance.Amount.String())
				if makeReport {
					saveWorkerError(topicId, workerName, fmt.Errorf("Balance Not Greater"))
					saveTopicError(topicId, fmt.Errorf("Balance Not Greater"))
				}
			} else {
				topicLog(topicId, "Worker ", workerName, " balance: ", balance.Amount.String())
				rewardedWorkersCount += 1
			}
		}
	}
	return rewardedWorkersCount, err
}