package main

import (
	"fmt"

	"github.com/um198/wallet/pkg/wallet"
)

func main() {
	s := &wallet.Service{}
	account, err := s.RegisterAccount("+992000000007")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = s.Deposit(account.ID, 90)
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

	pay, err := s.Pay(account.ID, 10, "auto")
	if err != nil {
		fmt.Println(wallet.ErrNotEnoughBalance)
		return
	}
	fmt.Println(pay)
	fmt.Println("Баланс после снятия 1: ", account.Balance)

	repp, err := s.Repeat(pay.ID)
	fmt.Println(repp)
	fmt.Println("Баланс после повтора снятия 1: ", account.Balance)
	pp, err := s.FindPaymentByID(repp.ID)
	if err != nil {
		return
	}

	fmt.Println(pp)

	fav, err := s.FavoritePayment(pp.ID, "myFv")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("FavoritePayment:= ", fav)

	favorite, err := s.FindFavoriteByID(fav.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("FindFavoriteByID:= ", favorite)

	pFav, err := s.PayFromFavorite(fav.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("PayFromFavorite:= ", pFav)

	fmt.Println("Баланс после платежа из избранных: ", account.Balance)

	_, err = s.RegisterAccount("+992000000008")
	_, err = s.RegisterAccount("+992000000009")
	_, err = s.RegisterAccount("+992000000010")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = s.ExportToFile("acc.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	// s.ImportFromFile("acc.txt")
	

}
