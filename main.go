package main

import (
    "context"
    "crypto/ecdsa"
    "crypto/rand"
    "fmt"
    "log"
    "math/big"

    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/ethclient"
)

func generateRandomAddress() common.Address {
    key, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
    if err != nil {
        log.Fatal(err)
    }
    return crypto.PubkeyToAddress(key.PublicKey)
}

func main() {
    client, err := ethclient.Dial("https://rpc-testnet.unit0.dev")
    if err != nil {
        log.Fatal(err)
    }

    privateKey, err := crypto.HexToECDSA("YOUPRIVATEKEYHERE!")
    if err != nil {
        log.Fatal(err)
    }

    publicKey := privateKey.Public()
    publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
    if !ok {
        log.Fatal("error casting public key to ECDSA")
    }

    fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
    nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
    if err != nil {
        log.Fatal(err)
    }

    gasPrice, err := client.SuggestGasPrice(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    for i := 0; i < 10000; i++ {
        toAddress := generateRandomAddress()
        value := big.NewInt(1000000) // 1 ETH in wei

        // Menggunakan EstimateGas untuk mendapatkan gas limit yang akurat
        msg := ethereum.CallMsg{
            From:  fromAddress,
            To:    &toAddress,
            Value: value,
            Data:  nil,
        }

        gasLimit, err := client.EstimateGas(context.Background(), msg)
        if err != nil {
            log.Fatal(err)
        }

        // Buat transaksi dengan gasLimit dari EstimateGas
        tx := types.NewTransaction(nonce+uint64(i), toAddress, value, gasLimit, gasPrice, nil)

        chainID, err := client.NetworkID(context.Background())
        if err != nil {
            log.Fatal(err)
        }

        signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
        if err != nil {
            log.Fatal(err)
        }

        err = client.SendTransaction(context.Background(), signedTx)
        if err != nil {
            log.Fatal(err)
        }

        fmt.Printf("Transaction sent to %s with tx hash: %s\n", toAddress.Hex(), signedTx.Hash().Hex())
    }
}
