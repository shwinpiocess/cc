package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Host struct {
	HostID          int       `orm:"column(id);auto"`
	AssetID         string    `orm:"column(asset_id);size(255);null"`
	AutoRenew       int8      `orm:"column(auto_renew);null"`
	BakOperator     int       `orm:"column(bak_operator);null"`
	BandWidth       int       `orm:"column(band_width);null"`
	Cpu             uint      `orm:"column(cpu)"`
	CreateTime      time.Time `orm:"column(create_time);type(timestamp)"`
	Description     string    `orm:"column(description);size(255);null"`
	DeadLineTime    time.Time `orm:"column(deadline_time);type(timestamp);null"`
	DeviceClass     string    `orm:"column(device_class);size(255);null"`
	HardMemo        string    `orm:"column(hard_memo);size(255);null"`
	HostName        string    `orm:"column(host_name);size(255);null"`
	IdcName         string    `orm:"column(idc_name);size(255);null"`
	InnerIP         string    `orm:"column(inner_ip);size(32)"`
	InnerSwitchPort int       `orm:"column(inner_switchport);null"`
	ImageID         int       `orm:"column(image_id);null"`
	LastTime        time.Time `orm:"column(last_time);type(timestamp);null"`
	Mem             int       `orm:"column(mem);null"`
	Operator        string    `orm:"column(operator);size(255);null"`
	OSName          string    `orm:"column(os_name);size(255);null"`
	OuterIP         string    `orm:"column(outer_ip);size(32);null"`
	OuterSwitchPort string    `orm:"column(outer_switchport);size(255);null"`
	PosCode         string    `orm:"column(poscode);size(255);null"`
	Price           float64   `orm:"column(price);null;digits(10);decimals(0)"`
	ProjectId       int       `orm:"column(project_id);null"`
	Region          string    `orm:"column(region);size(255);null"`
	ServerRack      string    `orm:"column(server_rack);size(255);null"`
	ServerType      string    `orm:"column(server_type);size(255);null"`
	SN              int64     `orm:"column(sn);null"`
	Source          int8      `orm:"column(source)"`
	Status          string    `orm:"column(status);size(255);null"`
	StorageId       int       `orm:"column(storage_id);null"`
	StorageSize     float64   `orm:"column(storage_size);null"`
	StorageType     string    `orm:"column(storage_type);size(255);null"`
	Uuid            int       `orm:"column(uuid);null"`
	ZoneID          int       `orm:"column(zone_id);null"`
	ZoneName        string    `orm:"column(zone_name);size(255);null"`
	GseProxy        string    `orm:"column(gse_proxy);size(255);null"`
	VIP             string    `orm:"column(vip);size(255);null"`
	ModName         string    `orm:"column(mod_name);size(255);null"`
	ModuleID        int       `orm:"column(module_id);null"`
	ModuleName      string    `orm:"column(module_name);size(255);null"`
	SetID           int       `orm:"column(set_id);null"`
	SetName         string    `orm:"column(set_name);size(255);null"`
	ApplicationID   int       `orm:"column(application_id);null"`
	ApplicationName string    `orm:"column(application_name);size(255);null"`
	Owner           string    `orm:"column(owner);size(255);null"`
	Checked         string    `orm:"column(checked);size(255);null"`
	IsDistributed   bool      `orm:"column(is_distributed)"`
}

func (t *Host) TableName() string {
	return "host"
}

func init() {
	orm.RegisterModel(new(Host))
}

// AddHost insert a new Host into database and returns
// last inserted Id on success.
func AddHost(hosts []*Host) (err error) {
	o := orm.NewOrm()
	err = o.Begin()

	_, err = o.InsertMulti(len(hosts), hosts)

	if err == nil {
		o.Commit()
	} else {
		o.Rollback()
	}
	return
}

// GetHostById retrieves Host by Id. Returns error if
// Id doesn't exist
func GetHostById(id int) (v *Host, err error) {
	o := orm.NewOrm()
	v = &Host{HostID: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetHostById retrieves Host by Id. Returns error if
// Id doesn't exist
func GetHostBySn(sn int64) bool {
	o := orm.NewOrm()
	return o.QueryTable("host").Filter("sn", sn).Exist()
}

// GetHostById retrieves Host by Id. Returns error if
// Id doesn't exist
func GetHostByInnerIp(inner_ip string) bool {
	o := orm.NewOrm()
	return o.QueryTable("host").Filter("inner_ip", inner_ip).Exist()
}

// GetAllHost retrieves all Host matches certain condition. Returns empty list if
// no records exist
func GetAllHost(query map[string]interface{}, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Host))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Host
	qs = qs.OrderBy(sortFields...)
	if _, err := qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateHost updates Host by Id and returns error if
// the record to be updated doesn't exist
func UpdateHostById(m *Host) (err error) {
	o := orm.NewOrm()
	v := Host{HostID: m.HostID}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// 分配主机
func UpdateHostToApp(ids []int, appID int) (num int64, err error) {
	var set Set
	var mod Module
	o := orm.NewOrm()
	if err = o.QueryTable("set").Filter("ApplicationID", appID).Filter("SetName", "空闲机池").One(&set); err != nil {
		return
	}

	if err = o.QueryTable("module").Filter("ApplicationId", appID).Filter("SetId", set.SetID).One(&mod); err != nil {
		return
	}

	num, err = o.QueryTable("host").Filter("HostID__in", ids).Update(orm.Params{
		"ApplicationID": appID,
		//		"ApplicationName": app
		"SetID": set.SetID,
		//		"SetName": set.SetName,
		"ModuleID":      mod.Id,
		"IsDistributed": true,
	})
	return
}

// 上交主机
func ResHostModule(ids []int, appID int) (num int64, err error) {
	var set Set
	var mod Module
	o := orm.NewOrm()
	if err = o.QueryTable("set").Filter("ApplicationID", appID).Filter("SetName", "空闲机池").One(&set); err != nil {
		return
	}

	if err = o.QueryTable("module").Filter("ApplicationId", appID).Filter("SetId", set.SetID).One(&mod); err != nil {
		return
	}

	num, err = o.QueryTable("host").Filter("HostID__in", ids).Update(orm.Params{
		"ApplicationID": appID,
		//		"ApplicationName": app
		"SetID": set.SetID,
		//		"SetName": set.SetName,
		"ModuleID":      mod.Id,
		"IsDistributed": false,
	})
	return
}

// DeleteHost deletes Host by Id and returns error if
// the record to be deleted doesn't exist
func DeleteHost(id int) (err error) {
	o := orm.NewOrm()
	v := Host{HostID: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Host{HostID: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func DeleteHosts(id []int) (num int64, err error) {
	o := orm.NewOrm()
	num, err = o.QueryTable("host").Filter("HostID__in", id).Delete()
	return
}

func GetHostCount(id int, field string) (cnt int64, err error) {
	o := orm.NewOrm()
	if field == "ApplicationID" {
		cnt, err = o.QueryTable("host").Exclude("ModuleName", "空闲机").Filter(field, id).Count()
	} else {
		cnt, err = o.QueryTable("host").Filter(field, id).Count()
	}
	return
}
