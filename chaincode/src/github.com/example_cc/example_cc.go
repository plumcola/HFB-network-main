/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"strconv"
	"errors"	
	log "github.com/sirupsen/logrus" 

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer")

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) (pb.Response,error) {
	//fmt.Println("ex02 Init!!!")
	log.Infof("[%s][modbuschannel][example_cc][Init] ex02 Init",uuidgen())
	_, args := stub.GetFunctionAndParameters()
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var err error

	if len(args) != 4 {
		//return shim.Error("Incorrect argument numbers. Expecting 4")
		log.Errorf("[%s][modbuschannel][example_cc][valueIssuer] Incorrect argument numbers. Expecting 4",uuidgen())
		return shim.Error("") , errors.New(ERRORWrongNumberArgs)
	}

	// Initialize the chaincode
	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		//return shim.Error("Expecting integer value for asset holding")
		log.Errorf("[%s][modbuschannel][example_cc][valueIssuer] Expecting integer value for asset holding",uuidgen())
		return shim.Error("") , errors.New(ERRORParsingData)

	}
	B = args[2]
	Bval, err = strconv.Atoi(args[3])
	if err != nil {
		//return shim.Error("Expecting integer value for asset holding")
		log.Errorf("[%s][modbuschannel][example_cc][valueIssuer] Expecting integer value for asset holding",uuidgen())
		return shim.Error("") , errors.New(ERRORParsingData)
	}
	//fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)
	log.Infof("[%s][modbuschannel][example_cc][Init] Initialize the chaincode with Aval = %d, Bval = %d",uuidgen(), Aval, Bval)

	// Write the state to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		//return shim.Error(err.Error())
		log.Errorf("[%s][modbuschannel][example_cc][stateIssuer] Error in writing the state to the ledger",uuidgen())
		return shim.Error("") , errors.New(ERRORPutState)
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		//return shim.Error(err.Error())
		log.Errorf("[%s][modbuschannel][example_cc][stateIssuer] Error in writing the state to the ledger",uuidgen())
		return shim.Error("") , errors.New(ERRORPutState)
	}

	log.Infof("[%s][modbuschannel][example_cc][PutState] Succeed to write the state to the ledger",uuidgen())
	return shim.Success(nil) , errors.New("") 
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) (pb.Response,error) {
	//fmt.Println("ex02 Invoke")
	log.Infof("[%s][modbuschannel][example_cc][Invoke] ex02 Invoke",uuidgen())

	function, args := stub.GetFunctionAndParameters()
	if function == "invoke" {
		// Make payment of X units from A to B
		return t.invoke(stub, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	} else if function == "query" {
		// the old "Query" is now implemtned in invoke
		return t.query(stub, args)
	}

	//return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
	log.Errorf("[%s][modbuschannel][example_cc][invokeIssuer] Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"",uuidgen())
	return shim.Error("") , errors.New(ERRORServiceNotExists)
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) (pb.Response,error) {
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var X int          // Transaction value
	var err error

	log.Infof("[%s][modbuschannel][example_cc][Invoke] ex02 invoke",uuidgen())

	if len(args) != 3 {
		//return shim.Error("Incorrect number of arguments. Expecting 3")
		log.Errorf("[%s][modbuschannel][example_cc][valueIssuer] Incorrect number of arguments. Expecting 3",uuidgen())
		return shim.Error("") , errors.New(ERRORWrongNumberArgs)
	}

	A = args[0]
	B = args[1]

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		//return shim.Error("Failed to get state")
		log.Errorf("[%s][modbuschannel][example_cc][stateIssuer] Failed to get state",uuidgen())
		return shim.Error("") , errors.New(ERRORGetState)
	}
	if Avalbytes == nil {
		//return shim.Error("Entity not found")	
		log.Errorf("[%s][modbuschannel][example_cc][idIssuer] Entity not found",uuidgen())	
		return shim.Error("") , errors.New(EERRORnotID)	
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		//return shim.Error("Failed to get state")
		log.Errorf("[%s][modbuschannel][example_cc][stateIssuer] Failed to get state",uuidgen())
		return shim.Error("") , errors.New(ERRORGetState)
	}
	if Bvalbytes == nil {
		//return shim.Error("Entity not found")
		log.Errorf("[%s][modbuschannel][example_cc][idIssuer] Entity not found",uuidgen())	
		return shim.Error("") , errors.New(EERRORnotID)
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	// Perform the execution
	X, err = strconv.Atoi(args[2])
	if err != nil {
		//return shim.Error("Invalid transaction amount, expecting a integer value")
		log.Errorf("[%s][modbuschannel][example_cc][valueIssuer] Invalid transaction amount, expecting a integer value",uuidgen())
		return shim.Error("") , errors.New(ERRORParsingData)	
	}
	Aval = Aval - X
	Bval = Bval + X
	//fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)
	log.Infof("[%s][modbuschannel][example_cc][Transaction] Aval = %d, Bval = %d after performing the transaction",uuidgen(), Aval, Bval)	

	// Write the state back to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		//return shim.Error(err.Error())
		log.Errorf("[%s][modbuschannel][example_cc][stateIssuer] Failed to write the state back to the ledger",uuidgen())
		return shim.Error("") , errors.New(ERRORPutState)	
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		//return shim.Error(err.Error())
		log.Errorf("[%s][modbuschannel][example_cc][stateIssuer] Failed to write the state back to the ledger",uuidgen())
		return shim.Error("") , errors.New(ERRORPutState)
	}

	payloadAsBytes := []byte(strconv.Itoa(Bval))	
	log.Infof("[%s][modbuschannel][example_cc][Transaction] Transaction makes payment of X units from A to B",uuidgen())
	return shim.Success(payloadAsBytes) , errors.New("")
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) (pb.Response,error){
	if len(args) != 1 {
		//return shim.Error("Incorrect number of arguments. Expecting 1")
		log.Errorf("[%s][modbuschannel][example_cc][valueIssuer] Incorrect number of arguments. Expecting 1",uuidgen())
		return shim.Error("") , errors.New(ERRORWrongNumberArgs)
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		//return shim.Error("Failed to delete state")
		log.Errorf("[%s][modbuschannel][example_cc][stateIssuer] Failed to delete state",uuidgen())
		return shim.Error("") , errors.New(ERRORDelState)
	}

	log.Infof("[%s][modbuschannel][example_cc][DelState] Succeed to delete an entity from state",uuidgen())
	return shim.Success(nil), errors.New("")
	
}

// query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) (pb.Response,error){
	var A string // Entities
	var err error

	if len(args) != 1 {
		//return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
		log.Errorf("[%s][modbuschannel][example_cc][valueIssuer] Incorrect number of arguments. Expecting name of the person to query",uuidgen())
		return shim.Error("") , errors.New(ERRORWrongNumberArgs)
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		//return shim.Error(jsonResp)
		log.Errorf("[%s][modbuschannel][example_cc][stateIssuer] %s",uuidgen(),jsonResp)
		return shim.Error("") , errors.New(ERRORGetState)	
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		//return shim.Error(jsonResp)
		log.Errorf("[%s][modbuschannel][example_cc][valueIssuer] %s",uuidgen(),jsonResp)	
		return shim.Error("") , errors.New(ERRORParsingData)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	//fmt.Printf("Query Response:%s\n", jsonResp)
	log.Infof("[%s][modbuschannel][example_cc][Query] Query Response: %s",uuidgen(),jsonResp)
	return shim.Success(Avalbytes) , errors.New("")
}

func main() {
	customFormatter := new(log.TextFormatter)
    	customFormatter.TimestampFormat = "2006-01-02T15:04:05Z"
    	log.SetFormatter(customFormatter)
   	customFormatter.FullTimestamp = true

	//pbr , _ = new(SimpleChaincode) 
	//err := shim.Start(pbr)

	err := shim.Start(new(SimpleChaincode))

	if err != nil {
		//fmt.Printf("Error starting Simple chaincode: %s", err)
		log.Errorf("[%s][modbuschannel][example_cc][Init] Error starting Simple chaincode: %s",uuidgen(),err)
	} else {
		log.Infof("[%s][modbuschannel][example_cc][Init] Succeed to start Simple chaincode: %s",uuidgen(),err)
	}

}
