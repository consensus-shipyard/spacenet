package data

import "time"

type FundRequest struct {
	Address string `json:"address"`
}

type AddrInfo struct {
	Amount           uint64    `json:"amount"`
	LatestWithdrawal time.Time `json:"latest_withdrawal"`
}

type TotalInfo struct {
	Amount           uint64    `json:"amount"`
	LatestWithdrawal time.Time `json:"latest_withdrawal"`
}
