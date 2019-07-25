package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rs/zerolog/log"
)

// Wallets file name
const walletFile = "./tmp/wallets.data"

// Wallets struct with map of wallets
type Wallets struct {
	Wallets map[string]*Wallet
}

// Creates wallets struct containing wallets map.
func CreateWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = map[string]*Wallet{}
	err := wallets.LoadFile()
	return &wallets, err
}

// Adds a newly created wallet.
func (ws *Wallets) AddWallet() string {
	wallet := MakeWallet()
	address := fmt.Sprintf("%s", wallet.Address())
	ws.Wallets[address] = wallet
	return address
}

// Gets all addresses for all stored wallets.
func (ws *Wallets) GetAllAddresses() []string {
	var addresses []string
	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}
	return addresses
}

// Gets a wallet by its address.
func (ws *Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

// Loads a decoded wallets from a file.
func (ws *Wallets) LoadFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	var wallets Wallets

	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		return err
	}

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	if err = decoder.Decode(&wallets); err != nil {
		return err
	}
	ws.Wallets = wallets.Wallets
	return nil
}

// Saves an encoded wallets struct to the file
func (ws *Wallets) SaveFile() {
	var content bytes.Buffer

	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	if err := encoder.Encode(ws); err != nil {
		log.Panic().Err(err)
	}

	if err := ioutil.WriteFile(walletFile, content.Bytes(), 0644); err != nil {
		log.Panic().Err(err)
	}
}
