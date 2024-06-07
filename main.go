package main

import (
	"encoding/json"
	"log"

	"go.sia.tech/core/types"
	"go.sia.tech/walletd/api"
)

const (
	apiAddress = "http://localhost:9980/api"
	password   = "sia is cool"
)

var (
	sk = types.NewPrivateKeyFromSeed(make([]byte, 32))
)

func main() {
	client := api.NewClient(apiAddress, password)

	sp := types.PolicyPublicKey(sk.PublicKey()) // change me
	addr := sp.Address()

	log.Printf("Mining Address: %s", addr)

	balance, err := client.AddressBalance(addr)
	if err != nil {
		panic(err)
	}
	log.Printf("Balance: %s", balance.Siacoins)

	state, err := client.ConsensusTipState()
	if err != nil {
		panic(err)
	}

	utxos, err := client.AddressSiacoinOutputs(addr, 0, 100)
	if err != nil {
		panic(err)
	}
	filtered := utxos[:0]
	for _, utxo := range utxos {
		if utxo.MaturityHeight > state.Index.Height {
			continue
		}
		filtered = append(filtered, utxo)
	}
	utxos = filtered

	if len(utxos) == 0 {
		log.Println("No spendable UTXOs")
		return
	}

	/*fee, err := client.TxpoolFee()
	if err != nil {
		panic(err)
	}

	// Could estimate the weight before fundind, but it's
	// simpler to use a const and fees are, generally, negligible.
	const txnWeight = 1000*/

	// Create the transaction
	// minerFee := fee.Mul64(txnWeight)
	// sendAmount := types.Siacoins(100)
	// changeAmount := utxos[0].SiacoinOutput.Value.Sub(sendAmount).Sub(minerFee) // panics on underflow
	txn := types.V2Transaction{
		// MinerFee: minerFee,
		SiacoinInputs: []types.V2SiacoinInput{
			{
				Parent: utxos[0],
			},
		},
		SiacoinOutputs: []types.SiacoinOutput{
			{Address: types.VoidAddress, Value: utxos[0].SiacoinOutput.Value}, // send to the void
			// {Address: addr, Value: changeAmount},                              // send change back
		},
	}

	// get the transaction hash
	sigHash := state.InputSigHash(txn)
	// sign the transaction
	txn.SiacoinInputs[0].SatisfiedPolicy = types.SatisfiedPolicy{
		Policy:     sp,
		Signatures: []types.Signature{sk.SignHash(sigHash)},
	}

	req := api.TxpoolBroadcastRequest{
		V2Transactions: []types.V2Transaction{txn},
	}
	jsonBuf, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		panic(err)
	}
	log.Println(string(jsonBuf))

	if err := client.TxpoolBroadcast(nil, []types.V2Transaction{txn}); err != nil {
		panic(err)
	}
	log.Println("Transaction sent", txn.ID())
}
