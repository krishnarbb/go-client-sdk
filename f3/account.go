package f3

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Account represenation in Form3
type Account struct {
	Data Data `json:"data"`
}

// Attributes of an account in Form3
type Attributes struct {
	Country                     string   `json:"country"`
	BaseCurrency                string   `json:"base_currency"`
	AccountNumber               string   `json:"account_number"`
	BankID                      string   `json:"bank_id"`
	BankIDCode                  string   `json:"bank_id_code"`
	Bic                         string   `json:"bic"`
	Iban                        string   `json:"iban"`
	Title                       string   `json:"title"`
	FirstName                   string   `json:"first_name"`
	BankAccountName             string   `json:"bank_account_name"`
	AlternativeBankAccountNames []string `json:"alternative_bank_account_names"`
	AccountClassification       string   `json:"account_classification"`
	JointAccount                bool     `json:"joint_account"`
	AccountMatchingOptOut       bool     `json:"account_matching_opt_out"`
	SecondaryIdentification     string   `json:"secondary_identification"`
}

// Data in Form3 account
type Data struct {
	Type           string     `json:"type"`
	ID             string     `json:"id"`
	OrganisationID string     `json:"organisation_id"`
	Version        int        `json:"version"`
	Attributes     Attributes `json:"attributes"`
}

// CreateAccount registers a new Account in Form3.
func (c *Client) CreateAccount(ctx context.Context, account *Account) error {
	req, err := c.newRequest("POST", "/v1/organisation/accounts", account)
	if err != nil {
		return err
	}

	resp, err := c.do(ctx, req, nil)
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		log.Println("CreateAccount :: Successful")
	} else {
		return fmt.Errorf("CreateAccount :: unexpected status code: %d", resp.StatusCode)
	}

	return err
}

// DeleteAccount deletes an account with a given accountID and version.
func (c *Client) DeleteAccount(ctx context.Context, accountID string, version int) error {
	req, err := c.newRequest("DELETE", "/v1/organisation/accounts/"+accountID+"?version="+strconv.Itoa(version), nil)
	if err != nil {
		return err
	}

	resp, err := c.do(ctx, req, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		log.Println("DeleteAccount :: Successful")
	} else {
		return fmt.Errorf("DeleteAccount :: unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// FetchAccount returns an account with a given accountID.
func (c *Client) FetchAccount(ctx context.Context, accountID string) (*Account, error) {
	req, err := c.newRequest("GET", "/v1/organisation/accounts/"+accountID, nil)
	if err != nil {
		return nil, err
	}

	r := new(Account)
	resp, err := c.do(ctx, req, &r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		log.Println("FetchAccount :: Successful")
	} else {
		return nil, fmt.Errorf("FetchAccount :: unexpected status code: %d", resp.StatusCode)
	}

	return r, nil
}

// ListAccounts returns accounts listing. If pagination is requested, only the requested accounts are returned.
func (c *Client) ListAccounts(ctx context.Context) ([]*Account, error) {
	u := c.servicePathURL
	if c.pagination != "" {
		u = c.servicePathURL + c.pagination
	}

	req, err := c.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Accounts []*Account `json:"data"`
	}

	_, err = c.do(ctx, req, &response)

	return response.Accounts, err
}

// AccountServiceCheck waits till Account service is available
func (c *Client) AccountServiceCheck(ctx context.Context) error {
	req, err := c.newRequest("GET", "/v1/health", nil)
	if err != nil {
		return err
	}

	serviceAvailable := false
	for serviceAvailable != true {
		resp, err := c.do(ctx, req, nil)
		if err == nil && resp.StatusCode == http.StatusOK {
			serviceAvailable = true
			break
		}
		log.Print("Waiting for Account service to be available")
	}
	return err
}
