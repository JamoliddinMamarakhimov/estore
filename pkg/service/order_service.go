package service

import (
	"products/models"
	"products/pkg/repository"
)

func CreateOrder(userID uint, orderItems []models.OrderItem) (models.Order, error) {
	if err := repository.CreateOrder(&models.Order{}); err != nil {
		return models.Order{}, err
	}

	return models.Order{}, nil
}

func GetUserOrdersByID(id uint) ([]models.Order, error) {
	var orders []models.Order
	orders, err := repository.GetOrdersByUserID(id)
	if err != nil {
		return orders, err
	}
	return orders, nil
}

func GetAllUserOrders() ([]models.Order, error) {
	orders, err := repository.GetAllOrders()
	if err != nil {
		return orders, err
	}
	return orders, nil
}

func UpdateItem(items models.Order, id, userID uint) error {
	err := repository.UpdateOrder(items.ProductID, items.Quantity, id, userID)
	if err != nil {
		return err
	}
	return nil
}

func RemoveItem(oID uint) error {
	err := repository.DeleteOrder(oID)
	if err != nil {
		return err
	}
	return nil
}
