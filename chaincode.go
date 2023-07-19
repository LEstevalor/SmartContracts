package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type SimpleContract struct {
}

type Node struct {
	// 存储节点的ID和信誉分数
	ID         string
	Reputation int
}

func (s *SimpleContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	/**
	必须有init和invoke
	链码的初始化方法，在链码部署时被调用
	*/

	return shim.Success(nil)
}

func (s *SimpleContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	/**
	链码的主要入口点，用于处理链码的调用
	*/

	function, args := stub.GetFunctionAndParameters()
	// 检查调用的函数名是否为编写的有效函数名，如果是则调用并返回，否则返回Error
	if function == "submitTransaction" {
		return s.submitTransaction(stub, args)
	} else if function == "updateNodeReputation" {
		return s.updateNodeReputation(stub, args)
	}

	return shim.Error("Invalid function name.")
}

func (s *SimpleContract) submitTransaction(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	/**
	验证交易的时间戳、输入参数和信誉分数
	*/

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4.")
	}

	// 验证交易时间戳，防止女巫攻击
	// 从链码的shim.ChaincodeStubInterface获取交易的时间戳，并对其进行验证
	txTimestamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error("Error getting transaction timestamp: " + err.Error())
	}
	// 类型装换
	txTime := time.Unix(txTimestamp.Seconds, 0)
	// 计算从交易时间戳txTime到当前时间的时间间隔，若大于一分钟则声明无效
	if time.Since(txTime) > 1*time.Minute {
		return shim.Error("Transaction is too old.")
	}

	// 通过验证输入参数与信誉，防止日蚀攻击
	transactionReputation, err := strconv.Atoi(args[2])
	if err != nil || transactionReputation < 0 {
		return shim.Error("Invalid transaction reputation argument")
	}

	nodeId := args[3]
	nodeBytes, err := stub.GetState(nodeId)
	if err != nil {
		return shim.Error("Failed to get node reputation: " + err.Error())
	}
	var node Node
	err = json.Unmarshal(nodeBytes, &node)
	if err != nil {
		return shim.Error("Failed to unmarshal node: " + err.Error())
	}
	if transactionReputation > node.Reputation {
		return shim.Error("Invalid reputation")
	}

	// 给定的键值对存储到链码的状态数据库中
	key := args[0]
	value, err := strconv.Atoi(args[1])
	// 输入参数无效
	if err != nil || value < 0 {
		return shim.Error("Invalid value. Value should be a non-negative integer.")
	}

	err = stub.PutState(key, []byte(strconv.Itoa(value)))
	if err != nil {
		return shim.Error("Failed to update state: " + err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleContract) updateNodeReputation(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	/**
	创建和更新节点的信誉分数
	*/

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	nodeId := args[0]
	reputation, err := strconv.Atoi(args[1])
	if err != nil || reputation < 0 {
		return shim.Error("Invalid reputation value. Value should be a non-negative integer.")
	}

	node := Node{ID: nodeId, Reputation: reputation}
	nodeBytes, err := json.Marshal(node)
	if err != nil {
		return shim.Error("Failed to marshal node: " + err.Error())
	}

	err = stub.PutState(nodeId, nodeBytes)
	if err != nil {
		return shim.Error("Failed to update node reputation: " + err.Error())
	}

	return shim.Success(nil)
}

func main() {

	// 创建实例
	simpleContract := new(SimpleContract)
	// 启动链码
	err := shim.Start(simpleContract)
	if err != nil {
		fmt.Printf("Error starting SimpleContract chaincode: %s", err)
	}

}
