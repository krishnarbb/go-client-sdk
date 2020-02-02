package f3

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestListAccounts(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := loadJSON("testdata/account_list_response.json")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write(b); err != nil {
			t.Fatal(err)
		}
	}))

	c := NewClient(
		WithBaseURL(ts.URL),
	)

	r, err := c.ListAccounts(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if len(r) != 2 {
		t.Errorf("unexpected number of accounts = %d; want = %d", len(r), 2)
	}

}

func TestListAccountsWithPagination(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := loadJSON("testdata/account_list_response.json")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write(b); err != nil {
			t.Fatal(err)
		}
	}))

	c := NewClient(
		WithBaseURL(ts.URL),
		WithServicePath("/v1/organisation/accounts"),
		WithPagination(1, 1),
	)

	r, err := c.ListAccounts(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if len(r) != 2 {
		t.Errorf("unexpected number of accounts = %d; want = %d", len(r), 2)
	}
}

func TestCreateAccount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		want, err := loadJSON("testdata/account_create_request.json")
		if err != nil {
			t.Fatal(err)
		}
		assertEqualJSON(t, got, want)

		w.WriteHeader(http.StatusOK)
	}))

	c := NewClient(
		WithBaseURL(ts.URL),
	)
	r := &Account{
		Data: Data{
			Type:           "accounts",
			ID:             "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
			OrganisationID: "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
			Version:        0,
			Attributes: Attributes{
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

	if err := c.CreateAccount(context.Background(), r); err != nil {
		t.Fatal(err)
	}
}

func TestFetchAccount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := loadJSON("testdata/account_fetch_response.json")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write(b); err != nil {
			t.Fatal(err)
		}
	}))

	c := NewClient(
		WithBaseURL(ts.URL),
	)

	want := &Account{
		Data: Data{
			Type:           "accounts",
			ID:             "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
			OrganisationID: "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
			Version:        0,
			Attributes: Attributes{
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

	got, err := c.FetchAccount(context.Background(), want.Data.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got = %v; want = %v", got, want)
	}
	fmt.Printf("FetchAccount :: with ID %s \n", got.Data.ID)
}

func TestDeleteAccount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	c := NewClient(
		WithBaseURL(ts.URL),
	)

	if err := c.DeleteAccount(context.Background(), "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc", 0); err != nil {
		t.Fatal(err)
	}
}
