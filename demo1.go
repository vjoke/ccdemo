/**
 * @author: hsiung
 * @date: 2016/12/3
 * @desc: a simple demo
 */

package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const (
	TICKET_PREFIX = "ticket-"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// Ticket ticket
type Ticket struct {
	TxID   string // use transaction id
	Number string
}

// Init - called when the chaincode is deployed to the chain.
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Chain code deployed successfully\n")
	return nil, nil
}

// Invoke - called to handle invoke transactions
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "putTicket" {
		return putTicket(stub, args)
	}

	return nil, errors.New("function not supported")
}

// Query - callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "getTicket" {
		return getTicket(stub, args)
	}

	return nil, errors.New("function not supported")
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Help routines
func putTicket(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	// error checking
	var number = args[0]
	check, err := strconv.ParseBool(args[1])
	if err != nil {
		return nil, errors.New("Invalid parameter")
	}

	if check && len(number) != 7 {
		return nil, errors.New("invalid ticket number")
	}
	// FIXME: we use transaction as the id of ticket
	var ticket = Ticket{TxID: stub.GetTxID(), Number: number}
	return writeTicket(stub, ticket)
}

func writeTicket(stub shim.ChaincodeStubInterface, ticket Ticket) ([]byte, error) {
	ticketBytes, err := json.Marshal(&ticket)
	if err != nil {
		return nil, err
	}

	err = stub.PutState(TICKET_PREFIX+ticket.TxID, ticketBytes)
	if err != nil {
		return nil, errors.New("PutState error" + err.Error())
	}

	return ticketBytes, nil
}

func getTicket(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	ticketBytes, err := stub.GetState(TICKET_PREFIX + args[0])
	if err != nil {
		return nil, errors.New("Error retrieving ticket")
	}

	var ticket Ticket
	err = json.Unmarshal(ticketBytes, &ticket)
	if err != nil {
		return nil, errors.New("Error unmarshalling ticket")
	}

	return ticketBytes, nil
}
