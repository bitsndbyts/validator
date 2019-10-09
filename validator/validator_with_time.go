package main

import (
	"fmt"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/tendermint/tendermint/rpc/client"
)

func main() {
	rpcClient := client.NewHTTP("localhost:26657", "/websocket")
	err := rpcClient.OnStart()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to localhost:26657")
	for {
		time.Sleep(time.Minute * 30)

		cmdStr := fmt.Sprintf("echo 1234567890 | ic-cli tx distribution withdraw-rewards" +
			" interchangevaloper1yry29xzpp8yukzh9x0k5579t25vwlrs24scm67 --from interchange-4" +
			" --chain-id AlphaNet-1 --broadcast-mode block --yes")

		cmd := exec.Command("sh", "-c", cmdStr)

		_, err = cmd.Output()
		if err != nil {
			panic(err)
		}

		fmt.Println("***** Withdraw commission completed *****")
		time.Sleep(time.Second * 5)

		cmdQueryAccount := fmt.Sprintf("ic-cli query account interchange1yry29xzpp8yukzh9x0k5579t25vwlrs2ldemr5" +
			" --chain-id AlphaNet-1")
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
		cmdStr = fmt.Sprintf("echo 1234567890 | ic-cli tx staking delegate "+
			"interchangevaloper1yry29xzpp8yukzh9x0k5579t25vwlrs24scm67 %dintr --chain-id AlphaNet-1 "+
			"--from interchange-4 --broadcast-mode block --yes", randBal)
		cmd = exec.Command("sh", "-c", cmdStr)

		_, err = cmd.Output()
		if err != nil {
			panic(err)
		}

		fmt.Println("***** Delegation completed *****")
	}
}

