package database

import (
	"com668-backend/utility"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HostMachine struct {
	ID       uint    `gorm:"column:id;primaryKey;autoIncrement"`
	UUID     string  `gorm:"column:uuid;size:36;unique;not null"`
	OS       string  `gorm:"column:os;size:20;not null;check:os IN ('Windows','Linux','MacOS')"`
	Hostname string  `gorm:"column:hostname;size:40;not null;unique"`
	IP4      *string `gorm:"column:ip4;size:15;unique"`
	IP6      *string `gorm:"column:ip6;size:39;unique"`
	TeamID   uint    `gorm:"column:team_id;not null"`
	Team     Team    `gorm:"foreignKey:team_id;references:id"`
}

func (host *HostMachine) BeforeCreate(tx *gorm.DB) error {
	ctx := GetContext(tx)
	if host.UUID == "" {
		uuid, err := utility.GenerateRandomUUID()
		if err != nil {
			if ctx != nil {
				ctx.Set("errorCode", http.StatusInternalServerError)
			}
			return errors.New("failed to create a host uuid")
		}
		host.UUID = uuid
	}
	return nil
}

type GetHostsFilters struct {
	UUIDs    []string
	Page     *int
	PageSize *int
	Hostname *string
}

func GetHost(ctx *gin.Context, filters GetHostsFilters) (*HostMachine, error) {
	hosts, count, err := GetHosts(ctx, GetHostsFilters{
		UUIDs:    filters.UUIDs,
		PageSize: utility.Pointer(len(filters.UUIDs)),
	})
	if err != nil {
		return nil, err
	}
	if count == 0 {
		ctx.Set("errorCode", http.StatusNotFound)
		return nil, errors.New("host not found")
	}
	return hosts[0], nil
}

func GetHosts(ctx *gin.Context, filters GetHostsFilters) ([]*HostMachine, int64, error) {
	tx := GetDBTransaction(ctx).Model(&HostMachine{})
	tx = tx.Preload("Team").Preload("Team.Users")

	if len(filters.UUIDs) > 0 {
		tx = tx.Where("uuid IN (?)", filters.UUIDs)
	}
	if filters.Hostname != nil {
		tx = tx.Where("hostname = ?", *filters.Hostname)
	}

	var count int64
	tx.Count(&count)
	if filters.PageSize != nil {
		tx = tx.Limit(*filters.PageSize)
		if filters.Page != nil {
			tx = tx.Offset(*filters.PageSize * (*filters.Page - 1))
		}
	}

	var hosts []*HostMachine
	tx = tx.Find(&hosts)
	if tx.Error != nil {
		return nil, -1, handleError(ctx, tx.Error)
	}
	return hosts, count, nil
}

func CreateHost(ctx *gin.Context, host *HostMachine) error {
	tx := GetDBTransaction(ctx).Model(&HostMachine{})
	tx = tx.Create(host)
	if tx.Error != nil {
		return handleError(ctx, tx.Error)
	}
	return nil
}

func UpdateHost(ctx *gin.Context, host *HostMachine) error {
	tx := GetDBTransaction(ctx).Model(&HostMachine{})
	tx = tx.Where("uuid = ?", host.UUID)
	fields := map[string]any{"os": host.OS, "hostname": host.Hostname, "ip4": host.IP4, "ip6": host.IP6, "team_id": host.TeamID}
	tx = tx.Updates(fields)
	if tx.Error != nil {
		return handleError(ctx, tx.Error)
	}
	return nil
}

func DeleteHost(ctx *gin.Context, uuid string) error {
	tx := GetDBTransaction(ctx)
	if err := tx.Model(&IncidentHost{}).
		Where("host_machine_id IN (?)", tx.Table("tbl_host_machine").Select("id").Where("uuid = ?", uuid)).
		Delete(&IncidentHost{}).Error; err != nil {
		return handleError(ctx, err)
	}
	if err := tx.Model(&HostMachine{}).Where("uuid = ?", uuid).Delete(&HostMachine{}).Error; err != nil {
		return handleError(ctx, err)
	}
	return nil
}
