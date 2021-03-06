package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"os"
	"runtime"

	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	genesisData = "First Transaction from Genesis"
	lastHashKey = "last_hash"
)

// A BlockChain definition with the BadgerDB configured as a DB.
type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

// A BlockChain iterator allowing to iterate over items in a BlockChain DB.
type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

// Checks if a BlockChain DB exist by checking for a presence of a dbFile.
func DBexists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

// Initialise a Blockchain from an existing DB.
func ContinueBlockChain(address string) *BlockChain {
	if DBexists() == false {
		log.Warn().Msg("No existing blockchain found. A new BlockChain needs to be created.")
		runtime.Goexit()
	}
	var lastHash []byte
	opts := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opts)
	HandleError(err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(lastHashKey))
		HandleError(err)
		err = item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)
			return nil
		})
		HandleError(err)
		return err
	})
	HandleError(err)
	return &BlockChain{lastHash, db}
}

// Initialise a new Blockchain using an address data provided.
func InitBlockChain(address string) *BlockChain {
	var lastHash []byte

	if DBexists() {
		log.Warn().Msg("Blockchain already exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opts)
	HandleError(err)

	err = db.Update(func(txn *badger.Txn) error {
		cbtx := CoinbaseTx(address, genesisData)
		genesis := Genesis(cbtx)
		log.Debug().Msg("Genesis created")
		err = txn.Set(genesis.Hash, genesis.Serialize())
		HandleError(err)
		err = txn.Set([]byte(lastHashKey), genesis.Hash)
		lastHash = genesis.Hash
		return err
	})
	HandleError(err)
	return &BlockChain{lastHash, db}
}

// Adds a new block into the BlockChain DB.
func (chain *BlockChain) AddBlock(transactions []*Transaction) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(lastHashKey))
		HandleError(errors.Wrap(err, ""))
		err = item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)
			return nil
		})
		return err
	})
	HandleError(errors.Wrap(err, ""))

	newBlock := CreateBlock(transactions, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		HandleError(errors.Wrap(err, ""))
		err = txn.Set([]byte(lastHashKey), newBlock.Hash)
		chain.LastHash = newBlock.Hash
		return err
	})
	HandleError(errors.Wrap(err, ""))
}

// Creates a new BlockChainIterator allowing easy iteration over a BlockChain DB.
// This iterator is meant to iterates backwards by tracking a prevHash.
func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}
	return iter
}

// Provides a next item from a currently used BlockChain DB using an BlockChainIterator.
func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		HandleError(errors.Wrap(err, ""))
		var encodedBlock []byte
		err = item.Value(func(val []byte) error {
			encodedBlock = append([]byte{}, val...)
			return nil
		})
		block = Deserialize(encodedBlock)
		return err
	})
	HandleError(errors.Wrap(err, ""))

	iter.CurrentHash = block.PrevHash

	return block
}

// Finds all unspent transactions which are assigned to an address/user.
// An unspent transactions are transactions that have outputs which are not
// referenced by other inputs. These transactions are important because if an output
// hasn't been spent then that means that those tokens still exist for a certain user.
// By counting all the unspent outputs that are assigned to a certain user we can
// find out how many tokens are assigned to that user.

// Inputs are like a debits against an account
// Outputs are like a credits against an account
// A exchanges 100 amount for cash => TxOutput of 100 amount created to an A address
// A paying 10 amount to B => a new Tx created where previous TxOutput (100) is used as TxInput and 10 amount used as an TxOutput
// A's payment to B's account uses a previous transaction's output as its input.

//      Inputs          Outputs
//TX0     0               100
//TX1     100              10
func (chain *BlockChain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
	var unspentTxs []Transaction
	spentTXOs := make(map[string][]int)
	iter := chain.Iterator()

	for {
		// Backward iteration, from the latest Block to the Genesis
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
		Outputs:
			for outIdx, out := range tx.Outputs { // check for outputs which are not referenced by other inputs
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.IsLockedWithKey(pubKeyHash) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}
			if tx.IsCoinbase() == false { // track (non coinbase) inputs in the spentTXOs
				for _, in := range tx.Inputs {
					if in.UsesKey(pubKeyHash) {
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}
	return unspentTxs
}

// Finds all unspent transactions outputs.
func (chain *BlockChain) FindUTXO(pubKeyHash []byte) []TxOutput {
	var UTXOs []TxOutput
	unspentTransactions := chain.FindUnspentTransactions(pubKeyHash)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs
}

// Finds all unspent tokens which can be used as an spendable outputs.
func (chain *BlockChain) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := chain.FindUnspentTransactions(pubKeyHash)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOuts
}

func (bc *BlockChain) FindTransaction(ID []byte) (Transaction, error) {
	iter := bc.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction does not exist")
}

func (bc *BlockChain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTX, err := bc.FindTransaction(in.ID)
		HandleError(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	prevTXs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTX, err := bc.FindTransaction(in.ID)
		HandleError(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}
