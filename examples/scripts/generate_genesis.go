package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	"znn-sdk-go/wallet"

	"github.com/zenon-network/go-zenon/chain/genesis"
	"github.com/zenon-network/go-zenon/common"
	"github.com/zenon-network/go-zenon/common/types"
	"github.com/zenon-network/go-zenon/vm/embedded/definition"
)

// args: genesis.template.json producer.json "producer_password" total_pillars genesis.json
func main() {
	args := os.Args
	keyStore, err := wallet.ReadKeyFile(args[2], args[3], "")
	common.DealWithErr(err)
	pillars := make([]types.Address, 0)
	numPillars, err := strconv.Atoi(args[4])
	common.DealWithErr(err)
	for i := 1; i <= numPillars; i++ {
		_, file, err := keyStore.DeriveForIndexPath(uint32(i))
		if err != nil {
			continue
		}
		fmt.Println(file.Address)
		pillars = append(pillars, file.Address)
	}

	jsonFile, err := os.Open(args[1])
	common.DealWithErr(err)
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var config genesis.GenesisConfig
	err = json.Unmarshal(byteValue, &config)
	common.DealWithErr(err)

	pillar := config.PillarConfig.Pillars[0]

	config.PillarConfig.Pillars = make([]*definition.PillarInfo, 0)
	config.PillarConfig.Delegations = make([]*definition.DelegationInfo, 0)
	for index, address := range pillars {
		newPillar := &definition.PillarInfo{
			Name:                         fmt.Sprintf("PILLAR_%d", index+1),
			StakeAddress:                 address,
			BlockProducingAddress:        address,
			RewardWithdrawAddress:        address,
			Amount:                       pillar.Amount,
			RegistrationTime:             pillar.RegistrationTime,
			RevokeTime:                   pillar.RevokeTime,
			GiveBlockRewardPercentage:    pillar.GiveBlockRewardPercentage,
			GiveDelegateRewardPercentage: pillar.GiveDelegateRewardPercentage,
			PillarType:                   pillar.PillarType,
		}
		config.PillarConfig.Pillars = append(config.PillarConfig.Pillars, newPillar)

		newDelegation := &definition.DelegationInfo{
			Name:   fmt.Sprintf("PILLAR_%d", index+1),
			Backer: address,
		}
		config.PillarConfig.Delegations = append(config.PillarConfig.Delegations, newDelegation)
	}

	for _, address := range pillars {
		balance := make(map[types.ZenonTokenStandard]*big.Int, 0)
		balance[types.ZnnTokenStandard] = big.NewInt(10000000000000)
		balance[types.QsrTokenStandard] = big.NewInt(100000000000000)
		newBlock := &genesis.GenesisBlockConfig{
			Address:     address,
			BalanceList: balance,
		}
		config.GenesisBlocks.Blocks = append(config.GenesisBlocks.Blocks, newBlock)
	}

	for _, address := range pillars {
		fusion := &definition.FusionInfo{
			Owner:            address,
			Id:               config.PlasmaConfig.Fusions[0].Id,
			Amount:           big.NewInt(1000000000000),
			ExpirationHeight: 1,
			Beneficiary:      address,
		}
		config.PlasmaConfig.Fusions = append(config.PlasmaConfig.Fusions, fusion)
	}

	qsrTotalAmount := big.NewInt(0)

	for _, fusion := range config.PlasmaConfig.Fusions {
		if fusion == nil {
			return
		}
		qsrTotalAmount.Add(qsrTotalAmount, fusion.Amount)
	}

	for _, block := range config.GenesisBlocks.Blocks {
		if block.Address != types.PlasmaContract {
			continue
		}
		block.BalanceList[types.QsrTokenStandard] = qsrTotalAmount
	}

	znnTotalAmount := big.NewInt(0)

	for _, el := range config.PillarConfig.Pillars {
		znnTotalAmount.Add(znnTotalAmount, el.Amount)
	}

	for _, block := range config.GenesisBlocks.Blocks {
		if block.Address != types.PillarContract {
			continue
		}
		block.BalanceList[types.ZnnTokenStandard] = znnTotalAmount
	}

	given := make(map[types.ZenonTokenStandard]*big.Int)
	for _, block := range config.GenesisBlocks.Blocks {
		for zts, amount := range block.BalanceList {
			total, ok := given[zts]
			if ok == false {
				given[zts] = new(big.Int).Set(amount)
			} else {
				total.Add(total, amount)
			}
		}
	}

	for _, token := range config.TokenConfig.Tokens {
		total, ok := given[token.TokenStandard]
		if ok == false {
			return
		}
		token.TotalSupply = total
	}

	err = genesis.CheckGenesis(&config)
	common.DealWithErr(err)

	file, _ := json.MarshalIndent(config, "", "\t")
	_ = ioutil.WriteFile(args[5], file, 0644)
}
