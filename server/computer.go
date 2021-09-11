package server

import (
	"adp_backend/model"

	"gorm.io/gorm"
)

func GetMachineUser(cli *gorm.DB, userID int32) (*model.MachineUser, error) {
	var machine model.MachineUser

	db := cli.Table(machine.TableName()).Where("ID = ?", userID).First(&machine)
	if db.Error != nil {
		return nil, db.Error
	}
	return &machine, nil
}

func FindMachineNameList(cli *gorm.DB, searchName, startTime, endTime string, index, limit, sortTime int32) ([]model.MachineUser, int64, error) {
	var list = []model.MachineUser{}
	tbName := (&model.MachineUser{}).TableName()

	order := "expired_at DESC"
	if sortTime == 1 {
		order = "expired_at ASC"
	}

	db := cli.Table(tbName)
	if startTime != "" {

		db = db.Where("  expired_at > ? and expired_at < ?", startTime, endTime)
	}
	if searchName != "" {
		db = db.Where("machine_name LIKE ?", "%"+searchName+"%")
	}



	var count int64
	db.Count(&count)
	db = db.Limit(int(limit)).Offset(int((index - 1) * limit)).Order(order).Order("machine_name ASC").Find(&list)
	if db.Error != nil {
		return nil, 0, db.Error
	}

	return list, count, nil
}

func FindDomainList(cli *gorm.DB) ([]model.Domain, error) {
	var domainList []model.Domain

	db := cli.Find(&domainList)
	if db.Error != nil {
		return nil, db.Error
	}

	return domainList, nil
}

//UpsetMachineUserByName 更新或新增一条数据(基于计算机名称)
func UpsetMachineUserByName(cli *gorm.DB, machine model.MachineUser) error {

	testName := model.MachineUser{}
	db := cli.Table(machine.TableName()).Where("machine_name = ? and domain = ?", machine.MachineName, machine.Domain).First(&testName)
	if db.Error == nil {
		db := cli.Table(machine.TableName()).Where("machine_name = ? and domain = ?", machine.MachineName, machine.Domain).Updates(&machine) // Update(&machine)
		if db.Error != nil {
			return db.Error
		}
	} else if db.Error == gorm.ErrRecordNotFound {
		db := cli.Table(machine.TableName()).Create(&machine)
		if db.Error != nil {
			return db.Error
		}
	}

	return nil
}

//func UpsetMachineUserByName(cli *gorm.DB, machineList []model.MachineUser) error {
//	tbName := (&model.MachineUser{}).TableName()
//	var newList []model.MachineUser
//	//var oldList []model.MachineUser
//	for _, user := range machineList {
//		testName := model.MachineUser{}
//		db := cli.Table(tbName).Where("machine_name = ? and domain = ?", user.MachineName, user.Domain).First(&testName)
//		if db.Error == gorm.ErrRecordNotFound {
//			newList = append(newList, user)
//		}
//	}
//
//	db := cli.Create(&newList)
//	if db.Error != nil {
//		return db.Error
//	}
//
//	return nil
//}

func UpdateMachineUserRemark(cli *gorm.DB, id int32, remark string) error {
	var user model.MachineUser

	db := cli.Table(user.TableName()).Where("ID = ?", id).Update("remark", remark)
	if db.Error != nil {
		return db.Error
	}

	return nil
}

func GetDomainByName(cli *gorm.DB, domainName string) (*model.Domain, error) {
	var domain model.Domain

	db := cli.Table(domain.TableName()).Where("name = ?", domainName).First(&domain)
	if db.Error != nil {
		return nil, db.Error
	}
	return &domain, db.Error
}

func AddDomain(cli *gorm.DB, domain model.Domain) error {
	db := cli.Table(domain.TableName()).Create(&domain)
	if db.Error != nil {
		return db.Error
	}
	return db.Error
}
