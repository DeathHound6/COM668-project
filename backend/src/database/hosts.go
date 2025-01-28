package database

import (
	"com668-backend/utility"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HostMachine struct {
	ID       uint   `gorm:"column:id;primaryKey;autoIncrement"`
	UUID     string `gorm:"column:uuid;size:36;unique;not null"`
	OS       string `gorm:"column:os;size:20;not null"`
	Hostname string `gorm:"column;hostname;not null"`
	IP4      string `gorm:"column:ip4;size:15"`
	IP6      string `gorm:"column:ip6;size:39"`
	TeamID   uint   `gorm:"column:team_id;not null"`
	Team     Team   `gorm:"foreignKey:team_id;references:id"`
}

func (host *HostMachine) BeforeCreate(tx *gorm.DB) error {
	ctx := GetContext(tx)
	uuid, err := utility.GenerateRandomUUID()
	if err != nil {
		ctx.Set("errorCode", http.StatusInternalServerError)
		return errors.New("failed to create a host uuid")
	}
	if len(host.OS) > 20 {
		ctx.Set("errorCode", http.StatusBadRequest)
		return errors.New("os cannot be greater than 20 characters")
	}
	if len(host.Hostname) > 255 {
		ctx.Set("errorCode", http.StatusBadRequest)
		return errors.New("hostname cannot be greater than 255 characters")
	}
	if len(host.IP4) > 15 {
		ctx.Set("errorCode", http.StatusBadRequest)
		return errors.New("ipv4 address cannot be greater than 15 characters")
	}
	if len(host.IP6) > 39 {
		ctx.Set("errorCode", http.StatusBadRequest)
		return errors.New("ipv6 address cannot be greater than 39 characters")
	}
	host.UUID = uuid
	return nil
}

func (host *HostMachine) BeforeUpdate(tx *gorm.DB) error {
	ctx := GetContext(tx)
	if len(host.OS) > 20 {
		ctx.Set("errorCode", http.StatusBadRequest)
		return errors.New("os cannot be greater than 20 characters")
	}
	if len(host.Hostname) > 255 {
		ctx.Set("errorCode", http.StatusBadRequest)
		return errors.New("hostname cannot be greater than 255 characters")
	}
	if len(host.IP4) > 15 {
		ctx.Set("errorCode", http.StatusBadRequest)
		return errors.New("ipv4 address cannot be greater than 15 characters")
	}
	if len(host.IP6) > 39 {
		ctx.Set("errorCode", http.StatusBadRequest)
		return errors.New("ipv6 address cannot be greater than 39 characters")
	}
	return nil
}

type GetHostsFilters struct {
	Page     *int
	PageSize *int
}

func GetHosts(ctx *gin.Context, filters GetHostsFilters) ([]*HostMachine, int64, error) {
	tx := GetDBTransaction(ctx).Model(&HostMachine{})

	var count int64
	tx.Count(&count)
	if filters.PageSize != nil {
		tx = tx.Limit(*filters.PageSize)
		if filters.Page != nil {
			tx = tx.Offset(*filters.PageSize * (*filters.Page - 1))
		}
	}

	tx = tx.Preload("Team")
	var hosts []*HostMachine
	tx = tx.Find(&hosts)
	if tx.Error != nil {
		ctx.Set("errorCode", http.StatusInternalServerError)
		return nil, -1, tx.Error
	}
	return hosts, count, nil
}
