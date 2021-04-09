package wallet

import (
	"errors"
	"github.com/google/uuid"
	"github.com/um198/wallet/pkg/types"
)

var ErrPhoneRegistered = errors.New("phone alredy registered")
var ErrAmmountMustBePositive = errors.New("ammount must be greater then zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("not enough balance ")
var ErrPaymentNotFound = errors.New("payment not found")

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}
	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)

	return account, nil
}

func (s *Service) Deposit(accountID int64, ammount types.Money) error {
	if ammount <= 0 {
		return ErrAmmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return ErrAccountNotFound
	}

	account.Balance += ammount
	return nil
}


func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}
	return account, nil
}


func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	var payment *types.Payment
	for _, pay := range s.payments {
		if pay.ID == paymentID {
			payment = pay
			break
		}
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}
	return payment, nil
}


func (s *Service) Reject(paymentID string) error  {
	payment,err:=s.FindPaymentByID(paymentID)
	if err != nil {
		return nil
	}
	acc,err:=s.FindAccountByID(payment.AccountID)
	if err != nil {
		return nil 
	}

	payment.Status=types.PaymentStatusFail
	acc.Balance+=payment.Amount
	return nil
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmmountMustBePositive
	}
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}
	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}
	account.Balance -= amount
	PaymntID := uuid.New().String()
	payment := &types.Payment{
		ID:        PaymntID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil

}


func (s *Service) Repeat(paymentID string) (*types.Payment, error)  {
	
	payment,err:=s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	PaymntID := uuid.New().String()
	newPayment := &types.Payment{
		ID:        PaymntID,
		AccountID: payment.AccountID,
		Amount:    payment.Amount,
		Category:  payment.Category,
		Status:    types.PaymentStatusInProgress,
	}

	acc, err:=s.FindAccountByID(s.nextAccountID)
	if err != nil {
		return nil, err
	}

	acc.Balance-=payment.Amount

	s.payments = append(s.payments, payment)
	return newPayment, nil	


}




