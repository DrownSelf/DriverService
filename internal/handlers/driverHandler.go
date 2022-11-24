package handlers

import (
	"encoding/json"
	"net/http"

	pb "github.com/DrownSelf/OrderService/pkg/grpc"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/DrownSelf/DriverService/internal/appErrors"
	"github.com/DrownSelf/DriverService/internal/auth"
	"github.com/DrownSelf/DriverService/internal/configs"
	"github.com/DrownSelf/DriverService/internal/entities"
	"github.com/DrownSelf/DriverService/internal/services"
	"github.com/DrownSelf/DriverService/models"
	"github.com/DrownSelf/DriverService/restapi/operations"
	"github.com/DrownSelf/DriverService/restapi/operations/actions_with_token"
	"github.com/DrownSelf/DriverService/restapi/operations/logging_actions"
)

type IDriverHandler interface {
	GetDriver(params actions_with_token.GetDriverParams) middleware.Responder
	Logout(params actions_with_token.GetDriverLogoutParams) middleware.Responder
	LogIn(params logging_actions.PostDriverLoginParams) middleware.Responder
	Register(params logging_actions.PostDriverRegisterParams) middleware.Responder
	EndRide(params actions_with_token.PostDriverEndRideParams) middleware.Responder
	TakeOrder(params actions_with_token.PostDriverTakeOrderParams) middleware.Responder
	RateUser(params actions_with_token.PostDriverRateUserParams) middleware.Responder
	RateDriverFromOrder(params actions_with_token.PostOrderRateDriverParams) middleware.Responder
}

type CustomResponder func(http.ResponseWriter, runtime.Producer)

func (c CustomResponder) WriteResponse(w http.ResponseWriter, p runtime.Producer) {
	c(w, p)
}

type HandlerDependencies struct {
	DriverService services.IDriverService
	OrderClient   pb.OrderServiceClient
	Forger        auth.TokenForger
	Config        *configs.Config
}

type DriverHandler struct {
	driverService services.IDriverService
	orderClient   pb.OrderServiceClient
	forger        auth.TokenForger
	config        *configs.Config
}

func NewDriverHandler(dependencies HandlerDependencies) (*DriverHandler, error) {
	return &DriverHandler{
		driverService: dependencies.DriverService,
		orderClient:   dependencies.OrderClient,
		forger:        dependencies.Forger,
		config:        dependencies.Config,
	}, nil
}

func ConfigureHandlers(swaggerApi *operations.DriverServiceAPI, handler IDriverHandler) {
	swaggerApi.ActionsWithTokenGetDriverHandler = actions_with_token.GetDriverHandlerFunc(handler.GetDriver)
	swaggerApi.LoggingActionsPostDriverLoginHandler = logging_actions.PostDriverLoginHandlerFunc(handler.LogIn)
	swaggerApi.LoggingActionsPostDriverRegisterHandler = logging_actions.PostDriverRegisterHandlerFunc(handler.Register)
	swaggerApi.ActionsWithTokenPostDriverRateUserHandler = actions_with_token.PostDriverRateUserHandlerFunc(handler.RateUser)
	swaggerApi.ActionsWithTokenPostOrderRateDriverHandler = actions_with_token.PostOrderRateDriverHandlerFunc(handler.RateDriverFromOrder)
	swaggerApi.ActionsWithTokenGetDriverLogoutHandler = actions_with_token.GetDriverLogoutHandlerFunc(handler.Logout)
	swaggerApi.ActionsWithTokenPostDriverEndRideHandler = actions_with_token.PostDriverEndRideHandlerFunc(handler.EndRide)
	swaggerApi.ActionsWithTokenPostDriverTakeOrderHandler = actions_with_token.PostDriverTakeOrderHandlerFunc(handler.TakeOrder)
}

