// +build integration

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"gopkg.in/h2non/baloo.v3"

	_ "github.com/lib/pq"
)

const (
	// TODO: defines the remaining required fieds
	paymentSchema = `{
		"title": "payment resource",
		"type": "object",
		"properties":{
			"id": {
				"type": "string"
			},
			"type":{
				"type": "string"
			},
			"version":{
				"type": "integer"
			},
			"organisation_id":{
				"type": "string"
			},
			"created_at":{
				"type": "string"
			},
			"updated_at":{
				"type": "string"
			}
		},
		"required":["id","type", "version", "organisation_id", "created_at"]
	}`

	filteredPaymentsSchema = `{
		"title": "paginated payments resource",
		"type": "object",
		"properties":{
			"limit": {
				"type": "integer"
			},
			"offset": {
				"type": "integer"
			},
			"total_count": {
				"type": "integer"
			},
			"results": {
				"type": "array"
			}
		},
		"required":["limit","offset","total_count","results"]
	}`

	errorSchema = `{
		"title": "error response",
		"type": "object",
		"properties": {
			"error": {
				"type": "object",
				"properties": {
					"code": {
						"type": "string"
					},
					"message": {
						"type": "string"
					}
				},
				"required": ["code", "message"]
			}
		},
		"required": ["error"]
}`

	requestBody = `{
		"type": "Withdraw",
		"version": 0,
		"organisation_id": "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb",
		"attributes": {
		"amount": "356",
		"beneficiary_party": {
			"account_name": "C Parisi",
			"account_number": "5431230",
			"account_number_code": "BBAN",
			"account_type": 1,
			"address": "10 Avenue de Rome",
			"bank_id": "403000",
			"bank_id_code": "GBDSC",
			"name": "Cédric Parisi"
		},
		"charges_information": {
			"bearer_code": "SHAR",
			"sender_charges": [
			{ "amount": "5.00", "currency": "GBP" },
			{ "amount": "10.00", "currency": "USD" },
			{ "amount": "1.00", "currency": "EUR" }
			],
			"receiver_charges_amount": "3.00",
			"receiver_charges_currency": "USD"
		},
		"currency": "GBP",
		"debtor_party": {
			"account_name": "G Orwell",
			"account_number": "GB29XABC10161234567801",
			"account_number_code": "IBAN",
			"address": "18 Westside",
			"bank_id": "203301",
			"bank_id_code": "GBDSC",
			"name": "George Orwell"
		},
		"end_to_end_reference": "Feb garden table",
		"fx": {
			"contract_reference": "CP123",
			"exchange_rate": "1.00000",
			"original_amount": "190.42",
			"original_currency": "EUR"
		},
		"numeric_reference": "1002001",
		"payment_purpose": "Paying for goods/services",
		"payment_scheme": "FPS",
		"payment_type": "Credit",
		"processing_date": "2017-01-18",
		"reference": "Payment for Jo's garden table",
		"scheme_payment_sub_type": "InternetBanking",
		"scheme_payment_type": "ImmediatePayment",
		"sponsor_party": {
			"account_number": "56781234",
			"bank_id": "456456",
			"bank_id_code": "GBDSC"
		}
		}
	}`

	updateRequest = `{
		"id": "216d4da9-e59a-4cc6-8df3-3da6e7580b77",
		"type": "Payment",
		"version": 1,
		"organisation_id": "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb",
		"attributes": {
			"id": "08be82ea-d433-4db7-9924-a057450d082b",
			"payment_id": "216d4da9-e59a-4cc6-8df3-3da6e7580b77",
			"amount": "100.21",
			"beneficiary_party": {
				"attribute_id": "08be82ea-d433-4db7-9924-a057450d082b",
				"account_name": "W Owens",
				"account_number": "31926819",
				"account_number_code": "BBAN",
				"account_type": 0,
				"address": "1 The Beneficiary Localtown SE2",
				"bank_id": "403000",
				"bank_id_code": "GBDSC",
				"name": "Wilfred Jeremiah Owens"
			},
			"charges_information": {
				"id": "a0d5dd15-fa62-4db1-8669-97bc74f00eb4",
				"attribute_id": "08be82ea-d433-4db7-9924-a057450d082b",
				"bearer_code": "SHAR",
				"sender_charges": [
					{
						"charges_information_id": "a0d5dd15-fa62-4db1-8669-97bc74f00eb4",
						"amount": "5.00",
						"currency": "GBP"
					}
				],
				"receiver_charges_amount": "1.00",
				"receiver_charges_currency": "USD"
			},
			"currency": "GBP",
			"debtor_party": {
				"attribute_id": "08be82ea-d433-4db7-9924-a057450d082b",
				"account_name": "EJ Brown Black",
				"account_number": "GB29XABC10161234567801",
				"account_number_code": "IBAN",
				"address": "10 Debtor Crescent Sourcetown NE1",
				"bank_id": "203301",
				"bank_id_code": "GBDSC",
				"name": "Emelia Jane Brown"
			},
			"end_to_end_reference": "Wil piano Jan",
			"fx": {
				"attribute_id": "08be82ea-d433-4db7-9924-a057450d082b",
				"contract_reference": "FX123",
				"exchange_rate": "2.00000",
				"original_amount": "200.42",
				"original_currency": "USD"
			},
			"numeric_reference": "1002001",
			"payment_purpose": "Paying for goods/services",
			"payment_scheme": "FPS",
			"payment_type": "Credit",
			"processing_date": "2017-01-18",
			"reference": "Payment for Em's piano lessons",
			"scheme_payment_sub_type": "InternetBanking",
			"scheme_payment_type": "ImmediatePayment",
			"sponsor_party": {
				"attribute_id": "08be82ea-d433-4db7-9924-a057450d082b",
				"account_number": "56781234",
				"bank_id": "123123",
				"bank_id_code": "GBDSC"
			}
		},
		"created_at": "2019-02-06T13:03:18.411342Z",
		"updated_at": "0001-01-01T00:00:00Z"
	}`

	baseURL = "http://localhost:8000"
)

