package main

import (
    "encoding/json"
    "fmt"

    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type CrossChainContract struct {
    contractapi.Contract
}

type Token struct {
    Owner  string `json:"owner"`
    Amount uint64 `json:"amount"`
}

// Mint tokens
func (c *CrossChainContract) Mint(ctx contractapi.TransactionContextInterface, ethereumAddress string, amount uint64) error {
    tokenKey := ethereumAddress
    tokenJSON, err := ctx.GetStub().GetState(tokenKey)
    if err != nil {
        return fmt.Errorf("failed to read token: %v", err)
    }

    var token Token
    if tokenJSON == nil {
        token = Token{
            Owner:  ethereumAddress,
            Amount: amount,
        }
    } else {
        err = json.Unmarshal(tokenJSON, &token)
        if err != nil {
            return fmt.Errorf("failed to unmarshal token: %v", err)
        }
        token.Amount += amount
    }

    tokenJSON, err = json.Marshal(token)
    if err != nil {
        return fmt.Errorf("failed to marshal token: %v", err)
    }

    return ctx.GetStub().PutState(tokenKey, tokenJSON)
}

// Burn tokens
func (c *CrossChainContract) Burn(ctx contractapi.TransactionContextInterface, ethereumAddress string, amount uint64) error {
    tokenKey := ethereumAddress
    tokenJSON, err := ctx.GetStub().GetState(tokenKey)
    if err != nil {
        return fmt.Errorf("failed to read token: %v", err)
    }
    if tokenJSON == nil {
        return fmt.Errorf("token does not exist")
    }

    var token Token
    err = json.Unmarshal(tokenJSON, &token)
    if err != nil {
        return fmt.Errorf("failed to unmarshal token: %v", err)
    }
    if token.Amount < amount {
        return fmt.Errorf("insufficient token balance")
    }

    token.Amount -= amount
    tokenJSON, err = json.Marshal(token)
    if err != nil {
        return fmt.Errorf("failed to marshal token: %v", err)
    }

    return ctx.GetStub().PutState(tokenKey, tokenJSON)
}

func main() {
    chaincode, err := contractapi.NewChaincode(new(CrossChainContract))
    if err != nil {
        panic(fmt.Sprintf("Error creating CrossChainContract chaincode: %v", err))
    }
    if err := chaincode.Start(); err != nil {
        panic(fmt.Sprintf("Error starting CrossChainContract chaincode: %v", err))
    }
}
