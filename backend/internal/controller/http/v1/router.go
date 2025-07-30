package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"travel-planner/internal/usecase"
	"travel-planner/pkg/logger"
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

	tripGroup := apiV1Group.Group("/trips")

	{
		tripGroup.Get("", r.getTrips)
		tripGroup.Get("/:trip_id", r.getTrip)
	}
}

func NewFlightRoutes(apiV1Group fiber.Router, uc usecase.Flights, log logger.Interface) {
	r := &FlightsV1{uc: uc, log: log, v: validator.New(validator.WithRequiredStructEnabled())}

	flightGroup := apiV1Group.Group("/flights")

	{
		flightGroup.Get("", r.getFlights)
		flightGroup.Post("", r.postFlight)
		flightGroup.Get("/:flight_id", r.getFlight)
	}
}

func NewActivityRoutes(apiV1Group fiber.Router, uc usecase.Activities, log logger.Interface) {
	r := &ActivitiesV1{uc: uc, log: log, v: validator.New(validator.WithRequiredStructEnabled())}

	activityGroup := apiV1Group.Group("/activities")

	{
		activityGroup.Get("", r.getActivities)
		activityGroup.Get("/:activity_id", r.getActivity)
	}
}

func NewAccommodationRoutes(apiV1Group fiber.Router, uc usecase.Accommodation, log logger.Interface) {
	r := &AccommodationV1{uc: uc, log: log, v: validator.New(validator.WithRequiredStructEnabled())}

	accommodationGroup := apiV1Group.Group("/accommodation")

	{
		accommodationGroup.Get("", r.getAllAccommodation)
		accommodationGroup.Get("/:accommodation_id", r.getAccommodationByID)
	}
}
