package dao

const (
	LoginTokenTable = "table_login_token"
)

type LoginToken struct {
	Token              string `gorm:"column:token"`
	Uid                int64  `gorm:"column:uid"`
	OpenId             string `gorm:"column:openid"`
	UnionId            string `gorm:"column:unionid"`
	LoginType          int    `gorm:"column:login_type"`
	LoginScene         int    `gorm:"column:login_scene"`
	Rgister            int    `gorm:"column:register"`
	AccessToken        string `gorm:"column:access_token"`
	RefreshToken       string `gorm:"column:refresh_token"`
	AccessTokenExpire  int64  `gorm:"column:access_token_expire"`
	RefreshTokenExpire int64  `gorm:"column:refresh_token_expire"`
	Status             int    `gorm:"column:status"`
	ExpireTime         int64  `gorm:"column:expire_time"`
	CreateTIme         int64  `gorm:"column:create_time"`
}

func AddLoginToken(record *LoginToken) error {
	result := gormDb.Table(LoginTokenTable).Save(record)
	return result.Error
}

func GetLoginToken(token string) (*LoginToken, error) {
	var loginToken LoginToken
	result := gormDb.Table(LoginTokenTable).First(&loginToken).Where("token = ?", token)
	return &order, result.Error
}

// func GetOrderList(uid int64, lastOrderId string, pageSize int) ([]OrderStorage, error) {
// 	var orders []OrderStorage
// 	result := gormDb.Table(OrderTable).Where("uid = ? and order_id > ?", uid, lastOrderId).Limit(pageSize).Find(&orders)
// 	return orders, result.Error
// }
