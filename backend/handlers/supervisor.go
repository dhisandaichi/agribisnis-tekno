package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"agritec/backend/middleware"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (h *H) Summary(c *gin.Context) {
	var s struct {
		Farmers        int64   `json:"farmers"`
		Visits7d       int64   `json:"visits_7d"`
		VisitsToday    int64   `json:"visits_today"`
		Distributors   int64   `json:"distributors"`
		OverdueCount   int64   `json:"overdue_distributors"`
		OverdueAmount  float64 `json:"overdue_amount"`
		Outstanding    float64 `json:"outstanding_amount"`
		PendingOrders  int64   `json:"pending_orders"`
	}
	err := h.DB.Raw(`select
		(select count(*) from farmers) as farmers,
		(select count(*) from visits where visit_time >= now() - interval '7 days') as visits7d,
		(select count(*) from visits where (visit_time at time zone 'Asia/Jakarta')::date = (now() at time zone 'Asia/Jakarta')::date) as visits_today,
		(select count(*) from distributors) as distributors,
		(select count(distinct distributor_id) from orders where payment_status = 'unpaid' and due_date < current_date) as overdue_count,
		(select coalesce(sum(total_amount),0)::float8 from orders where payment_status = 'unpaid' and due_date < current_date) as overdue_amount,
		(select coalesce(sum(total_amount),0)::float8 from orders where payment_status = 'unpaid') as outstanding,
		(select count(*) from orders where total_amount = 0 and payment_status = 'unpaid') as pending_orders`).Scan(&s).Error
	if err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *H) AllVisits(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days < 1 || days > 3650 {
		days = 30
	}
	var out []VisitRow
	err := h.DB.Raw(visitSelect+`where v.visit_time >= (date_trunc('day', now() at time zone 'Asia/Jakarta') - make_interval(days => ?)) at time zone 'Asia/Jakarta'
		order by v.visit_time desc limit 2000`, days-1).Scan(&out).Error
	if err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, rows(out))
}

// ---------- Distributor ----------
type DistributorRow struct {
	ID              int64    `json:"id"`
	Name            string   `json:"name"`
	OwnerID         *int64   `json:"owner_id"`
	OwnerName       *string  `json:"owner_name"`
	Address         *string  `json:"address"`
	Phone           *string  `json:"phone"`
	CreditLimitDays int      `json:"credit_limit_days"`
	Outstanding     float64  `json:"outstanding"`
	Overdue         float64  `json:"overdue"`
	DaysLate        *int     `json:"days_late"`
}

const distributorSelect = `select d.id, d.name, d.owner_id, coalesce(nullif(u.full_name,''), u.email) as owner_name, d.address, d.phone, d.credit_limit_days,
	coalesce(sum(o.total_amount) filter (where o.payment_status = 'unpaid'), 0)::float8 as outstanding,
	coalesce(sum(o.total_amount) filter (where o.payment_status = 'unpaid' and o.due_date < current_date), 0)::float8 as overdue,
	max(current_date - o.due_date) filter (where o.payment_status = 'unpaid' and o.due_date < current_date) as days_late
	from distributors d left join users u on u.id = d.owner_id left join orders o on o.distributor_id = d.id `

func (h *H) ListDistributors(c *gin.Context) {
	var out []DistributorRow
	if err := h.DB.Raw(distributorSelect + `group by d.id, u.full_name, u.email order by d.name`).Scan(&out).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, rows(out))
}

type distributorInput struct {
	Name            string `json:"name" binding:"required,max=255"`
	OwnerID         *int64 `json:"owner_id"`
	Address         string `json:"address"`
	Phone           string `json:"phone" binding:"max=50"`
	CreditLimitDays int    `json:"credit_limit_days" binding:"gte=0,lte=365"`
}

func (h *H) checkOwner(c *gin.Context, ownerID *int64) bool {
	if ownerID == nil {
		return true
	}
	var role string
	h.DB.Raw(`select role from users where id = ?`, *ownerID).Scan(&role)
	if role != "distributor" {
		fail(c, http.StatusBadRequest, "Akun pemilik harus ber-role distributor")
		return false
	}
	return true
}

func (h *H) SaveDistributor(c *gin.Context) {
	var in distributorInput
	if !bindJSON(c, &in) || !h.checkOwner(c, in.OwnerID) {
		return
	}
	var err error
	if c.Param("id") == "" {
		err = h.DB.Exec(`insert into distributors (name, owner_id, address, phone, credit_limit_days) values (?, ?, nullif(?, ''), nullif(?, ''), ?)`,
			in.Name, in.OwnerID, in.Address, in.Phone, in.CreditLimitDays).Error
	} else {
		id, ok := idParam(c)
		if !ok {
			return
		}
		err = h.DB.Exec(`update distributors set name = ?, owner_id = ?, address = nullif(?, ''), phone = nullif(?, ''), credit_limit_days = ? where id = ?`,
			in.Name, in.OwnerID, in.Address, in.Phone, in.CreditLimitDays, id).Error
	}
	if err != nil {
		if strings.Contains(err.Error(), "distributors_owner_id_key") {
			fail(c, http.StatusConflict, "Akun tersebut sudah terhubung ke kios lain")
			return
		}
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tersimpan"})
}

// ---------- Orders ----------
type OrderRow struct {
	ID              int64     `json:"id"`
	DistributorID   int64     `json:"distributor_id"`
	DistributorName string    `json:"distributor_name"`
	Product         string    `json:"product"`
	Quantity        int       `json:"quantity"`
	TotalAmount     float64   `json:"total_amount"`
	OrderDate       time.Time `json:"order_date"`
	DueDate         *time.Time `json:"due_date"`
	PaymentStatus   string    `json:"payment_status"`
}

const orderSelect = `select o.id, o.distributor_id, d.name as distributor_name, o.product, o.quantity, o.total_amount::float8 as total_amount,
	o.order_date, o.due_date, o.payment_status from orders o join distributors d on d.id = o.distributor_id `

func (h *H) ListOrders(c *gin.Context) {
	var out []OrderRow
	if err := h.DB.Raw(orderSelect + `order by o.order_date desc, o.id desc limit 1000`).Scan(&out).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, rows(out))
}

