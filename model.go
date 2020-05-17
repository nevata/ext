package ext

//Model 基本模型的定义
type Model struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt JSONTime  `json:"create_at"`
	UpdatedAt JSONTime  `json:"update_at"`
	DeletedAt *JSONTime `sql:"index" json:"-"`
}
