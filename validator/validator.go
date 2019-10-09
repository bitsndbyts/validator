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

func main() {
	rpcClient := client.NewHTTP("localhost:26657", "/websocket")
	err := rpcClient.OnStart()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to localhost:26657")

	ctx := context.Background()
	block, err := rpcClient.Subscribe(ctx, "subscribe", "tm.event='NewBlock'")
	if err != nil {
		panic(err)
	}

	_block := tmTypes.EventDataNewBlock{}

	for event := range block {
		bz, err := json.Marshal(event.Data)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(bz, &_block)
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
}
