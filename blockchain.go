package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Block struct {
	index        int
	previousHash string
	timestamp    int64
	transactions []string
	thisHash     string
}

type Blockchain struct {
	blocks []*Block
}

func calculateHash(b *Block) string {
	data := fmt.Sprintf("%d%s%d%s", b.index, b.previousHash, b.timestamp, b.transactions)
	hash := sha256.Sum256([]byte(data))
	hashString := hex.EncodeToString(hash[:])
	return hashString
}

func newBlock(index int, previousHash string, transactions []string) *Block {

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
	fmt.Println("Block Transactions: ", obj.transactions)
	fmt.Println("Current Block Hash: ", obj.thisHash)

}

func NewBlockchain() *Blockchain {
	genesisBlock := newBlock(0, "0", []string{"Genesis Block"})

	return &Blockchain{
		blocks: []*Block{genesisBlock},
	}
}

func (bc *Blockchain) addBlock(transactions []string) {
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
			return false
		}
		if currentBlock.previousHash != previousBlock.thisHash {
			return false
		}
	}
	return true
}

func (bc *Blockchain) modifyBlockChain(index int, newTransactions []string) error {

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
	blockchain.addBlock([]string{"10 sCoin Alice to Bob", "15 sCoin Bob to Charlie"})
	blockchain.addBlock([]string{"8 sCoin Bob to Alice", "1 sCoin Charlie to Alice"})

	blockchain.printBlockchain()

	if blockchain.verifyChain() {
		fmt.Println("Blockchain is valid.")
	} else {
		fmt.Println("Blockchain is invalid.")
	}

	/* Tampering the chain by Bob
	blockchain.blocks[1].transactions[0] = "100 coins to Bob"

	if blockchain.verifyChain() {
		fmt.Println("Blockchain is valid.")
	} else {
		fmt.Println("Blockchain is invalid.")
	}
	*/

	newTransactions := []string{"100 sCoin to Bob"}
	err := blockchain.modifyBlockChain(1, newTransactions)

	if err != nil {
		fmt.Println("Error Modifying Chain", err)
	}

	if blockchain.verifyChain() {
		fmt.Println("Blockchain is Valid")
	} else {
		fmt.Println("Blockchain is Invalid")
	}

}

