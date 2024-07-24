package enet

import (
	"EagleNet/eiface"
)

//实现router时,先嵌入这个BaseRouter,然后根据需要重写
type BaseRouter struct{}

func (b *BaseRouter) PreHandler(req eiface.IRequest) {
}

func (b *BaseRouter) Handler(req eiface.IRequest) {
}

func (b *BaseRouter) PostHandler(req eiface.IRequest) {
}
