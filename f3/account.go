package f3

import (
	"context"
	"fmt"
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

// AccountListOptions is used for pagination
type AccountListOptions struct {
	PageNumber int `url:"page_number"`
	PageSize   int `url:"page_size"`
}

// CreateAccount registers a new Account in Form3.
func (c *Client) CreateAccount(ctx context.Context, account *Account) error {
	req, err := c.newRequest("POST", "/v1/organisation/accounts", account)
	if err != nil {
		return err
	}

	resp, err := c.do(ctx, req, nil)
	if resp.StatusCode != http.StatusCreated {
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("DeleteAccount :: unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// FetchAccount returns an account with a given accountID.
func (c *Client) FetchAccount(ctx context.Context, accountID string) (*Account, error) {
	req, err := c.newRequest("GET", "/v1/organisation/accounts"+accountID, nil)
	if err != nil {
		return nil, err
	}

	r := new(Account)
	resp, err := c.do(ctx, req, &r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("FetchAccount :: unexpected status code: %d", resp.StatusCode)
	}

	return r, nil
}

// ListAccounts returns all accounts.
func (c *Client) ListAccounts(ctx context.Context) ([]*Account, error) {
	req, err := c.newRequest("GET", "/v1/organisation/accounts", nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Accounts []*Account `json:"data"`
	}

	resp, err := c.do(ctx, req, &response)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ListAccounts :: unexpected status code: %d", resp.StatusCode)
	}

	return response.Accounts, err
}

// ListAccountsWithPagination returns paginated accounts listing
func (c *Client) ListAccountsWithPagination(ctx context.Context, opts AccountListOptions) ([]*Account, error) {

	u := "/v1/organisation/accounts" + "?page[number]=" + strconv.Itoa(opts.PageNumber) + "&page[size]=" + strconv.Itoa(opts.PageSize)

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
