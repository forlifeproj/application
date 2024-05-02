package dao

const (
	OrderTable = "table_order_record"
)

type OrderStorage struct {
	OrderId   string `gorm:"column:order_id"`
	Uid       int64  `gorm:"column:uid"`
	ProductID string `gorm:"column:product_id"`
	Status    int    `gorm:"column:status"`
	Ts        int    `gorm:"column:ts"`
}

func AddOrder(record *OrderStorage) error {
	result := gormDb.Table(OrderTable).Save(record)
	return result.Error
}

func GetOrderRecord(orderId string) (*OrderStorage, error) {
	var order OrderStorage
	result := gormDb.Table(OrderTable).First(&order).Where("order_id = ?", orderId)
	return &order, result.Error
}

func GetOrderList(uid int64, lastOrderId string, pageSize int) ([]OrderStorage, error) {
	var orders []OrderStorage
	result := gormDb.Table(OrderTable).Where("uid = ? and order_id > ?", uid, lastOrderId).Limit(pageSize).Find(&orders)
	return orders, result.Error
}
