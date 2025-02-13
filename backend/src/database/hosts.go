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
	uuid, err := utility.GenerateRandomUUID()
	if err != nil {
		if ctx != nil {
			ctx.Set("errorCode", http.StatusInternalServerError)
		}
		return errors.New("failed to create a host uuid")
	}
	if !utility.SliceHasElement([]string{"Windows", "Linux", "MacOS"}, host.OS) {
		if ctx != nil {
			ctx.Set("errorCode", http.StatusBadRequest)
		}
		return errors.New("os must be one of 'Windows', 'Linux', or 'MacOS'")
	}
	if len(host.Hostname) > 255 {
		if ctx != nil {
			ctx.Set("errorCode", http.StatusBadRequest)
		}
		return errors.New("hostname cannot be greater than 255 characters")
	}
	if ip4 := host.IP4; ip4 != nil {
		if *ip4 == "" {
			host.IP4 = nil
		} else if len(*ip4) > 15 {
			if ctx != nil {
				ctx.Set("errorCode", http.StatusBadRequest)
			}
			return errors.New("ipv4 address cannot be greater than 15 characters")
		}
	}
	if ip6 := host.IP6; ip6 != nil {
		if *ip6 == "" {
			host.IP6 = nil
		} else if len(*ip6) > 39 {
			if ctx != nil {
				ctx.Set("errorCode", http.StatusBadRequest)
			}
			return errors.New("ipv6 address cannot be greater than 39 characters")
		}
	}
	if host.IP4 == nil && host.IP6 == nil {
		if ctx != nil {
			ctx.Set("errorCode", http.StatusBadRequest)
		}
		return errors.New("host must contain either IPv4 or IPv6 address")
	}
	host.UUID = uuid
	return nil
}

func (host *HostMachine) BeforeUpdate(tx *gorm.DB) error {
	ctx := GetContext(tx)
	if !utility.SliceHasElement([]string{"Windows", "Linux", "MacOS"}, host.OS) {
		if ctx != nil {
			ctx.Set("errorCode", http.StatusBadRequest)
		}
		return errors.New("os must be one of 'Windows', 'Linux', or 'MacOS'")
	}
	if len(host.Hostname) > 255 {
		if ctx != nil {
			ctx.Set("errorCode", http.StatusBadRequest)
		}
		return errors.New("hostname cannot be greater than 255 characters")
	}
	if ip4 := host.IP4; ip4 != nil {
		if *ip4 == "" {
			host.IP4 = nil
		} else if len(*ip4) > 15 {
			if ctx != nil {
				ctx.Set("errorCode", http.StatusBadRequest)
			}
			return errors.New("ipv4 address cannot be greater than 15 characters")
		}
	}
	if ip6 := host.IP6; ip6 != nil {
		if *ip6 == "" {
			host.IP6 = nil
		} else if len(*ip6) > 39 {
			if ctx != nil {
				ctx.Set("errorCode", http.StatusBadRequest)
			}
			return errors.New("ipv6 address cannot be greater than 39 characters")
		}
	}
	if host.IP4 == nil && host.IP6 == nil {
		if ctx != nil {
			ctx.Set("errorCode", http.StatusBadRequest)
		}
		return errors.New("host must contain either IPv4 or IPv6 address")
	}
	return nil
}

type GetHostsFilters struct {
	UUIDs    []string
	Page     *int
	PageSize *int
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
	tx = tx.Preload("Team")

	if len(filters.UUIDs) > 0 {
		tx = tx.Where("uuid IN (?)", filters.UUIDs)
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
	tx := GetDBTransaction(ctx).Model(&HostMachine{})
	tx = tx.Where("uuid = ?", uuid).Delete(&HostMachine{})
	if tx.Error != nil {
		return handleError(ctx, tx.Error)
	}
	return nil
}
