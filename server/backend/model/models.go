package model

import "time"

// Prompt 提示词管理表
type Prompt struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Category  string    `gorm:"type:varchar(50);default:'default'" json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Task 任务执行记录表
type Task struct {
	ID         uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PromptID   uint64     `gorm:"not null" json:"prompt_id"`
	Status     string     `gorm:"type:enum('pending','processing','completed','failed');default:'pending'" json:"status"`
	LastRun    *time.Time `json:"last_run"`
	RetryCount int        `gorm:"default:0" json:"retry_count"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	Prompt     Prompt     `gorm:"foreignKey:PromptID" json:"prompt,omitempty"`
}

// Result 任务执行结果与分析表
type Result struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID         uint64    `gorm:"uniqueIndex;not null" json:"task_id"`
	ResponseText   string    `gorm:"type:longtext" json:"response_text"`
	BrandScore     float64   `gorm:"type:decimal(5,2);default:0.00" json:"brand_score"`
	AnalysisReport string    `gorm:"type:json" json:"analysis_report"`
	CreatedAt      time.Time `json:"created_at"`
}

// Citation 大模型回复中的引用链接明细表
type Citation struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID    uint64    `gorm:"not null" json:"task_id"`
	URL       string    `gorm:"type:text;not null" json:"url"`
	Title     string    `gorm:"type:varchar(255)" json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

// Article 生成的优化品牌曝光的文章表
type Article struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Title          string    `gorm:"type:varchar(255);not null" json:"title"`
	Content        string    `gorm:"type:longtext;not null" json:"content"`
	TargetKeywords string    `gorm:"type:varchar(255)" json:"target_keywords"`
	PublishStatus  string    `gorm:"type:enum('pending','published');default:'pending'" json:"publish_status"`
	PublishedURL   string    `gorm:"type:text" json:"published_url"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
