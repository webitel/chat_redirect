package model

import "time"

type GeneralQuestions struct {
	B2Id        string `json:"b2Id,omitempty"`
	FullName    string `json:"fullName,omitempty"`
	TaxId       string `json:"taxId,omitempty"`
	Passport    string `json:"passport,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

type GeneralBody struct {
	//Type    string `json:"type"`
	Gateway string `json:"gateway"`
	Payload any    `json:"payload"`
}

type CardTransactionQuestions struct {
	B2Id              string    `json:"b2Id,omitempty"`
	FullName          string    `json:"fullName,omitempty"`
	TaxId             string    `json:"taxId,omitempty"`
	Passport          string    `json:"passport,omitempty"`
	PhoneNumber       string    `json:"phoneNumber,omitempty"`
	Location          string    `json:"location,omitempty"`
	CardNumber        string    `json:"cardNumber,omitempty"`
	OperationTime     time.Time `json:"operationTime"`
	OperationAmount   string    `json:"operationAmount,omitempty"`
	OperationCurrency string    `json:"operationCurrency,omitempty"`
}

type TransferDetailsQuestions struct {
	B2Id              string `json:"b2Id,omitempty"`
	FullName          string `json:"fullName,omitempty"`
	TaxId             string `json:"taxId,omitempty"`
	Passport          string `json:"passport,omitempty"`
	PhoneNumber       string `json:"phoneNumber,omitempty"`
	Location          string `json:"location,omitempty"`
	IBAN              string `json:"IBAN,omitempty"`
	OperationAmount   string `json:"operationAmount,omitempty"`
	OperationCurrency string `json:"operationCurrency,omitempty"`
}
