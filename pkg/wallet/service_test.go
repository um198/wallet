package wallet

import (
	// "reflect"
	"fmt"
	"testing"
	// "github.com/um198/wallet/pkg/types"
)

func TestService_FindAccountByID_found(t *testing.T) {
	svc := Service{}
	account, err := svc.RegisterAccount("+992000000007")
	if err != nil {
		fmt.Println(err)
		return
	}

	account11, err := svc.FindAccountByID(account.ID)
	if err != nil {
		if account != account11 {
			t.Error(err)
		}

	}
}


func TestService_FindAccountByID_notFound(t *testing.T) {
	svc := Service{}
	account, err := svc.RegisterAccount("+992000000007")
	if err != nil {
		fmt.Println(err)
		return
	}

	account0, err := svc.FindAccountByID(account.ID+1)
	if err != nil {
		if err != ErrAccountNotFound {
			t.Error(account0)
		}

	}
}


func TestService_Reject_ok(t *testing.T) {
	svc := Service{}
	account, err := svc.RegisterAccount("+992000000007")
	if err != nil {
		fmt.Println(err)
		return
	}
	account.Balance=99
	pay, err := svc.Pay(account.ID, 1, "auto")
	err = svc.Reject(pay.ID)
	if account.Balance!=99{
		t.Error("Деньги не были сняты")
	}

}


func TestService_Reject_notFound(t *testing.T) {
	svc := &Service{}
	
	err:= svc.Reject("0")
	if err == nil {
		t.Error(ErrPaymentNotFound)
		return 
	}

}