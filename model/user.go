package model

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

//"Users" is struct for user table
type Users struct {
	ID        uint   `gorm:"primary_key;auto_increment" json:"id"`
	Email     string `gorm:"unique" json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `gorm:"unique" json:"username"`
	Status    bool   `json:"status"`
	ClientID  uint   `json:"client_id"`

	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (u *Users) Disable() {
	u.Status = false
}

func (u *Users) Enable() {
	u.Status = true
}

func (u *Users) BeforeSave() error {
	pwd := GetMD5Hash(u.Password)

	u.Password = pwd
	return nil

}

func (u *Users) VerifyPassword(hashedPassword string, password string) bool {

	hasher := GetMD5Hash(password)
	if hashedPassword == hasher {
		return true
	}
	return false
}

func (u *Users) SaveUser(db *gorm.DB) (*Users, error) {
	var err error

	err = db.Debug().Create(&u).Error
	if err != nil {
		return &Users{}, err
	}
	db.Save(&u)
	return u, nil

}

func (u *Users) UpdateUser(db *gorm.DB, uid uint64) (*Users, error) {
	var err error
	db = db.Debug().Model(&Users{}).Where("id=?", uid).Take(&Users{}).UpdateColumns(
		map[string]interface{}{
			"email":      u.Email,
			"phone":      u.Phone,
			"first_name": u.FirstName,
			"last_name":  u.LastName,
			"username":   u.Username,
			"update_at":  time.Now(),
			"status":     u.Status,
		},
	)
	if db.Error != nil {
		return &Users{}, db.Error
	}

	err = db.Debug().Model(&Users{}).Where("id=?", uid).Take(&u).Error
	if err != nil {
		return &Users{}, err
	}
	return u, nil

}

func (u *Users) FindAll(db *gorm.DB) (*[]UserResponse, error) {

	var err error

	users := []Users{}
	usersRes := UserResponse{}
	Result := []UserResponse{}

	err = db.Debug().Model(&Users{}).Limit(100).Find(&users).Error

	if err != nil {
		return &[]UserResponse{}, err
	}

	if len(users) > 0 {
		for i := range users {
			usersRes.User = users[i]
			err = db.Debug().Model(&Client{}).Where("id = ?", users[i].ClientID).Take(&usersRes.ClientRef).Error
			if err != nil {
				return &[]UserResponse{}, err
			}
			Result = append(Result, usersRes)
		}
	}

	return &Result, nil
}

func (u *Users) FindOne(db *gorm.DB, uid uint64) (*UserResponse, error) {
	var err error
	// var c Client
	user := Users{}
	Result := UserResponse{}

	row := db.Debug().Model(&Users{}).Where("id = ?", uid).First(&user).Row()
	if err != nil {
		return &UserResponse{}, err
	}
	log.Print(row)
	return &Result, err
}

// func DBMigrate(db *gorm.DB) *gorm.DB {
// 	db.AutoMigrate(&Users{})
// 	return db
// }
