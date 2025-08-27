package utxo

import (
	"fmt"
)

type UTXO struct {
	TransactionID string `json:"transactionID"`
	OutputIndex   uint   `json:"outputIndex"`
	Address       string `json:"address"`
	Amount        uint64 `json:"amount"`
}

type UTXOSet struct {
	UTXOs []*UTXO `json:"UTXOs"`
}

func NewUTXOSet() *UTXOSet {
	return &UTXOSet{}
}

func (us *UTXOSet) AddUTXO(u *UTXO) {
	us.UTXOs = append(us.UTXOs, u)
}

func (us *UTXOSet) RemoveUTXO(u *UTXO) {
	for i, utxo := range us.UTXOs {
		if utxo.TransactionID == u.TransactionID && utxo.OutputIndex == u.OutputIndex {
			us.UTXOs = append(us.UTXOs[:i], us.UTXOs[i+1:]...)
			break
		}
	}
}

func (us *UTXOSet) RemoveUTXOByID(transactionID string, outputIndex uint) {
	for i, utxo := range us.UTXOs {
		if utxo.TransactionID == transactionID && utxo.OutputIndex == outputIndex {
			us.UTXOs = append(us.UTXOs[:i], us.UTXOs[i+1:]...)
			break
		}
	}
}

func (us *UTXOSet) UTXOExists(transactionID string, outputIndex uint) bool {
	for _, u := range us.UTXOs {
		if u.TransactionID == transactionID && u.OutputIndex == outputIndex {
			return true
		}
	}

	return false
}

func (us *UTXOSet) GetUTXO(transactionID string, outputIndex uint) *UTXO {
	for _, u := range us.UTXOs {
		if u.TransactionID == transactionID && u.OutputIndex == outputIndex {
			return u
		}
	}

	return nil
}

func (us *UTXOSet) GetUTXOsByAddress(address string) []*UTXO {
	utxos := make([]*UTXO, 0)

	for _, u := range us.UTXOs {
		if u.Address == address {
			utxos = append(utxos, u)
		}
	}

	return utxos
}

func (us *UTXOSet) GetAddressBalance(address string) uint64 {
	var balance uint64

	for _, u := range us.GetUTXOsByAddress(address) {
		balance += u.Amount
	}

	return balance
}

// Input: address + desired qty | Output: UTXO list that sum >= desired qty
func (us *UTXOSet) FindSpendableUTXOsForAddress(address string, desiredQty uint64) ([]*UTXO, error) {
	var currQty uint64 = 0
	utxos := make([]*UTXO, 0)

	for _, u := range us.GetUTXOsByAddress(address) {
		if currQty >= desiredQty {
			break
		}

		currQty += u.Amount
		utxos = append(utxos, u)
	}

	if currQty < desiredQty {
		return nil, fmt.Errorf("insuficient funds")
	}

	return utxos, nil
}

func (us *UTXOSet) HasSufficientFunds(address string, amount uint64) bool {
	spendableUTXOs, err := us.FindSpendableUTXOsForAddress(address, amount)
	if err != nil {
		return false
	}

	if len(spendableUTXOs) > 0 {
		return true
	}

	return false
}
