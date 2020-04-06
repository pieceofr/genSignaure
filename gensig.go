package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	sdk "github.com/bitmark-inc/bitmark-sdk-go"
	"github.com/bitmark-inc/bitmark-sdk-go/account"
)

const seedFileName = "accountseed.sign"

func main() {
	//fmt.Println("signature:", genSignature())
	//seed, filepath := newBitmarkAccount()
	//fmt.Println("New Account:  Seed:", seed, " Filepath:", filepath)
	genSignature(seedFileName)
}

func newBitmarkAccount() (string, string) {
	sdk.Init(&sdk.Config{
		Network:    sdk.Testnet,
		APIToken:   "bmk - lljpzkhqdkzmblhg",
		HTTPClient: http.DefaultClient,
	})

	a, err := account.New()
	fmt.Println("Account Number:", a.AccountNumber())
	seed := a.Seed()
	seedFile := filepath.Join(seedFileName)

	f, err := os.OpenFile(seedFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("SEED:%s", seed))
	if err != nil {
		panic(err)
	}

	return seed, seedFile
}

func getSeedFromFile(seedFile string) (string, error) {
	f, err := os.Open(seedFile)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, f)
	if err != nil {
		return "", err
	}
	seed := strings.Trim(strings.Split(buf.String(), ":")[1], "\n")
	return seed, nil
}

func genSignature(filepath string) string {
	seed, err := getSeedFromFile(filepath)
	if err != nil {
		panic(err)
	}
	fmt.Println("Get Frome file Seed:", seed)
	timestamp := fmt.Sprintf("%d", time.Now().Unix()*1000)
	fmt.Println("timestamp:", timestamp)
	sdk.Init(&sdk.Config{
		Network:    sdk.Testnet,
		APIToken:   "bmk - lljpzkhqdkzmblhg",
		HTTPClient: http.DefaultClient,
	})

	sender, err := account.FromSeed(seed)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("sender accounut:", sender.AccountNumber())
	fmt.Println("enc pub key:", hex.EncodeToString(sender.(*account.AccountV2).EncrKey.PublicKeyBytes()))
	signature := sender.Sign([]byte(timestamp))
	fmt.Println("signature:", hex.EncodeToString(signature))
	return hex.EncodeToString(signature)
}
