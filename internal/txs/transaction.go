package txs

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	cryptoRand "crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
)

type Transaction struct {
	FromU   string
	ToU     string
	Amount  int
	private *ecdsa.PrivateKey
	Public  string
	Sign    string
	ID      string
}

var BalanceU = make(map[string]int)

func GetTxsData() []Transaction {
	var ntx int
	transactions := make([]Transaction, 0)
	fmt.Println("Enter how many transactions you want to perform:")
	fmt.Scan(&ntx)
	for i := 0; i < ntx; i++ {
		var newT Transaction
		fmt.Println("From user:")
		fmt.Scan(&newT.FromU)
		fmt.Println("To user:")
		fmt.Scan(&newT.ToU)
		fmt.Println("VC Amount:")
		fmt.Scan(&newT.Amount)

		if newT.checkValid() {
			fromB, _ := BalanceU[newT.FromU]
			toB, _ := BalanceU[newT.ToU]
			fmt.Printf("Transaction is valid\n"+
				"%s remaining balance: %v VC\n"+
				"%s new balance: %v VC\n"+"\n", newT.FromU, fromB, newT.ToU, toB)
			transactions = append(transactions, newT.ConfirmTx())
		} else {
			fmt.Println("Transaction is not valid")
			if newT.FromU == newT.ToU {
				fmt.Println(
					"You can't send VC from one user to the same user" + "\n")
			} else {
				fromB, _ := BalanceU[newT.FromU]
				fmt.Printf("User %s doesn't have enough VC to send\n"+
					"%s current balance: %v VC\n"+"\n", newT.FromU, newT.FromU, fromB)
			}

		}
	}
	return transactions
}

func (t *Transaction) ConfirmTx() Transaction {
	t.generatePrivateKey()
	t.getPublicKey()
	if t.FromU != "Blockchain" {
		t.signTx()
	}
	t.generateTxID()
	return *t
}

func (t *Transaction) getPublicKey() {
	publicKey, err := x509.MarshalPKIXPublicKey(&t.private.PublicKey)
	if err != nil {
		log.Fatal(err)
	}
	t.Public = base64.StdEncoding.EncodeToString(publicKey)
}

func (t *Transaction) generatePrivateKey() {
	var err error
	t.private, err = ecdsa.GenerateKey(elliptic.P256(), cryptoRand.Reader)
	if err != nil {
		log.Fatal(err)
	}
}

func (t *Transaction) signTx() {
	sha256Hash := sha256.New()
	sha256Hash.Write([]byte(fmt.Sprintf("%s%v%s", t.FromU, t.Amount, t.ToU)))
	hash := sha256Hash.Sum(nil)

	bytes, err := ecdsa.SignASN1(cryptoRand.Reader, t.private, hash[:])
	if err != nil {
		log.Fatal(err)
	}
	t.Sign = base64.StdEncoding.EncodeToString(bytes)
}

func (t *Transaction) generateTxID() {
	data := fmt.Sprintf("%s%v%s", t.FromU, t.Amount, t.ToU)
	binaryData := []byte(fmt.Sprintf("%s%s", data, t.Public))
	if t.Sign != "" {
		binaryData = []byte(fmt.Sprintf("%s%s", binaryData, t.Sign))
	}

	sha256Hash1 := sha256.New()
	sha256Hash1.Write(binaryData)

	sha256Hash2 := sha256.New()
	sha256Hash2.Write(sha256Hash1.Sum(nil))

	t.ID = fmt.Sprintf("%x", sha256Hash2.Sum(nil))
}

func (t *Transaction) checkValid() bool {
	if t.FromU == t.ToU {
		return false
	}
	if _, ok := BalanceU[t.FromU]; !ok {
		BalanceU[t.FromU] = 100
	} else {
		if val, _ := BalanceU[t.FromU]; val < t.Amount {
			return false
		}
	}
	if _, ok := BalanceU[t.ToU]; !ok {
		BalanceU[t.ToU] = 100
	}
	BalanceU[t.ToU] += t.Amount
	BalanceU[t.FromU] -= t.Amount
	return true

}

func (t *Transaction) VerifySignature() error {
	decodedSignature, err := base64.StdEncoding.DecodeString(t.Sign)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %w", err)
	}

	decodedPublicKey, err := base64.StdEncoding.DecodeString(t.Sign)
	if err != nil {
		return fmt.Errorf("failed to decode public key: %w", err)
	}

	pubKeyInterface, err := x509.ParsePKIXPublicKey(decodedPublicKey)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}
	pubKey := pubKeyInterface.(*ecdsa.PublicKey)

	hash := sha256.Sum256([]byte(t.Sign))

	valid := ecdsa.VerifyASN1(pubKey, hash[:], decodedSignature)
	if !valid {
		return errors.New("invalid signature")
	}
	return nil
}
