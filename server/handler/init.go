package handler

import (
	"github.com/ilya-korotya/solid/usecase"
)

type Handle struct {
	Usecase usecase.UserUsecase
}

func New(uc usecase.UserUsecase) *Handle {
	return &Handle{
		Usecase: uc,
	}
}
