package controllers

import (
	"errors"
	"net/http"
	"products/errs"
	"products/models"
	"products/pkg/service"
	"strconv"
	"products/logger"
	"github.com/gin-gonic/gin"
)

// CreateOrder
// @Summary Create Order
// @Security ApiKeyAuth
// @Tags orders
// @Description create new order
// @ID create-order
// @Accept json
// @Produce json
// @Param input body models.Order true "new order info"
// @Success 200 {object} defaultResponse
// @Failure 400 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /api/order [post]
func CreateOrder(c *gin.Context) {
	// userID := c.GetUint(userIDCtx)
	// urole := c.GetString(userRoleCtx)
	// if urole == "" {
	// 	handleError(c, errs.ErrValidationFailed)
	// 	return
	// }
	userID, exists := c.Get(userIDCtx)
	if !exists {
		handleError(c, errs.ErrUnauthorized)
		return
	}

	// Приводим userID к uint
	uid, ok := userID.(uint)
	if !ok {
		handleError(c, errs.ErrValidationFailed)
		return
	}
	var orderItems []models.OrderItem
	if err := c.BindJSON(&orderItems); err != nil {
		logger.Error.Printf("[controllers.AddOrder] error binding order items: %v\n", err)
		handleError(c, errs.ErrValidationFailed)
		return
	}

	logger.Info.Printf("User ID [%d] is creating an order\n", uid)

	// Создаем заказ через сервис
	order, err := service.CreateOrder(uid, orderItems)
	if err != nil {
		handleError(c, err)
		return
	}

	logger.Info.Printf("User ID [%d] successfully created an order with ID: %d\n", uid, order.ID)
	// c.JSON(http.StatusCreated, gin.H{"message": "Order created successfully!!!"})
	// Возвращаем сообщение и ID заказа
	c.JSON(http.StatusCreated, gin.H{
		"message":  "Order created successfully!!!",
		"order_id": order.ID,
	})
}

// GetUserOrderByID
// @Summary Get User Order By ID
// @Security ApiKeyAuth
// @Tags orders
// @Description get user order by ID
// @ID get-user-order-by-ID
// @Produce json
// @Param id path integer true "id of the order"
// @Success 200 {array} models.Order
// @Failure 400 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /api/order/{id} [get]
func GetUserOrderByID(c *gin.Context) {
	userID := c.GetUint(userIDCtx)
	if userID == 0 {
		handleError(c, errs.ErrRecordNotFound)
		return
	}
	urole := c.GetString(userRoleCtx)
	if urole == "" {
		handleError(c, errs.ErrValidationFailed)
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handleError(c, errs.ErrValidationFailed)
		return
	}
	orders, err := service.GetUserOrdersByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetAllOrders
// @Summary Get All Orders
// @Security ApiKeyAuth
// @Tags orders
// @Description get list of all orders
// @ID get-all-orders
// @Produce json
// @Param q query string false "fill if you need search"
// @Success 200 {array} models.Order
// @Failure 400 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /api/order [get]
func GetAllOrders(c *gin.Context) {
	userID := c.GetUint(userIDCtx)
	if userID == 0 {
		handleError(c, errs.ErrRecordNotFound)
		return
	}
	urole := c.GetString(userRoleCtx)
	if urole == "" {
		handleError(c, errs.ErrValidationFailed)
		return
	}
	if urole != "admin" {
		handleError(c, errs.ErrPermissionDenied)
		return
	}

	orders, err := service.GetAllUserOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// RemoveItem
// @Summary Delete Order By ID
// @Security ApiKeyAuth
// @Tags orders
// @Description delete order by ID
// @ID delete-order-by-ID
// @Param id path integer true "id of the order"
// @Success 200 {object} defaultResponse
// @Failure 400 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /api/order/{id} [delete]
func RemoveItem(c *gin.Context) {
	orderItemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.RemoveItem(uint(orderItemID)); err != nil {
		if errors.Is(err, errs.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
