package wallet

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/um198/wallet/pkg/types"
	// "github.com/um198/wallet/pkg/types"
)

type testService struct {
	*Service
}

type testAccount struct {
	phone    types.Phone
	balance  types.Money
	payments []struct {
		amount   types.Money
		category types.PaymentCategory
	}
}

var defaultTestAccount=testAccount {
	phone: "+992000000007",
	balance: 10_000_00,
	payments: []struct{
		amount types.Money
		category types.PaymentCategory
	}{{1_000_00, "auto"},
},
}




func newTestService() *testService {
	return &testService{Service: &Service{}}
}

func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can`t register account, erro = %v", err)
	}

	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can`t deposit account, error = %v", err)
	}

	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can`t make payment, error = %v", err)
		}
	}

	return account, payments, nil
}

func (s *testService) addAccountWithBalance(phone types.Phone, balance types.Money) (*types.Account, error) {
	account, err := s.RegisterAccount(phone)
	if err != nil {
		return nil, fmt.Errorf("can`t register account, error = %v", err)
	}

	err = s.Deposit(account.ID, balance)
	if err != nil {
		return nil, fmt.Errorf("can`t deposit account, error = %v", err)
	}

	return account, nil
}

func TestService_FindAccountByID_found(t *testing.T) {
	s := newTestService()
	_,payments, err:=s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	payment:=payments[0]
	got, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FindAccountByID(): error = %v", err)
		return
	}
	if !reflect.DeepEqual(payment, got) {
		t.Errorf("FindAccountByID(): wrong payment returned = %v", err)
		return
	}
}

func TestService_FindAccountByID_notFound(t *testing.T) {
	s := newTestService()

	_,_,err:=s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return 
	}

	_, err = s.FindPaymentByID(uuid.New().String())
	if err == nil {
		t.Error("FindAccountByID(): must return error, returnes nil")
		return
	}

	if err == ErrAccountNotFound {
		t.Errorf("FindAccountByID(): must return ErrAccountNotFound, returnes %v", err)
		return
	}
}

func TestService_Reject_ok(t *testing.T) {
	s := newTestService()
	_,payments,err:=s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return 
	}
	payment:=payments[0]
	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v", err)
		return
	}
	savedPayment, err:=s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Reject(): can not find payment by id, error = %v", err)
		return
	}
	if savedPayment.Status!=types.PaymentStatusFail{
		t.Errorf("Reject(): status did not changed, error = %v", err)
	}
	savedAccount, err:=s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("Reject(): can not find account by id, error = %v", err)
		return
	}
	if savedAccount.Balance!=defaultTestAccount.balance{
		t.Errorf("Reject(): balance did not changed, error = %v, %v, %v", err,savedAccount.Balance,defaultTestAccount.balance)
		return
	}
}

func TestService_Repeat_sucsses(t *testing.T) {
	s := newTestService()
	_,payments,err:=s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return 
	}

	payment:=payments[0]
	_,err = s.Repeat(payment.ID)
	if err != nil {
		t.Errorf("Repeat(): error = %v", err)
		return
	}
	savedPayment, err:=s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Repeat(): can not find payment by id, error = %v", err)
		return
	}
	if savedPayment.Status!=types.PaymentStatusInProgress{
		t.Errorf("Repeat(): status did not changed, error = %v", err)
	}
	savedAccount, err:=s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("Repeat(): can not find account by id, error = %v", err)
<<<<<<< HEAD
		return
	}
	if savedAccount.Balance==defaultTestAccount.balance{
		t.Errorf("Repeat(): balance did not changed, error = %v", err)
		return
	}
}


func TestService_FavoritePayment_sucsses(t *testing.T) {
	s := newTestService()
	_,payments,err:=s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return 
	}

	payment:=payments[0]
	_, err=s.FavoritePayment(payment.ID,"auto")
	if err != nil {
		t.Errorf("PayFromFavorite(): can not add favorite, error = %v", err)
		return 
	}
}

func TestService_PayFromFavorite_sucsses(t *testing.T) {
	s := newTestService()
	_,payments,err:=s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return 
	}

	payment:=payments[0]
	fv, err:=s.FavoritePayment(payment.ID,"auto")
	if err != nil {
		t.Errorf("PayFromFavorite(): can not add favorite, error = %v", err)
		return 
	}

	_,err=s.PayFromFavorite(fv.ID)
	if err != nil {
		t.Errorf("PayFromFavorite(): can not find favorite, error = %v", err)
		return 
	}

	savedAccount, err:=s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("PayFromFavorite(): can not find account by id, error = %v", err)
		return
	}
	if savedAccount.Balance==defaultTestAccount.balance{
		t.Errorf("PayFromFavorite(): balance did not changed, old = %v, new = %v", defaultTestAccount.balance,savedAccount.Balance)
=======
		return
	}
	if savedAccount.Balance==defaultTestAccount.balance{
		t.Errorf("Repeat(): balance did not changed, error = %v", err)
>>>>>>> master
		return
	}

}



