package service

import (
	"github.com/forlifeproj/application/gfriends/order_manager/conf"
	fllog "github.com/forlifeproj/msf/log"
	"github.com/forlifeproj/protocol/gfriends/json/order_manager"
)

var groupList []order_manager.GroupInfo

func InitGroupList() {
	// friends
	groupList = append(groupList, order_manager.GroupInfo{
		GroupId:     "c8b963bb7924e54859e3450fce5ec319",
		Price:       990,
		GroupQRCode: conf.GConf.GroupConf.FriendGroupUrl,
	})
	groupList = append(groupList, order_manager.GroupInfo{
		GroupId:     "830e818a1abb1b0c684cc50523cf7632",
		Price:       1990,
		GroupQRCode: conf.GConf.GroupConf.FriendGroupUrl,
	})
	groupList = append(groupList, order_manager.GroupInfo{
		GroupId:     "88b1b369b60f163c76072261a252f0c4",
		Price:       2990,
		GroupQRCode: conf.GConf.GroupConf.FriendGroupUrl,
	})
	// news
	groupList = append(groupList, order_manager.GroupInfo{
		GroupId:     "b00e5c9873c2a37766a94d11c19131e8",
		Price:       990,
		GroupQRCode: conf.GConf.GroupConf.NewsGroupUrl,
	})
	groupList = append(groupList, order_manager.GroupInfo{
		GroupId:     "04f886bf57f9070935c4a0f972032380",
		Price:       1990,
		GroupQRCode: conf.GConf.GroupConf.NewsGroupUrl,
	})
	groupList = append(groupList, order_manager.GroupInfo{
		GroupId:     "05959e4eb153d8bf8ad8d775fdddf08a",
		Price:       2990,
		GroupQRCode: conf.GConf.GroupConf.NewsGroupUrl,
	})
	fllog.Log().Debugf("groups:%+v", groupList)
}

func GetGroupInfoById(groupId string) *order_manager.GroupInfo {
	for _, item := range groupList {
		if groupId == item.GroupId {
			return &item
		}
	}
	return nil
}