func getToken() string {
	token := struct {
		Token string `json:"token"`
	}{}
	resp, _ := http.Post(baseURL+"/auth/", "application/json", strings.NewReader(`{
		"id":"cedric",
		"secret":"secret"
	}`))
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&token)
	return token.Token
}

func Test_CreatePayment(t *testing.T) {
	client := baloo.New(baseURL)
	tests := []struct {
		name               string
		body               string
		expectedStatusCode int
		expectedSchema     string
		callMocks          func() string
	}{
		{
			name:               "create payment ok",
			body:               requestBody,
			expectedStatusCode: http.StatusCreated,
			expectedSchema:     paymentSchema,
			callMocks: func() string {
				return getToken()
			},
		},
		{
			name:               "create payment failed due to not authentified",
			body:               requestBody,
			expectedStatusCode: http.StatusUnauthorized,
			expectedSchema:     errorSchema,
			callMocks: func() string {
				return "wrong_token"
			},
		},
		{
			name:               "create payment failed due empty body",
			body:               "",
			expectedStatusCode: http.StatusInternalServerError,
			expectedSchema:     errorSchema,
			callMocks: func() string {
				return getToken()
			},
		},
		{
			name: "create payment failed due invalid payment request",
			body: `{
				"type": "unknown_type"
			}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedSchema:     errorSchema,
			callMocks: func() string {
				return getToken()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.callMocks()

			client.Post("/payments/").
				Body(strings.NewReader(tt.body)).
				SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
				Expect(t).
				Status(tt.expectedStatusCode).
				JSONSchema(tt.expectedSchema).
				Done()
		})
	}
}

func Test_GetPayment(t *testing.T) {
	client := baloo.New(baseURL)
	tests := []struct {
		name               string
		id                 string
		expectedStatusCode int
		expectedSchema     string
	}{
		{
			// TODO: find a payment id instead of hard coding one
			name:               "get payment ok",
			id:                 "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
			expectedStatusCode: http.StatusOK,
			expectedSchema:     paymentSchema,
		},
		{
			name:               "get payment not found",
			id:                 "00000000-0000-0000-0000-000000000000",
			expectedStatusCode: http.StatusNotFound,
			expectedSchema:     errorSchema,
		},
		{
			name:               "get payment invalid id",
			id:                 "not_a_uuid",
			expectedStatusCode: http.StatusBadRequest,
			expectedSchema:     errorSchema,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.Get(fmt.Sprintf("/payments/%s", tt.id)).
				Expect(t).
				Status(tt.expectedStatusCode).
				JSONSchema(tt.expectedSchema).
				Done()
		})
	}
}

func Test_GetFilteredPayment(t *testing.T) {
	client := baloo.New(baseURL)
	tests := []struct {
		name               string
		params             map[string]string
		expectedStatusCode int
		expectedSchema     string
	}{
		{
			name: "get filtered payments ok",
			params: map[string]string{
				"offset": "0",
				"limit":  "2",
				"sort":   "-views",
			},
			expectedStatusCode: http.StatusOK,
			expectedSchema:     filteredPaymentsSchema,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.Get("/payments/").
				Params(tt.params).
				Expect(t).
				Status(tt.expectedStatusCode).
				JSONSchema(tt.expectedSchema).
				Done()
		})
	}
}

func Test_UpdatePayment(t *testing.T) {
	client := baloo.New(baseURL)
	tests := []struct {
		name               string
		id                 string
		body               string
		expectedStatusCode int
		expectedSchema     string
		callMocks          func() string
	}{
		{
			// TODO: find a payment, update a field
			name:               "update payment ok",
			id:                 "216d4da9-e59a-4cc6-8df3-3da6e7580b77",
			body:               updateRequest,
			expectedStatusCode: http.StatusNoContent,
			callMocks: func() string {
				return getToken()
			},
		},
		{
			name:               "update payment failed due to id mismatch",
			id:                 "516d4da9-e59a-4cc6-8df3-3da6e7580b75",
			body:               updateRequest,
			expectedStatusCode: http.StatusBadRequest,
			callMocks: func() string {
				return getToken()
			},
		},
		{
			name:               "update payment failed due to no body",
			id:                 "216d4da9-e59a-4cc6-8df3-3da6e7580b77",
			body:               "",
			expectedStatusCode: http.StatusInternalServerError,
			expectedSchema:     errorSchema,
			callMocks: func() string {
				return getToken()
			},
		},
		{
			name: "update payment failed due to invalid request",
			id:   "216d4da9-e59a-4cc6-8df3-3da6e7580b77",
			body: `{
				"type": "unknown_type"
			}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedSchema:     errorSchema,
			callMocks: func() string {
				return getToken()
			},
		},
	}
	for _, tt := range tests {
		token := tt.callMocks()

		t.Run(tt.name, func(t *testing.T) {
			client.Put(fmt.Sprintf("/payments/%s", tt.id)).
				Body(strings.NewReader(tt.body)).
				SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
				Expect(t).
				Status(tt.expectedStatusCode).
				Done()
			// TODO: assert updated_at
		})
	}
}

