# How To Run

1. Clone the repository
2. Run a `walletd` node in full index mode (`walletd -index.mode full` or change index mode in walletd.yml)
	2.a The API address and password are hardcoded, you may need to change them for this script to work
3. Mine some utxos into the hardcoded address (`addr:00069004b64a79f7bfa3698f0c3cca13b2fc4c1054b2b3b6c58bc5bd2095b6a65811d503b9ba`)
	4.a `walletd mine --addr addr:00069004b64a79f7bfa3698f0c3cca13b2fc4c1054b2b3b6c58bc5bd2095b6a65811d503b9ba -n 10`
4. Mine until the v2 require height and the utxos mature
5. Run the script `go run main.go`. A transaction will be broadcast that sends one UTXO from the address to the void.

You can make changes to `main.go` to test different spend policies. 

The transaction that is sent is whatever utxo is returned first. That's usually the oldest utxo, but could be different. Production scripts would need to select UTXOs, calculate miner fees, and add a change output.