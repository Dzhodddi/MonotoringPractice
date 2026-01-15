package models

import (
	productModels "github.com/Dzhodddi/EcommerceAPI/pkg/shared/products"
	"time"
)

type Order struct {
	ID            uint `gorm:"primaryKey;autoIncrement"`
	CreatedAt     time.Time
	TotalPrice    float64
	AccountID     uint64
	Status        string
	PaymentStatus string
	ProductsInfos []ProductsInfo                  `gorm:"foreignKey:OrderID"`
	Products      []*productModels.OrderedProduct `gorm:"-"`
}
