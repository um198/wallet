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
	err = svc.Deposit(account.ID, 90)
	if err != nil {
		switch err {
		case wallet.ErrAmmountMustBePositive:
			fmt.Println("Сумма должна быть позитив")
		case wallet.ErrAccountNotFound:
			fmt.Println("Аккаунт не найден")
		}
		return
	}

	fmt.Println("Баланс: ", account.Balance)

	pay, err := svc.Pay(account.ID, 22, "auto")
	if err != nil {
		fmt.Println(wallet.ErrNotEnoughBalance)
		return
	}
	fmt.Println(pay)
	fmt.Println("Баланс после снятия 1: ", account.Balance)

	repp, err := svc.Repeat(pay.ID)
	fmt.Println(repp)
	fmt.Println("Баланс после повтора снятия 1: ", account.Balance)
	
	pp,err:=svc.FindPaymentByID(repp.ID)
	if err != nil {
		fmt.Println(err, repp.ID)
		return
	}

	fmt.Println("mkmokm     ",pp)

}
