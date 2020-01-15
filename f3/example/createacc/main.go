package main

import (
	"context"
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

	// create an Account
	r := &form3.Account{
		Data: form3.Data{
			Type:           "accounts",
			ID:             "ad27e265-9605-4b4b-a0e5-3003ea9cc4da",
			OrganisationID: "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
			Version:        0,
			Attributes: form3.Attributes{
				Country:                     "GB",
				BaseCurrency:                "GBP",
				AccountNumber:               "41426819",
				BankID:                      "400300",
				BankIDCode:                  "GBDSC",
				Bic:                         "NWBKGB22",
				Iban:                        "GB11NWBK40030041426819",
				Title:                       "Ms",
				FirstName:                   "Samantha",
				BankAccountName:             "Samantha Holder",
				AlternativeBankAccountNames: []string{"Sam Holder"},
				AccountClassification:       "Personal",
				JointAccount:                false,
				AccountMatchingOptOut:       false,
				SecondaryIdentification:     "A1B2C3D4",
			},
		},
	}

	if err := client.CreateAccount(context.Background(), r); err != nil {
		log.Fatal(err)
	}

	accounts, err := client.ListAccounts(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for _, account := range accounts {
		fmt.Printf("ListAccount :: with ID %s \n", account.Data.ID)
	}

}
