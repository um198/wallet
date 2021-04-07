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

	account11, err := svc.FindAccountByID(account.ID+1)
	if err != nil {
		if account != account11 {
			t.Error(err)
		}

	}
}
