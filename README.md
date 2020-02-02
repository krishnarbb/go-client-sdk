# Go client for accessing a public API

Provides access to a public Account service REST API of a financial company.

## Installation

```
go get github.com/krishnarbb/go-client-sdk
```

## Run

```
docker-compose up

The above command should start a Mock Account service API and also run the 
go-client-sdk.

To run the client manually, 
FORM3_HOST environment variable should be set to the URL of the Form3 API Service.

For ex :
FORM3_HOST="http://localhost:8080" go run main.go

```
## Unit tests

```
go test github.com/krishnarbb/go-client-sdk/f3
```

## Design Decisions

```
1. Configurable base URL of the Form3 REST service. FORM3_HOST environment variable shoudl be set to the Endpoint URL.
   This makes it easy to test this locally.

2. Use of context object to the client methods, users can maintain control, say when a request runs too long, they can cancel it.
   
    ctx := context.Background()
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)   // timeout if the request takes more than 5 seconds
    defer cancel()
    account, err := client.FetchAccount(context.Background(), accountID)

3. Generic request handling (used by service). 
Refer to : client.go :: 
    func (c *Client) newRequest(method, path string, body interface{}) 
    func (c *Client) do(req *http.Request, v interface{}) 

Account service client APIs which uses the generic client.go 
    func (c *Client) CreateAccount(ctx context.Context, account *Account) error 
    func (c *Client) DeleteAccount(ctx context.Context, accountID string, version int) error 
    func (c *Client) FetchAccount(ctx context.Context, accountID string) (*Account, error) 
    func (c *Client) ListAccounts(ctx context.Context) ([]*Account, error) 
    func (c *Client) AccountServiceCheck(ctx context.Context) error 

3. Pagination of List Accounts : (refer main.go)
        paginationclient := form3.NewClient(
                form3.WithBaseURL(baseURL),
                form3.WithServicePath("/v1/organisation/accounts"),
                form3.WithPagination(1, 1),
        )
        accounts, err := paginationclient.ListAccounts(context.Background())
```

## Sample client

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	form3 "github.com/bujkri/go-client-sdk/f3"
)

func main() {
	var (
		baseURL = os.Getenv("FORM3_HOST")
	)

	client := form3.NewClient(
		form3.WithBaseURL(baseURL),
        )

    	accountID := "ad27e265-9605-4b4b-a0e5-3003ea9cc4da"
    	account, err := client.FetchAccount(context.Background(), accountID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("FetchAccount :: %+v\n", account)
}    
```
