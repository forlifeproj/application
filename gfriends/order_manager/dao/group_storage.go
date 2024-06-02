package dao

const (
	GroupTable = "table_group_info"
)

type GroupStorage struct {
	GroupId      string `gorm:"column:group_id"`      // 群加密id
	GroupName    string `gorm:"column:group_name"`    // 群名称
	GroupAvatar  string `gorm:"column:group_avatar"`  // 头像
	GroupQRCode  string `gorm:"column:group_qr_code"` // 群二维码链接
	GroupMembers int32  `gorm:"column:group_members"` // 群成员数量
	Status       int32  `gorm:"column:status"`        // 状态, 0-已下架，1-已上架
}

func GetAllGroupList() ([]GroupStorage, error) {
	var groups []GroupStorage
	result := gormDb.Table(OrderTable).Where("status = 1").Find(&groups)
	return groups, result.Error
}
