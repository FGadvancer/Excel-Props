package db

import "time"

func BatchInsertIntoGroupMember(toInsertInfoList []*db.GroupMember) error {
	for _, toInsertInfo := range toInsertInfoList {
		toInsertInfo.JoinTime = time.Now()
		if toInsertInfo.RoleLevel == 0 {
			toInsertInfo.RoleLevel = constant.GroupOrdinaryUsers
		}
		toInsertInfo.MuteEndTime = time.Unix(int64(time.Now().Second()), 0)
	}
	return db.DB.MysqlDB.DefaultGormDB().Create(toInsertInfoList).Error

}
