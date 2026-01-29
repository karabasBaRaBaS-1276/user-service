package handler

import (
	"net/http"

	"github.com/karabasBaRaBaS-1276/user-service/internal/service"
)

type PersonHandler struct {
	service service.PersonService // указатель на сервис, реализующий бизнес логику
}

func NewPersonHandler(s service.PersonService) *PersonHandler {
	return &PersonHandler{service: s}
}

// Обработка http запросов по работе с пользователем
func (h *PersonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.create(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *PersonHandler) create(w http.ResponseWriter, r *http.Request) {
	// TODO: decode request
	// TODO: call h.service.CreatePerson
	w.WriteHeader(http.StatusCreated)
}