func Test_DeletePayment(t *testing.T) {
	client := baloo.New(baseURL)
	tests := []struct {
		name               string
		id                 string
		body               string
		expectedStatusCode int
		expectedSchema     string
		callMocks          func() string
	}{
		{
			name:               "delete payment invalid id",
			id:                 "not_a_uuid",
			expectedStatusCode: http.StatusBadRequest,
			expectedSchema:     errorSchema,
			callMocks: func() string {
				return getToken()
			},
		},
		{
			// TODO: find a payment
			name:               "delete payment ok",
			id:                 "7eb8277a-6c91-45e9-8a03-a27f82aca350",
			expectedStatusCode: http.StatusNoContent,
			callMocks: func() string {
				return getToken()
			},
		},
		{
			name:               "delete payment not found",
			id:                 "00000000-0000-0000-0000-000000000000",
			expectedStatusCode: http.StatusNoContent,
			callMocks: func() string {
				return getToken()
			},
		},
	}
	for _, tt := range tests {
		token := tt.callMocks()

		t.Run(tt.name, func(t *testing.T) {
			client.Delete(fmt.Sprintf("/payments/%s", tt.id)).
				SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
				Expect(t).
				Status(tt.expectedStatusCode).
				Done()
			// TODO: assert deleted_at
		})
	}
}
