package base

import "github.com/huhx/common-go/types"

type CreatedResponse struct {
	Id int64 `json:"id" example:"1800783867820785152"`
}

type BatchCreatedResponse struct {
	Ids types.Int64Array `json:"ids" example:"1800783867820785152,1800783867820785153"`
}
