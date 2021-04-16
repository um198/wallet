package wallet

import (
	"errors"
	"github.com/google/uuid"
	"github.com/um198/wallet/pkg/types"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

var ErrPhoneRegistered = errors.New("phone alredy registered")
var ErrAmmountMustBePositive = errors.New("ammount must be greater then zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("not enough balance ")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrFavoriteNotFound = errors.New("payment not found")

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
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

	account, err := s.FindAccountByID(accountID)
	if err != nil {
		return ErrAccountNotFound
	}
	account.Balance += ammount
	return nil
}

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			return acc, nil
		}
	}
	return nil, ErrAccountNotFound
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, pay := range s.payments {
		if pay.ID == paymentID {
			return pay, nil
		}
	}
	return nil, ErrPaymentNotFound
}

func (s *Service) Reject(paymentID string) error {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil
	}
	acc, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return nil
	}

	payment.Status = types.PaymentStatusFail
	acc.Balance += payment.Amount
	return nil
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmmountMustBePositive
	}

	account, err := s.FindAccountByID(accountID)
	if err != nil {
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

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	newPaymntID := uuid.New().String()
	newPayment := &types.Payment{
		ID:        newPaymntID,
		AccountID: payment.AccountID,
		Amount:    payment.Amount,
		Category:  payment.Category,
		Status:    types.PaymentStatusInProgress,
	}

	acc, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return nil, err
	}

	acc.Balance -= payment.Amount

	s.payments = append(s.payments, newPayment)
	return newPayment, nil

}

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	favorite := uuid.New().String()
	newFavorite := &types.Favorite{
		ID:        favorite,
		AccountID: payment.AccountID,
		Name:      name,
		Amount:    payment.Amount,
		Category:  payment.Category,
	}

	s.favorites = append(s.favorites, newFavorite)

	return newFavorite, nil
}

func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {

	for _, fav := range s.favorites {
		if fav.ID == favoriteID {
			return fav, nil
		}
	}
	return nil, ErrFavoriteNotFound

}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, err
	}

	newPaymntID := uuid.New().String()
	newPayment := &types.Payment{
		ID:        newPaymntID,
		AccountID: favorite.AccountID,
		Amount:    favorite.Amount,
		Category:  favorite.Category,
		Status:    types.PaymentStatusInProgress,
	}

	acc, err := s.FindAccountByID(favorite.AccountID)
	if err != nil {
		return nil, err
	}

	acc.Balance -= favorite.Amount

	s.payments = append(s.payments, newPayment)
	return newPayment, nil

}

func (s *Service) ExportToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Print(err)
		return err
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()
	result := ""
	for _, account := range s.accounts {
		result = result + strconv.FormatInt(account.ID, 10) + ";" +
			string(account.Phone) + ";" + strconv.FormatInt(int64(account.Balance), 10) + "|"
	}

	_, err = file.Write([]byte(result))
	if err != nil {
		log.Print(err)
		return err
	}

	return err
}

func (s *Service) ImportFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Print(err)
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	buf := make([]byte, 4)
	content := make([]byte, 0)
	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			content = append(content, buf[:read]...)
			break
		}
		if err != nil {
			log.Print(err)
			return err
		}
		content = append(content, buf[:read]...)
	}
	data := strings.Split(string(content), "|")

	for _, acc := range data {
		account := (strings.Split(acc, ";"))
		if len(account) > 1 {
			// log.Print(len(account))
			id, err := strconv.ParseInt(account[0], 10, 64)
			if err != nil {
				return err
			}
			balance, err := strconv.ParseInt(account[2], 10, 64)
			if err != nil {
				return err
			}
			account := &types.Account{
				ID:      id,
				Phone:   types.Phone(account[1]),
				Balance: types.Money(balance),
			}
			s.accounts = append(s.accounts, account)
		}
	}

	return err
}
