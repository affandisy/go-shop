package domain

type User struct {
	BaseModel
	Email    string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password string `gorm:"type:varchar(255);not null" json:"-"` // "-" agar tidak muncul di JSON
	Name     string `gorm:"type:varchar(100)" json:"name"`
	Phone    string `gorm:"type:varchar(20)" json:"phone"`
	Role     string `gorm:"type:varchar(20);default:'customer'" json:"role"` // customer, admin
	IsActive bool   `gorm:"default:true" json:"is_active"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}