type orderInput struct {
	DistributorID int64   `json:"distributor_id" binding:"required"`
	Product       string  `json:"product" binding:"required,max=255"`
	Quantity      int     `json:"quantity" binding:"required,gt=0"`
	TotalAmount   float64 `json:"total_amount" binding:"gte=0"`
	OrderDate     string  `json:"order_date" binding:"required"`
	DueDate       string  `json:"due_date"`
}

func (h *H) SaveOrder(c *gin.Context) {
	var in orderInput
	if !bindJSON(c, &in) {
		return
	}
	if _, err := time.Parse("2006-01-02", in.OrderDate); err != nil {
		fail(c, http.StatusBadRequest, "Tanggal order tidak valid")
		return
	}
	if in.DueDate != "" {
		if _, err := time.Parse("2006-01-02", in.DueDate); err != nil {
			fail(c, http.StatusBadRequest, "Tanggal jatuh tempo tidak valid")
			return
		}
	}
	// jatuh tempo kosong = tanggal order + tempo distributor
	due := `coalesce(nullif(?, '')::date, ?::date + (select credit_limit_days from distributors where id = ?))`
	var err error
	if c.Param("id") == "" {
		err = h.DB.Exec(`insert into orders (distributor_id, created_by, product, quantity, total_amount, order_date, due_date)
			values (?, ?, ?, ?, ?, ?, `+due+`)`,
			in.DistributorID, middleware.User(c).ID, in.Product, in.Quantity, in.TotalAmount, in.OrderDate, in.DueDate, in.OrderDate, in.DistributorID).Error
	} else {
		id, ok := idParam(c)
		if !ok {
			return
		}
		err = h.DB.Exec(`update orders set distributor_id = ?, product = ?, quantity = ?, total_amount = ?, order_date = ?, due_date = `+due+` where id = ?`,
			in.DistributorID, in.Product, in.Quantity, in.TotalAmount, in.OrderDate, in.DueDate, in.OrderDate, in.DistributorID, id).Error
	}
	if err != nil {
		if strings.Contains(err.Error(), "foreign key") {
			fail(c, http.StatusBadRequest, "Distributor tidak ditemukan")
			return
		}
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tersimpan"})
}

func (h *H) SetOrderStatus(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var in struct {
		PaymentStatus string `json:"payment_status" binding:"required,oneof=paid unpaid"`
	}
	if !bindJSON(c, &in) {
		return
	}
	if err := h.DB.Exec(`update orders set payment_status = ? where id = ?`, in.PaymentStatus, id).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Status diperbarui"})
}

func (h *H) DeleteOrder(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	if err := h.DB.Exec(`delete from orders where id = ?`, id).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Dihapus"})
}

// ---------- Users ----------
func (h *H) ListUsers(c *gin.Context) {
	var out []struct {
		ID        int64     `json:"id"`
		Email     string    `json:"email"`
		FullName  *string   `json:"full_name"`
		Phone     *string   `json:"phone"`
		Role      string    `json:"role"`
		Active    bool      `json:"active"`
		CreatedAt time.Time `json:"created_at"`
	}
	if err := h.DB.Raw(`select id, email, full_name, phone, role, active, created_at from users order by created_at desc`).Scan(&out).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, rows(out))
}

func (h *H) UpdateUser(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	var in struct {
		Role   *string `json:"role" binding:"omitempty,oneof=sales supervisor farmer distributor"`
		Active *bool   `json:"active"`
	}
	if !bindJSON(c, &in) {
		return
	}
	if id == middleware.User(c).ID {
		fail(c, http.StatusBadRequest, "Tidak bisa mengubah role/status akun sendiri")
		return
	}
	if err := h.DB.Exec(`update users set role = coalesce(?, role), active = coalesce(?, active) where id = ?`, in.Role, in.Active, id).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Pengguna diperbarui"})
}

// ResetPassword membuat password sementara untuk user yang lupa password.
func (h *H) ResetUserPassword(c *gin.Context) {
	id, ok := idParam(c)
	if !ok {
		return
	}
	pw := randomPassword()
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		serverError(c, err)
		return
	}
	res := h.DB.Exec(`update users set password_hash = ? where id = ?`, string(hash), id)
	if res.Error != nil {
		serverError(c, res.Error)
		return
	}
	if res.RowsAffected == 0 {
		fail(c, http.StatusNotFound, "Pengguna tidak ditemukan")
		return
	}
	c.JSON(http.StatusOK, gin.H{"temporary_password": pw})
}