func (h *DriverHandler) GetDriver(params actions_with_token.GetDriverParams) middleware.Responder {
	return CustomResponder(func(w http.ResponseWriter, p runtime.Producer) {
		cookie, err := params.HTTPRequest.Cookie("Token")
		if err != nil {
			appErrors.HandleErr(w, appErrors.ErrInvalidSession)
			return
		}

		claims, err := h.forger.Decode(cookie.Value)
		if err != nil {
			appErrors.HandleErr(w, appErrors.ErrInvalidSession)
			return
		}

		driver, err := h.driverService.GetDriver(claims.PhoneNumber)
		if err != nil {
			appErrors.HandleErr(w, err)
			return
		}

		response := models.GetDriverResponse{
			Name:        driver.Name,
			PhoneNumber: driver.PhoneNumber,
			TaxiType:    driver.TaxiType,
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	})
}

func (h *DriverHandler) LogIn(params logging_actions.PostDriverLoginParams) middleware.Responder {
	return CustomResponder(func(w http.ResponseWriter, p runtime.Producer) {
		driver := entities.Driver{
			PhoneNumber: params.LogInRequest.PhoneNumber,
			Password:    params.LogInRequest.Password,
		}

		gottenDriver, err := h.driverService.LogInDriver(driver)
		if err != nil {
			appErrors.HandleErr(w, err)
			return
		}
		token, err := h.forger.Encode(
			auth.TokenClaims{
				Name:        gottenDriver.Name,
				Email:       gottenDriver.Email,
				PhoneNumber: gottenDriver.PhoneNumber,
			},
			*h.config)

		cookie := http.Cookie{
			Name:   "Token",
			Value:  token,
			MaxAge: 3600,
		}

		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode("Successfully logged in")
	})
}

func (h *DriverHandler) Logout(params actions_with_token.GetDriverLogoutParams) middleware.Responder {
	return CustomResponder(func(w http.ResponseWriter, producer runtime.Producer) {
		cookie, err := params.HTTPRequest.Cookie("Token")
		if err != nil {
			appErrors.HandleErr(w, appErrors.ErrInvalidSession)
			return
		}

		if _, err = h.forger.Decode(cookie.Value); err != nil {
			appErrors.HandleErr(w, appErrors.ErrInvalidSession)
			return
		}

		cookie.MaxAge = -1
		http.SetCookie(w, cookie)

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode("Successfully logged out")
	})
}

func (h *DriverHandler) Register(params logging_actions.PostDriverRegisterParams) middleware.Responder {
	return CustomResponder(func(w http.ResponseWriter, producer runtime.Producer) {
		driver := entities.Driver{
			Name:        params.RegisterRequest.Name,
			PhoneNumber: params.RegisterRequest.PhoneNumber,
			Password:    params.RegisterRequest.Password,
			Email:       params.RegisterRequest.Email,
			TaxiType:    params.RegisterRequest.TaxiType,
		}

		if err := h.driverService.RegisterDriver(driver); err != nil {
			appErrors.HandleErr(w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode("Successfully added")
	})
}

func (h *DriverHandler) EndRide(params actions_with_token.PostDriverEndRideParams) middleware.Responder {
	return CustomResponder(func(w http.ResponseWriter, producer runtime.Producer) {
		cookie, err := params.HTTPRequest.Cookie("Token")
		if err != nil {
			appErrors.HandleErr(w, appErrors.ErrInvalidSession)
			return
		}

		if _, err = h.forger.Decode(cookie.Value); err != nil {
			appErrors.HandleErr(w, appErrors.ErrInvalidSession)
			return
		}

		_, err = h.orderClient.EndRide(params.HTTPRequest.Context(), &pb.EndRideRequest{
			OrderId:     params.ChangeStatusOfRequest.ID,
			OrderStatus: false,
		})

		if err != nil {
			appErrors.HandleErr(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode("Ride Ended")
	})
}

func (h *DriverHandler) TakeOrder(params actions_with_token.PostDriverTakeOrderParams) middleware.Responder {
	return CustomResponder(func(w http.ResponseWriter, producer runtime.Producer) {
		cookie, err := params.HTTPRequest.Cookie("Token")
		if err != nil {
			appErrors.HandleErr(w, appErrors.ErrInvalidSession)
			return
		}

		if _, err = h.forger.Decode(cookie.Value); err != nil {
			appErrors.HandleErr(w, appErrors.ErrInvalidSession)
			return
		}

		order, err := h.orderClient.TakeOrder(params.HTTPRequest.Context(), &pb.ServeClientRequest{
			Driver: &pb.Driver{
				Name:        params.RequestToAssignInWatingQueue.Name,
				Email:       params.RequestToAssignInWatingQueue.Email,
				PhoneNumber: params.RequestToAssignInWatingQueue.PhoneNumber,
				TaxiType:    params.RequestToAssignInWatingQueue.TaxiType,
			},
		})
		if err != nil {
			appErrors.HandleErr(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(models.DriverRideResponse{
			From:            order.From,
			To:              order.To,
			OrderID:         order.OrderId,
			UserName:        order.User.Name,
			UserPhoneNumber: order.User.PhoneNumber,
		})
	})
}

func (h *DriverHandler) RateUser(params actions_with_token.PostDriverRateUserParams) middleware.Responder {
	return CustomResponder(func(w http.ResponseWriter, producer runtime.Producer) {
		cookie, err := params.HTTPRequest.Cookie("Token")
		if err != nil {
			appErrors.HandleErr(w, appErrors.ErrInvalidSession)
			return
		}

		if _, err = h.forger.Decode(cookie.Value); err != nil {
			appErrors.HandleErr(w, appErrors.ErrInvalidSession)
			return
		}

		_, err = h.orderClient.RateRideFromDriver(params.HTTPRequest.Context(), &pb.RateUserFromDriver{
			OrderId: params.RequestToRateUser.ID,
			Rating:  int32(params.RequestToRateUser.Rating),
		})

		if err != nil {
			appErrors.HandleErr(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode("Rating added successfully")
	})
}

func (h *DriverHandler) RateDriverFromOrder(params actions_with_token.PostOrderRateDriverParams) middleware.Responder {
	return CustomResponder(func(w http.ResponseWriter, producer runtime.Producer) {
		if err := h.driverService.RateDriverFromOrder(params.RequestToAddRatingToDriver.PhoneNumber, float64(params.RequestToAddRatingToDriver.Rating)); err != nil {
			appErrors.HandleErr(w, appErrors.ErrInvalidData)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode("rating added")
	})
}
