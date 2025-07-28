package v1

import (
	"backplate/internal/usecase"
	"backplate/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func NewUserRoutes(apiV1Group fiber.Router, uc usecase.Users, log logger.Interface) {
	r := &UsersV1{uc: uc, log: log, v: validator.New(validator.WithRequiredStructEnabled())}

	userGroup := apiV1Group.Group("/users")

	{
		userGroup.Get("", r.getUsers)
		userGroup.Get("/:user_id", r.getUser)
	}
}

func NewTripRoutes(apiV1Group fiber.Router, uc usecase.Trips, log logger.Interface) {
	r := &TripsV1{uc: uc, log: log, v: validator.New(validator.WithRequiredStructEnabled())}

	deviceGroup := apiV1Group.Group("/trips")

	{
		deviceGroup.Get("", r.getTrips)
		deviceGroup.Get("/:trip_id", r.getTrip)
	}
}

func NewFlightRoutes(apiV1Group fiber.Router, uc usecase.Flights, log logger.Interface) {
	r := &FlightsV1{uc: uc, log: log, v: validator.New(validator.WithRequiredStructEnabled())}

	imageGroup := apiV1Group.Group("/flights")

	{
		imageGroup.Get("", r.getFlights)
		imageGroup.Get("/:flight_id", r.getFlight)
	}
}

func NewActivityRoutes(apiV1Group fiber.Router, uc usecase.Activities, log logger.Interface) {
	r := &ActivitiesV1{uc: uc, log: log, v: validator.New(validator.WithRequiredStructEnabled())}

	imageGroup := apiV1Group.Group("/activities")

	{
		imageGroup.Get("", r.getActivities)
		imageGroup.Get("/:activity_id", r.getActivity)
	}
}

func NewAccommodationRoutes(apiV1Group fiber.Router, uc usecase.Accommodation, log logger.Interface) {
	r := &AccommodationV1{uc: uc, log: log, v: validator.New(validator.WithRequiredStructEnabled())}

	imageGroup := apiV1Group.Group("/accommodation")

	{
		imageGroup.Get("", r.getAllAccommodation)
		imageGroup.Get("/:accommodation_id", r.getAccommodationByID)
	}
}
