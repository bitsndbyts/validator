package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"

	"github.com/tendermint/tendermint/rpc/client"
	tmTypes "github.com/tendermint/tendermint/types"
)

type Account struct {
	Name    string
	Address string
}

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

	var accounts []Account
	for i := 0; i < 2; i++ {
		cmdCreateAccount := fmt.Sprintf("echo rgukt123 | ic-cli keys add accountl%d", i)
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
			cmd = exec.Command("sh", "-c", cmdSend)

			_, err = cmd.Output()
			if err != nil {
				panic(err)
			}

			fmt.Println("Random send transaction completed\n")
		}
	}
}
