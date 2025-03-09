package domain

type Payment struct {
	PaymentID        string
	CustomerID       string
	PaymentTimestamp int64
	TransactionValue float64
}
