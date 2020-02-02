package main

import (
	"context"
	"fmt"
	"log"
	"os"

	form3 "github.com/krishnarbb/go-client-sdk/f3"
)

func main() {
	var (
		baseURL = os.Getenv("FORM3_HOST")
	)

	client := form3.NewClient(
		form3.WithBaseURL(baseURL),
	)

	if err := client.AccountServiceCheck(context.Background()); err != nil {
		log.Fatal(err)
	}

	accountID := "ad27e265-9605-4b4b-a0e5-3003ea9cc4da"

	// create an Account
	r := &form3.Account{
		Data: form3.Data{
			Type:           "accounts",
			ID:             accountID,
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
	fmt.Printf("CreateAccount : with ID %s successful \n", accountID)

	if err := client.DeleteAccount(context.Background(), accountID, 0); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("DeleteAccount :: with ID %s successful \n", accountID)

	// Make a paginated request for Listing accounts
	paginationclient := form3.NewClient(
		form3.WithBaseURL(baseURL),
		form3.WithServicePath("/v1/organisation/accounts"),
		form3.WithPagination(0, 1),
	)
	if err := paginationclient.CreateAccount(context.Background(), r); err != nil {
		log.Fatal(err)
	}

	accounts, err := paginationclient.ListAccounts(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ListAccount :: %d accounts returned\n", len(accounts))
	
	if err := paginationclient.DeleteAccount(context.Background(), accountID, 0); err != nil {
		log.Fatal(err)
	}
}
