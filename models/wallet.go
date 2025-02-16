package models

//Задаём структуры для кошелька и операций

type Wallet struct {
	ID       string  `json:"walletId"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}

type OperationRequest struct {
	WalletID      string  `json:"walletId"`
	OperationType string  `json:"operationType"`
	Amount        float64 `json:"amount"`
}
