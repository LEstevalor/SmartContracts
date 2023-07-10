package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type SimpleContract struct {
}

func (s *SimpleContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (s *SimpleContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "submitTransaction" {
		return s.submitTransaction(stub, args)
	}
	return shim.Error("Invalid function name.")
}

func (s *SimpleContract) submitTransaction(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2.")
	}

	// 验证交易时间戳，防止女巫攻击
	txTimestamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error("Error getting transaction timestamp: " + err.Error())
	}
	txTime := time.Unix(txTimestamp.Seconds, 0)
	if time.Since(txTime) > 1*time.Minute {
		return shim.Error("Transaction is too old.")
	}

	// 验证输入参数，防止日蚀攻击
	key := args[0]
	value, err := strconv.Atoi(args[1])
	if err != nil || value < 0 {
		return shim.Error("Invalid value. Value should be a non-negative integer.")
	}

	err = stub.PutState(key, []byte(strconv.Itoa(value)))
	if err != nil {
		return shim.Error("Failed to update state: " + err.Error())
	}

	return shim.Success(nil)
}

func main() {

	err := shim.Start(new(SimpleContract))
	if err != nil {
		fmt.Printf("Error starting SimpleContract chaincode: %s", err)
	}
}
