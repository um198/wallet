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
	err = svc.Deposit(account.ID, 9)
	if err != nil {
		switch err {
		case wallet.ErrAmmountMustBePositive:
			fmt.Println("Сумма должна быть позитив")
		case wallet.ErrAccountNotFound:
			fmt.Println("Аккаунт не найден")
		}
		return
	}

	fmt.Println("Баланс после поплнения: ", account.Balance)

	pay, err := svc.Pay(account.ID, 1, "auto")
	fmt.Println(pay)
	fmt.Println("Баланс после снятия 1: ", account.Balance)

	err = svc.Reject(pay.ID)
	fmt.Println(pay)
	fmt.Println("Баланс после отмены 1: ",account.Balance)
	pp,err:=svc.FindPaymentByID(pay.ID)
	if err != nil {
		return 
	}

	fmt.Println(pp)

}
