package types

// Money представляет собой денежную сумму в минимальных единицах (центы, копейки, дирамы и т.д.)
type Money int64

// PaymentCategory представляет категорию, в которой был совершен платеж (авто, аптеки и т.д.)
type PaymentCategory string

// PaymentStatus представляет статус платежа
type PaymentStatus string

// Предопределеннве статусы платежей
const (
	StatusOk         PaymentStatus = "OK"
	StatusFail       PaymentStatus = "FAIL"
	StatusInProgress PaymentStatus = "INPROGRESS"
)

//Payment представляет информацию о платеже.
type Payment struct {
	ID       int
	Amount   Money
	Category PaymentCategory
	Status   PaymentStatus
}

type Phone string

//Account представляет информациюо счете пользователья
type Account struct {
	ID      int64
	Phone   Phone
	Balance Money
}
