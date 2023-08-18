package pkg

// пакет куда буду записывать новые ошибки для удобства
type Error string

func (e Error) Error() string { return string(e) }

const TokenNotExist = Error("Token does not exist")
const UniqueViolationCode = "23505"
const uniqueViolationOrders = Error(`ERROR: duplicate key value violates unique constraint "orders_ordernumber_userlogin_key (SQLSTATE 23505)`) 
const NoOrders = Error("User doesn't have any orders")
