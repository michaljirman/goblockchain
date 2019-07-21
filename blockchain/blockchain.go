package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	dbPath      = "./tmp/blocks"
	lastHashKey = "last_hash"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func InitBlockChain() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opts)
	HandleError(errors.Wrap(err, ""))

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte(lastHashKey)); err == badger.ErrKeyNotFound {
			fmt.Println("aaaa")
			log.Warn().Err(err).Msgf("No existing blockchain found")
			genesis := Genesis()
			log.Debug().Msg("Genesis proved")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			HandleError(errors.Wrap(err, "failed to set value"))
			err = txn.Set([]byte(lastHashKey), genesis.Hash)
			lastHash = genesis.Hash
			return err
		} else {
			item, err := txn.Get([]byte(lastHashKey))
			HandleError(errors.Wrap(err, "failed to get value"))
			err = item.Value(func(val []byte) error {
				lastHash = append([]byte{}, val...)
				return nil
			})
			HandleError(errors.Wrap(err, "failed to get value"))
			return err
		}
	})
	HandleError(errors.Wrap(err, ""))
	return &BlockChain{lastHash, db}
}

func (chain *BlockChain) AddBlock(data string) {
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

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		HandleError(errors.Wrap(err, ""))
		err = txn.Set([]byte(lastHashKey), newBlock.Hash)
		chain.LastHash = newBlock.Hash
		return err
	})
	HandleError(errors.Wrap(err, ""))
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}
	return iter
}

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
