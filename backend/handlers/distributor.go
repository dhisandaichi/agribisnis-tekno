package handlers

import (
	"net/http"

	"agritec/backend/middleware"

	"github.com/gin-gonic/gin"
)

func (h *H) myDistributor(c *gin.Context) (*DistributorRow, bool) {
	var d DistributorRow
	res := h.DB.Raw(distributorSelect+`where d.owner_id = ? group by d.id, u.full_name, u.email`, middleware.User(c).ID).Scan(&d)
	if res.Error != nil {
		serverError(c, res.Error)
		return nil, false
	}
	if res.RowsAffected == 0 {
		return nil, true
	}
	return &d, true
}

func (h *H) DistributorOverview(c *gin.Context) {
	d, ok := h.myDistributor(c)
	if !ok {
		return
	}
	if d == nil {
		c.JSON(http.StatusOK, gin.H{"distributor": nil, "orders": []OrderRow{}})
		return
	}
	var orders []OrderRow
	if err := h.DB.Raw(orderSelect+`where o.distributor_id = ? order by o.order_date desc, o.id desc`, d.ID).Scan(&orders).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"distributor": d, "orders": rows(orders)})
}

// Distributor mengajukan order; nilai (total) dikonfirmasi supervisor.
func (h *H) DistributorCreateOrder(c *gin.Context) {
	var in struct {
		Product  string `json:"product" binding:"required,max=255"`
		Quantity int    `json:"quantity" binding:"required,gt=0,lte=100000"`
	}
	if !bindJSON(c, &in) {
		return
	}
	d, ok := h.myDistributor(c)
	if !ok {
		return
	}
	if d == nil {
		fail(c, http.StatusForbidden, "Akun Anda belum terhubung ke data kios")
		return
	}
	err := h.DB.Exec(`insert into orders (distributor_id, created_by, product, quantity, total_amount, order_date, due_date)
		values (?, ?, ?, ?, 0, current_date, current_date + ?::int)`, d.ID, middleware.User(c).ID, in.Product, in.Quantity, d.CreditLimitDays).Error
	if err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Order terkirim"})
}
