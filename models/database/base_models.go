package database

import (
	"gorm.io/gorm"
	"time"
)

// BaseModel gorm.Model 的定义
type BaseModel struct {
	ID        uint           `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type StringArray []string // 字符串数组类型
type UintArray []uint     // uint数组类型
