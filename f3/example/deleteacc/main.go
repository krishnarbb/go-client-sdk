package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	form3 "f3"
)

func main() {
	var (
		baseURL = "http://192.168.99.100:8080/"
	)

	client := form3.NewClient(
		form3.WithBaseURL(baseURL),
	)

	accounts, err := client.ListAccounts(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for _, account := range accounts {
		resAcc, _ := json.Marshal(account)
		fmt.Println("ListAccount : " + string(resAcc))
		//fmt.Printf("ListAccount :: with ID %+v\n", account)
	}
	for _, account := range accounts {
		if err := client.DeleteAccount(context.Background(), account.Data.ID, account.Data.Version); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("DeleteAccount :: with ID %s \n", account.Data.ID)
	}

}
