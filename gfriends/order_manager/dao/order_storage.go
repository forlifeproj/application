package dao

import (
	"time"
)

const (
	OrderTable = "table_order_record"
)

type OrderStorage struct {
	OrderId    string `gorm:"column:order_id"`     // 商户订单号
	Openid     string `gorm:"column:openid"`       // 微信openid
	PayOrderid string `gorm:"column:pay_order_id"` // 支付订单号
	GroupID    string `gorm:"column:group_id"`     // 群ID
	Amount     int64  `gorm:"column:amount"`       // 金额
	Status     int    `gorm:"column:status"`
	CreateTs   int    `gorm:"column:create_ts"`
	UpdateTs   int    `gorm:"column:update_ts"`
}

func AddOrder(record *OrderStorage) error {
	result := gormDb.Table(OrderTable).Where("order_id = ?", record.OrderId).Save(record)
	return result.Error
}

func UpdateOrderRecord(orderId string, status int, payOrderId string) error {
	result := gormDb.Table(OrderTable).Model(&OrderStorage{}).Where("order_id = ?", orderId).Updates(&OrderStorage{Status: status, PayOrderid: payOrderId, UpdateTs: int(time.Now().Unix())})
	return result.Error
}

func GetOrderRecord(orderId string) (*OrderStorage, error) {
	var order OrderStorage
	result := gormDb.Table(OrderTable).Where("order_id = ?", orderId).First(&order)
	return &order, result.Error
}

func GetOrderList(uid int64, lastOrderId string, pageSize int) ([]OrderStorage, error) {
	var orders []OrderStorage
	result := gormDb.Table(OrderTable).Where("uid = ? and order_id > ?", uid, lastOrderId).Limit(pageSize).Find(&orders)
	return orders, result.Error
}
