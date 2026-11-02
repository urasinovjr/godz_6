package ledger

import "errors"

type Transaction struct {
	ID          int
	Amount      float64
	Category    string
	Description string
	Date        string
}

func (t Transaction) Validate() error {
	if t.Amount == 0 {
		return errors.New("сумма транзакции не может быть нулевой")
	}
	if t.Amount < 0 {
		return errors.New("сумма транзакции не может быть отрицательной")
	}
	if t.Category == "" {
		return errors.New("у транзакции не указана категория")
	}
	return nil
}

var ErrBudgetExceeded = errors.New("budget exceeded")

var transactions = []Transaction{}

func AddTransaction(tx Transaction) error {
	err := tx.Validate()
	if err != nil {
		return err
	}

	budget, ok := budgets[tx.Category]
	if ok {
		spent := 0.0
		for _, t := range transactions {
			if t.Category == tx.Category {
				spent = spent + t.Amount
			}
		}

		if spent+tx.Amount > budget.Limit {
			return ErrBudgetExceeded
		}
	}

	tx.ID = len(transactions) + 1
	transactions = append(transactions, tx)
	return nil
}

func ListTransactions() []Transaction {
	result := make([]Transaction, len(transactions))
	copy(result, transactions)
	return result
}

func Reset() {
	transactions = []Transaction{}
	budgets = map[string]Budget{}
}
