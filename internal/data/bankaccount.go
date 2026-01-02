package data

import (
	"fmt"
	"strings"
)

type BankAccount struct {
	Name             string `json:"name"`
	Amount           string `json:"amount,omitempty"`
	AccountNumber    string `json:"account_number"`
	Balance          string `json:"balance"`
	CreditCardNumber string `json:"credit_card_number,omitempty"`
	RoutingNumber    string `json:"routing_number"`
}

// GenerateDeterministicAccountName generates a deterministic bank account name
// Examples: "Checking Account", "Savings Account", "Money Market"
func GenerateDeterministicAccountName(id, realName string, index int) string {
	if realName == "" {
		return ""
	}
	accountTypes := []string{
		"Checking Account",
		"Savings Account",
		"Money Market Account",
		"Deposit Account",
		"Investment Account",
		"Business Account",
		"Interest-Bearing Account",
		"Premium Savings",
		"High-Yield Savings",
		"Retirement Account",
	}

	fieldType := fmt.Sprintf("account_name_%d", index)
	hash := hashField(id, fieldType, realName)
	return selectFromList(hash, 0, accountTypes)
}

// GenerateDeterministicAmount generates a deterministic dollar amount
// Returns amount as a string (e.g., "1234.56")
func GenerateDeterministicAmount(id, realAmount string, index int) string {
	if realAmount == "" {
		return ""
	}

	fieldType := fmt.Sprintf("account_amount_%d", index)
	hash := hashField(id, fieldType, realAmount)

	// Generate realistic amounts: $100 to $999,999.99
	// Use multiple bytes to create a larger number
	cents := (int(hash[0])<<8 | int(hash[1])) % 100
	dollars := (int(hash[2])<<16 | int(hash[3])<<8 | int(hash[4])) % 1000000
	if dollars < 100 {
		dollars += 100
	}

	return fmt.Sprintf("%.2f", float64(dollars)+float64(cents)/100.0)
}

// GenerateDeterministicAccountNumber generates a deterministic bank account number (10-12 digits)
func GenerateDeterministicAccountNumber(id, realAccountNumber string, index int) string {
	if realAccountNumber == "" {
		return ""
	}
	fieldType := fmt.Sprintf("account_number_%d", index)
	hash := hashField(id, fieldType, realAccountNumber)
	// Generate 10-12 digit number
	num := int(hash[0])<<24 | int(hash[1])<<16 | int(hash[2])<<8 | int(hash[3])
	num = num % 900000000000
	if num < 1000000000 {
		num += 1000000000
	}
	return fmt.Sprintf("%d", num)
}

// GenerateDeterministicBalance generates a deterministic balance as a string (e.g., "12345.67")
func GenerateDeterministicBalance(id, realBalance string, index int) string {
	if realBalance == "" {
		return ""
	}
	fieldType := fmt.Sprintf("balance_%d", index)
	hash := hashField(id, fieldType, realBalance)
	cents := (int(hash[0])<<8 | int(hash[1])) % 100
	dollars := (int(hash[2])<<16 | int(hash[3])<<8 | int(hash[4])) % 1000000
	if dollars < 100 {
		dollars += 100
	}
	return fmt.Sprintf("%.2f", float64(dollars)+float64(cents)/100.0)
}

// GenerateDeterministicRoutingNumber generates a deterministic 9-digit routing number
func GenerateDeterministicRoutingNumber(id, realRoutingNumber string, index int) string {
	if realRoutingNumber == "" {
		return ""
	}
	fieldType := fmt.Sprintf("routing_number_%d", index)
	hash := hashField(id, fieldType, realRoutingNumber)
	num := int(hash[0])<<16 | int(hash[1])<<8 | int(hash[2])
	num = num % 900000000
	if num < 100000000 {
		num += 100000000
	}
	return fmt.Sprintf("%09d", num)
}

// GenerateDeterministicCreditCardNumber generates a deterministic credit card number
// Uses Luhn algorithm to ensure validity
func GenerateDeterministicCreditCardNumber(id, realCCNumber string, index int) string {
	if realCCNumber == "" {
		return ""
	}

	fieldType := fmt.Sprintf("credit_card_number_%d", index)
	hash := hashField(id, fieldType, realCCNumber)

	// Generate a 16-digit credit card number
	ccNum := make([]int, 16)

	// Set first digit to 4 (Visa) for simplicity
	ccNum[0] = 4

	// Fill in digits 1-14 using hash bytes
	for i := 1; i < 15; i++ {
		ccNum[i] = int(hash[i%8]) % 10
	}

	// Calculate Luhn check digit for the last digit
	sum := 0
	for i := range 15 {
		digit := ccNum[14-i]
		if i%2 == 0 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	checkDigit := (10 - (sum % 10)) % 10
	ccNum[15] = checkDigit

	// Convert to string
	var ccNumberStr strings.Builder
	for _, digit := range ccNum {
		fmt.Fprintf(&ccNumberStr, "%d", digit)
	}

	return ccNumberStr.String()
}

// ObscureBankAccounts takes real bank account data and returns deterministic fake accounts
func ObscureBankAccounts(id string, real []*BankAccount) []*BankAccount {
	if len(real) == 0 {
		return nil
	}

	fake := make([]*BankAccount, len(real))
	for i, account := range real {
		fake[i] = &BankAccount{
			Name:             GenerateDeterministicAccountName(id, account.Name, i),
			Amount:           GenerateDeterministicAmount(id, account.Amount, i),
			AccountNumber:    GenerateDeterministicAccountNumber(id, account.AccountNumber, i),
			Balance:          GenerateDeterministicBalance(id, account.Balance, i),
			CreditCardNumber: GenerateDeterministicCreditCardNumber(id, account.CreditCardNumber, i),
			RoutingNumber:    GenerateDeterministicRoutingNumber(id, account.RoutingNumber, i),
		}
	}
	return fake
}
