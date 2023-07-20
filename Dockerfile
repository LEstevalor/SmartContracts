FROM hyperledger/fabric-ccenv:2.3.3

# 将智能合约文件（chaincode）复制到容器中的/root/go/src/chaincode/chaincode目录下
COPY chaincode /root/go/src/chaincode/chaincode

# 将工作目录更改为智能合约目录
WORKDIR /root/go/src/chaincode

# 设置核心链码ID名称和对等地址环境变量
ENV CORE_CHAINCODE_ID_NAME=mychain:1.0
ENV CORE_PEER_ADDRESS=peer0.org1.example.com:7051

# 设置容器启动时执行的命令
CMD ["./mychain/chaincode"]
