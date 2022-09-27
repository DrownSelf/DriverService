package swagger

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"

	"DriverService/internal/app/appErrors"
	"DriverService/internal/app/services"
	"DriverService/internal/pkg/configs"
	"DriverService/internal/pkg/dto"
)

type Handler struct {
	driverService services.IDriverService
	store         *sessions.CookieStore
}

func NewHanlder(service services.IDriverService, config configs.Config) *Handler {
	return &Handler{driverService: service, store: sessions.NewCookieStore([]byte(config.SecretKey))}
}

func (h *Handler) ApiV1DriverLoginPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var request dto.LogInDriverRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		appErrors.HandleErr(w, appErrors.ErrInvalidData)
		return
	}

	response, err := h.driverService.LogInDriver(request)
	if err != nil {
		appErrors.HandleErr(w, err)
		return
	}

	session, err := h.store.Get(r, "driver-session")
	if err != nil {
		appErrors.HandleErr(w, err)
		return
	}

	session.Values["auth"] = true
	session.Values["phoneNumber"] = response.PhoneNumber
	session.Options = &sessions.Options{MaxAge: 3600}
	if err = session.Save(r, w); err != nil {
		appErrors.HandleErr(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode("Successfully logged in")
}

func (h *Handler) ApiV1DriverLogOut(w http.ResponseWriter, r *http.Request) {
	session, err := h.store.Get(r, "driver-session")
	if err != nil {
		appErrors.HandleErr(w, err)
		return
	}

	if ok := session.Values["auth"]; ok == nil || !ok.(bool) {
		appErrors.HandleErr(w, appErrors.ErrInvalidSession)
		return
	}
	session.Values["auth"] = false
	if err = session.Save(r, w); err != nil {
		appErrors.HandleErr(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode("Successfully logged out")
}

func (h *Handler) ApiV1DriverRegisterPost(w http.ResponseWriter, r *http.Request) {
	var driver dto.Driver
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewDecoder(r.Body).Decode(&driver); err != nil {
		appErrors.HandleErr(w, appErrors.ErrInvalidData)
		return
	}

	if err := h.driverService.RegisterDriver(driver); err != nil {
		appErrors.HandleErr(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode("Successfully added")
}

func (h *Handler) ApiV1DriverStatusPost(w http.ResponseWriter, r *http.Request) {
	session, err := h.store.Get(r, "driver-session")
	if err != nil {
		appErrors.HandleErr(w, err)
		return
	}

	if ok := session.Values["auth"]; ok == nil || !ok.(bool) {
		appErrors.HandleErr(w, appErrors.ErrInvalidSession)
		return
	}

	var status dto.ChangeStatusRequest
	phoneNumber := fmt.Sprint(session.Values["phoneNumber"])
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err = json.NewDecoder(r.Body).Decode(&status); err != nil {
		appErrors.HandleErr(w, appErrors.ErrInvalidData)
		return
	}

	err = h.driverService.SetStatus(status.Status, phoneNumber)
	if err != nil {
		appErrors.HandleErr(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode("Status changed")
}

func (h *Handler) ApiV1DriverGet(w http.ResponseWriter, r *http.Request) {
	session, err := h.store.Get(r, "driver-session")
	if err != nil {
		appErrors.HandleErr(w, err)
		return
	}

	if ok := session.Values["auth"]; ok == nil || !ok.(bool) {
		appErrors.HandleErr(w, appErrors.ErrInvalidSession)
		return
	}

	response, err := h.driverService.GetDriver(fmt.Sprint(session.Values["phoneNumber"]))
	if err != nil {
		appErrors.HandleErr(w, err)
		return
	}

	driverResponse := dto.GetDriverResponse{
		Name:        response.Name,
		PhoneNumber: response.PhoneNumber,
		TaxiType:    response.TaxiType,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(driverResponse)
}
