package Order

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}
func (h *Handler) Serve(orders chi.Router) {
	orders.Post("/", h.AddOrder)
	orders.Get("/{orderID}", h.FetchOrder)
}
func (h *Handler) AddOrder(w http.ResponseWriter, r *http.Request) {
	var order AddOrderModel
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		log.Printf("AddOrderHandler: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	logrus.Info("AddOrder: ", order.UserID)
	err, orderID := h.service.AddOrderHelper(order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte(orderID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}

func (h *Handler) FetchOrder(w http.ResponseWriter, r *http.Request) {
	var order string
	order = chi.URLParam(r, "orderID")
	logrus.Info("FetchOrder: ", order)
	err, orderDetails := h.service.FetchOrderHelper(order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonData, err := json.Marshal(orderDetails)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}
