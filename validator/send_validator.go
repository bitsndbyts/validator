package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/tendermint/tendermint/rpc/client"
	tmTypes "github.com/tendermint/tendermint/types"
)

type Account struct {
	Name    string
	Address string
}

func main() {
	rpcClient1 := client.NewHTTP("localhost:26657", "/websocket")
	err := rpcClient1.OnStart()
	if err != nil {
		panic(err)
	}

	fmt.Println("client1 connected to localhost:26657")

	rpcClient2 := client.NewHTTP("localhost:26657", "/websocket")
	err = rpcClient2.OnStart()
	if err != nil {
		panic(err)
	}

	fmt.Println("client2 connected to localhost:26657")

	var accounts []Account
	for i := 0; i < 2; i++ {
		cmdCreateAccount := fmt.Sprintf("echo rgukt123 | ic-cli keys add accountz%d", i)
		cmd := exec.Command("sh", "-c", cmdCreateAccount)

		res, err := cmd.Output()
		if err != nil {
			panic(err)
		}

		resSlice := strings.Split(string(res), ":")
		nameSlice := strings.Split(resSlice[1], "\n")
		addressSlice := strings.Split(resSlice[3], "\n")

		accounts = append(accounts, Account{Name: nameSlice[0], Address: addressSlice[0]})

		cmdSend := fmt.Sprintf("echo rgukt123 | ic-cli tx send "+
			"interchange1jn936et0r04rlqdt5jjuh46vyt024xe9kwx7p4 %s 10000stake"+
			" --chain-id ic --broadcast-mode block --yes", addressSlice[0])

		cmd = exec.Command("sh", "-c", cmdSend)

		_, err = cmd.Output()
		if err != nil {
			panic(err)
		}

		fmt.Println("Send transaction completed\n")
	}

	ctx1 := context.Background()
	ctx2 := context.Background()
	block, err := rpcClient1.Subscribe(ctx1, "subscribe", "tm.event='NewBlock'")
	if err != nil {
		panic(err)
	}

	blk, err := rpcClient2.Subscribe(ctx2, "subscribe", "tm.event='NewBlock'")
	if err != nil {
		panic(err)
	}

	_block := tmTypes.EventDataNewBlock{}

	go func() {
		for event := range block {
			bz, err := json.Marshal(event.Data)
			if err != nil {
				panic(err)
			}

			err = json.Unmarshal(bz, &_block)
			if err != nil {
				panic(err)
			}

			if _block.Block.Header.Height%1 == 0 {
				randFromAccount := rand.Intn(2)
				randToAccount := rand.Intn(2)

				cmdQueryAccount := fmt.Sprintf("ic-cli query account "+
					"%s --chain-id ic", accounts[randFromAccount].Address)
				cmd := exec.Command("sh", "-c", cmdQueryAccount)

				res, err := cmd.Output()
				if err != nil {
					panic(err)
				}

				balSlice := strings.Split(string(res), "\"")
				bal, err := strconv.Atoi(balSlice[1])
				if err != nil {
					panic(err)
				}

				randBal := rand.Intn(bal)
				cmdSend := fmt.Sprintf("echo rgukt123 | ic-cli  tx send %s %s %dstake"+
					" --chain-id ic --broadcast-mode block --yes",
					accounts[randFromAccount].Address, accounts[randToAccount].Address, randBal)
				fmt.Println(cmdSend)
				cmd = exec.Command("sh", "-c", cmdSend)

				_, err = cmd.Output()
				if err != nil {
					fmt.Println(err)
					panic(err)
				}

				fmt.Println("Random send transaction completed\n")
			}
		}
	}()

	_blk := tmTypes.EventDataNewBlock{}

	go func() {
		for event := range blk {
			bz, err := json.Marshal(event.Data)
			if err != nil {
				panic(err)
			}

			err = json.Unmarshal(bz, &_blk)
			if err != nil {
				panic(err)
			}

			if _block.Block.Header.Height%5 == 0 {
				cmdStr := fmt.Sprintf("echo rgukt123 | ic-cli tx distribution withdraw-rewards" +
					" interchangevaloper1jn936et0r04rlqdt5jjuh46vyt024xe9un87cl --from genesis" +
					" --chain-id ic --broadcast-mode block --yes")

				fmt.Println(cmdStr)
				cmd := exec.Command("sh", "-c", cmdStr)

				_, err := cmd.Output()
				if err != nil {
					panic(err)
				}

				fmt.Println("***** Withdraw commission completed *****")
				time.Sleep(time.Second * 5)

				cmdQueryAccount := fmt.Sprintf("ic-cli query account interchange1jn936et0r04rlqdt5jjuh46vyt024xe9kwx7p4" +
					" --chain-id ic")
				cmd = exec.Command("sh", "-c", cmdQueryAccount)

				res, err := cmd.Output()
				if err != nil {
					panic(err)
				}

				balSlice := strings.Split(string(res), "\"")
				bal, err := strconv.Atoi(balSlice[1])
				if err != nil {
					panic(err)
				}

				randBal := rand.Intn(bal)
				cmdStr = fmt.Sprintf("echo rgukt123 | ic-cli tx staking delegate "+
					"interchangevaloper1jn936et0r04rlqdt5jjuh46vyt024xe9un87cl %dstake --chain-id ic "+
					"--from genesis --broadcast-mode block --yes", randBal)
				cmd = exec.Command("sh", "-c", cmdStr)

				_, err = cmd.Output()
				if err != nil {
					panic(err)
				}

				fmt.Println("***** Delegation completed *****")
			}
		}
	}()
	select {}
}
