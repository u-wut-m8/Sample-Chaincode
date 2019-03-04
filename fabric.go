package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type SampleChaincode struct{}

func (t *SampleChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	//Get the arguments from transaction proposal
	args := stub.GetStringArgs()
	if len(args) != 2 {
		return shim.Error("Incorrect arguments. Expecting a key and a value.")
	}
	//We store the key and the value on the ledger
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
	}
	return shim.Success(nil)
}

func (t *SampleChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	//Extract the function and arguments from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()
	var result string
	var err error
	if fn == "set" {
		result, err = set(stub, args)
	} else { //Assume 'get' even when fn is nil
		result, err = get(stub, args)
	}
	if err != nil { //Failed to get function and/or arguments from transaction proposal
		return shim.Error(err.Error())
	}
	//Return the result as success payload
	return shim.Success([]byte(result))
}

func set(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return args[1], nil
}

func get(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key.")
	}
	value, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	if value == nil {
		return "", fmt.Errorf("Asset not found: %s", args[0])
	}
	return string(value), nil
}

func main() {
	err := shim.Start(new(SampleChaincode))
	if err != nil {
		fmt.Println("Could not start SampleChaincode")
	} else {
		fmt.Println("SampleChaincode successfully started.")
	}
}
