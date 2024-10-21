package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

type Block struct {
	index        int
	previousHash string
	timestamp    int64
	transactions []*Transaction
	thisHash     string
}

type Transaction struct {
	TransactionID              string
	SenderBlockchainAddress    string
	RecipientBlockchainAddress string
	Value                      float32
}

type Blockchain struct {
	blocks          []*Block
	TransactionPool []*Transaction
}

func newTransaction(sender string, recipent string, Value float32) *Transaction {

	t := new(Transaction)
	t.SenderBlockchainAddress = sender
	t.RecipientBlockchainAddress = recipent
	t.Value = Value

	t.TransactionID = calHashTrans(t)

	return t
}

func calHashTrans(t *Transaction) string {
	data := fmt.Sprintf("%s%s%f", t.SenderBlockchainAddress, t.RecipientBlockchainAddress, t.Value)
	hash := sha256.Sum256([]byte(data))
	hashString := hex.EncodeToString(hash[:])
	return hashString
}

func (bc *Blockchain) AddTransaction(sender string, recipient string, Value float32) {
	transaction := newTransaction(sender, recipient, Value)

	bc.TransactionPool = append(bc.TransactionPool, transaction)
}

func calculateHash(b *Block) string {

	var transactionsData string

	for _, t := range b.transactions {
		transactionsData += fmt.Sprintf("%s%s%s%f", t.TransactionID, t.SenderBlockchainAddress, t.RecipientBlockchainAddress, t.Value)
	}

	data := fmt.Sprintf("%d%s%d%s", b.index, b.previousHash, b.timestamp, transactionsData)
	hash := sha256.Sum256([]byte(data))
	hashString := hex.EncodeToString(hash[:])
	return hashString
}

func newBlock(index int, previousHash string, transactions []*Transaction) *Block {

	b := new(Block)
	b.index = index
	b.previousHash = previousHash
	b.timestamp = time.Now().UnixNano()
	b.transactions = transactions
	b.thisHash = calculateHash(b)
	return b
}

func printBlock(obj Block) {
	fmt.Println("Index:              ", obj.index)
	fmt.Println("Previous Hash:      ", obj.previousHash)
	fmt.Println("Timestamp:          ", obj.timestamp)
	fmt.Println("Current Block Hash: ", obj.thisHash)

	fmt.Println("Block Transactions: ", obj.transactions)

	// for _, t := range obj.transactions {
	// 	fmt.Printf("\tTransaction ID: %s\n", t.TransactionID)
	// 	fmt.Printf("\tSender: %s\n", t.SenderBlockchainAddress)
	// 	fmt.Printf("\tRecipient: %s\n", t.RecipientBlockchainAddress)
	// 	fmt.Printf("\tValue: %f\n", t.Value)
	// 	fmt.Println()
	// }
	transactionsJSON, err := json.MarshalIndent(obj.transactions, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling transactions to JSON:", err)
		return
	}

	fmt.Println("Transactions: ")
	fmt.Println(string(transactionsJSON))
}

func NewBlockchain() *Blockchain {
	genesisTransaction := &Transaction{
		TransactionID:              "0",
		SenderBlockchainAddress:    "0",
		RecipientBlockchainAddress: "0",
		Value:                      0,
	}

	genesisBlock := newBlock(0, "0", []*Transaction{genesisTransaction})

	return &Blockchain{
		blocks: []*Block{genesisBlock},
	}
}

func (bc *Blockchain) addBlock(transactions []*Transaction) {
	previousBlock := bc.blocks[len(bc.blocks)-1]
	previousHash := previousBlock.thisHash

	newBlock := newBlock(len(bc.blocks), previousHash, transactions)

	bc.blocks = append(bc.blocks, newBlock)
}

func (bc *Blockchain) printBlockchain() {
	for _, block := range bc.blocks {
		printBlock(*block)
		fmt.Println()
	}
}

func (bc *Blockchain) verifyChain() bool {

	for i := 1; i < len(bc.blocks); i++ {
		currentBlock := bc.blocks[i]
		previousBlock := bc.blocks[i-1]

		currentBlockHash := calculateHash(currentBlock)

		if currentBlock.thisHash != currentBlockHash {

			fmt.Println(currentBlock.thisHash, currentBlockHash)
			return false
		}
		if currentBlock.previousHash != previousBlock.thisHash {
			fmt.Println(currentBlock.previousHash, previousBlock.thisHash)
			return false
		}
	}
	return true
}

func (bc *Blockchain) modifyBlockChain(index int, newTransactions []*Transaction) error {

	if index < 0 || index >= len(bc.blocks) {
		return fmt.Errorf("Invalid Block index: %d", index)
	}

	bc.blocks[index].transactions = newTransactions

	if index < len(bc.blocks)-1 {
		nextBlock := bc.blocks[index+1]
		nextBlock.previousHash = bc.blocks[index].thisHash
		updatedHash := calculateHash(nextBlock)
		nextBlock.thisHash = updatedHash
	}
	return nil
}

func main() {

	blockchain := NewBlockchain()

	blockchain.AddTransaction("Alice", "Bob", 30)
	blockchain.AddTransaction("Bob", "Charlie", 50)
	blockchain.AddTransaction("Charlie", "Alice", 10)

	fmt.Println("Initial Blockchain")
	blockchain.printBlockchain()

	blockchain.addBlock(blockchain.TransactionPool)

	blockchain.TransactionPool = []*Transaction{}

	fmt.Println("Updated Blockchain")
	blockchain.printBlockchain()

	blockchain.AddTransaction("David", "Eve", 10.0)
	blockchain.AddTransaction("Eve", "Frank", 5.75)

	blockchain.addBlock(blockchain.TransactionPool)
	blockchain.TransactionPool = []*Transaction{}

	blockchain.printBlockchain()

	fmt.Println("Verifying Blockchain Integrity:")
	if blockchain.verifyChain() {
		fmt.Println("Blockchain is valid!")
	} else {
		fmt.Println("Blockchain has been tampered with!")
	}

	err := blockchain.modifyBlockChain(1, []*Transaction{
		{TransactionID: "1234", SenderBlockchainAddress: "Eve", RecipientBlockchainAddress: "Alice", Value: 100.0},
	})

	if err != nil {
		fmt.Println("Error Modifying Blockchain")
	} else {
		blockchain.printBlockchain()
	}

	if blockchain.verifyChain() {
		fmt.Println("Blockchain is valid!")
	} else {
		fmt.Println("Blockchain has been tampered with!")
	}

}
