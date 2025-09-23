package category_model

type Category struct {
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string `json:"name" gorm:"type:varchar(200);not null;uniqueIndex;"`
	Slug      string `json:"slug" gorm:"type:varchar(100);not null;uniqueIndex;"`
	IsActive  bool   `json:"is_active" gorm:"default:true"`
	Order     int    `json:"order"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	DeletedAt int64  `json:"deleted_at"`
}

func (Category) TableName() string {
	return "categories"
}
