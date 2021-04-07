package main

import (
	"fmt"

	"github.com/um198/wallet/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}
	account, err := svc.RegisterAccount("+992000000007")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(account)

	account11, err:=svc.FindAccountByID(account.ID)
	if err != nil {
		return 
	}

	fmt.Println(account11)

}
