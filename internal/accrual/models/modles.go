package models

import "time"

type Order struct {
	OrderNumber int64  `json:"order_number"`
	StatusOrder string `json:"status_order"`
	Accrual     int64  `json:"accrual"`
}

type OrderForRegister struct {
	OrderNumber int64   `json:"order_number"`
	StatusOrder string  `json:"status_order"`
	Goods       []Goods `json:"goods"`
}

type Goods struct {
	Description string `json:"description"`
	Price       int64  `json:"price"`
}

type GoodsWithReward struct {
	Reward []Reward
	OrderForRegister
}
type Reward struct {
	Match      string `json:"match"`
	Reward     int64  `json:"reward"`
	RewardType string `json:"reward_type"`
}

const (
	RegisteredOrder = "REGISTERED"
	ProcessingOrder = "PROCESSING"
	InvalidOrder    = "INVALID"
	ProcessedOrder  = "PROCESSED"
)

const (
	RateLimit = 100
	TimeLimit = 1 * time.Minute
)

const (
	RewardDefault     = 10
	RewardTypeDefault = "%"
)
