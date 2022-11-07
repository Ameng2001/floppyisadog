package models

import "time"

// BaseGormModel mimixks GormModel but uses uuid's for ID, generated in go
type BaseGormModel struct {
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Company model
type Company struct {
	BaseGormModel
	Name                 string `sql:"type:varchar(255);not null"`
	Archived             bool   `sql:"default:false"`
	DefaultTimezone      string `sql:"type:varchar(255);not null"`
	DefaultDayWeekStarts string `sql:"type:varchar(20);not null;defalut 'Monday'"`
}

func (c *Company) TableName() string {
	return "companies"
}

// Dirictory model
type Directory struct {
	BaseGormModel
	CompanyId  string `sql:"type:varchar(255);not null;unique"`
	UserId     string `sql:"type:varchar(255);not null;unique"`
	InternalId string `sql:"type:varchar(255);not null;unique"`
}

func (c *Directory) TableName() string {
	return "directories"
}

// Admin model
type Admin struct {
	BaseGormModel
	CompanyId string `sql:"type:varchar(255);not null;unique"`
	UserId    string `sql:"type:varchar(255);not null;unique"`
}

func (c *Admin) TableName() string {
	return "admins"
}

// Team model
type Team struct {
	BaseGormModel
	CompanyId     string `sql:"type:varchar(255);not null;unique"`
	Name          string `sql:"type:varchar(255);not null"`
	Archived      bool   `sql:"default:false"`
	Timezone      string `sql:"type:varchar(255);not null"`
	DayWeekStarts string `sql:"type:varchar(20);not null;defalut:'Monday'"`
	Color         string `sql:"type:varchar(10);not null;default:'#48B7AB'"`
}

func (c *Team) TableName() string {
	return "teams"
}

// Worker model
type Worker struct {
	BaseGormModel
	TeamId string `sql:"type:varchar(255);not null;unique"`
	UserId string `sql:"type:varchar(255);not null;unique"`
}

func (c *Worker) TableName() string {
	return "workers"
}

//Job Model
type Job struct {
	BaseGormModel
	TeamId   string `sql:"type:varchar(255);not null;unique"`
	Name     string `sql:"type:varchar(255);not null"`
	Archived bool   `sql:"default:false"`
	Color    string `sql:"type:varchar(10);not null;default:'#48B7AB'"`
}

func (c *Job) TableName() string {
	return "jobs"
}

// Shift Model
type Shift struct {
	BaseGormModel
	TeamId    string    `sql:"type:varchar(255);not null"`
	JobId     string    `sql:"type:varchar(255);not null;unique"`
	UserId    string    `sql:"type:varchar(255);not null;unique"`
	Published bool      `sql:"default:false;not null"`
	Start     time.Time `sql:"not null;DEFAULT:current_timestamp"`
	Stop      time.Time `sql:"not null;DEFAULT:current_timestamp"`
}

func (c *Shift) TableName() string {
	return "shifts"
}
