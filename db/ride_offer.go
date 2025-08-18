package db

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

type RideOffer struct {
	ID       string `gorm:"primaryKey;size:191"`
	DriverID string `gorm:"size:191;index"`
	FromGeo  string `gorm:"size:64;index"`
	ToGeo    string `gorm:"size:64;index"`
	Fare     float64
	Time     time.Time `gorm:"index"`
	Seats    int
	Status   string `gorm:"size:32;index"` // active, matched, completed

	Driver *User `gorm:"foreignKey:DriverID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	Matches      []Match       `gorm:"foreignKey:RideID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ChatMessages []ChatMessage `gorm:"foreignKey:RideID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Reviews      []Review      `gorm:"foreignKey:RideID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (o *RideOffer) BeforeCreate(tx *gorm.DB) (err error) {
	if strings.TrimSpace(o.Status) == "" {
		o.Status = "active"
	}
	if o.Time.IsZero() {
		o.Time = time.Now()
	}
	return nil
}

func (o *RideOffer) BeforeSave(tx *gorm.DB) (err error) {
	o.FromGeo = strings.TrimSpace(o.FromGeo)
	o.ToGeo = strings.TrimSpace(o.ToGeo)
	return nil
}

func (o *RideOffer) AfterUpdate(tx *gorm.DB) (err error) {
	if o.Seats == 0 && o.Status == "active" {
		return tx.Model(o).Clauses(clause.Returning{}).
			Update("status", "matched").Error
	}
	return nil
}
