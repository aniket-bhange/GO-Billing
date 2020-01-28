package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

//Product is struct for table colums
type Product struct {
	ID       uint   `gorm:"primary_key;auto_increment" json:"id"`
	Name     string `json:"name"`
	ClientID uint   `json:"client_id"`
	MaxPrice uint   `json:"max_price"`
	MinPrice uint   `json:"min_price"`
	Status   bool   `json:"status"`

	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Product) Disable() {
	p.Status = false
}

func (p *Product) Enable() {
	p.Status = true
}

func (p *Product) Create(db *gorm.DB) (*Product, error) {
	var err error

	err = db.Debug().Create(&p).Error

	if err != nil {
		return &Product{}, err
	}

	db.Save(&p)
	return p, nil
}

func (p *Product) Update(db *gorm.DB, uid uint64) (*Product, error) {
	var err error

	err = db.Debug().Model(&Product{}).Where("id = ?", uid).Take(&Product{}).UpdateColumns(map[string]interface{}{
		"name":      p.Name,
		"client_id": p.ClientID,
		"max_price": p.MaxPrice,
		"min_price": p.MinPrice,
		"status":    p.Status,
	}).Error

	if err != nil {
		return &Product{}, err
	}

	err = db.Debug().Model(&Product{}).Where("id = ?", uid).Take(&p).Error

	if err != nil {
		return &Product{}, err
	}

	return p, nil

}

func (p *Product) FindAll(db *gorm.DB) (*[]Product, error) {
	var err error

	products := []Product{}

	err = db.Debug().Model(&Product{}).Limit(100).Find(&products).Error

	if err != nil {
		return &[]Product{}, err
	}

	return &products, nil
}
