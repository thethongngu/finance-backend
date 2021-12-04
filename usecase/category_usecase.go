package usecase

import (
	"finance/adaptor"
	"finance/entity"
)

type CategoryUsecaseInterface interface {
	GetCategory([]entity.Category, error)
}

type CategoryUsecase struct {
	adaptor *adaptor.CategoryAdaptorInterface
}

func (c *CategoryUsecase) GetCategory() {

}
