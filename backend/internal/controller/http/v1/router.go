package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"travel-planner/internal/usecase"
	"travel-planner/pkg/logger"
)

func NewUserRoutes(apiV1Group fiber.Router, uc usecase.Users, log logger.Interface) {
	r := &UsersV1{uc: uc, log: log, v: validator.New(validator.WithRequiredStructEnabled())}

	usersV1Group := apiV1Group.Group("/users")

	{
		usersV1Group.Get("", r.getUsers)
		usersV1Group.Get("/:user_id", r.getUser)
	}
}

func NewTripRoutes(apiV1Group fiber.Router, uc usecase.Trips, log logger.Interface) fiber.Router {
	r := &TripsV1{uc: uc, log: log, v: validator.New(validator.WithRequiredStructEnabled())}

	tripsV1Group := apiV1Group.Group("/trips")

	{
		tripsV1Group.Get("", r.getTrips)
		tripsV1Group.Get("/:trip_id", r.getTrip)
	}

	return tripsV1Group
}

func NewFlightRoutes(tripsV1Group fiber.Router, uc usecase.Flights, log logger.Interface) {
	r := &FlightsV1{uc: uc, log: log, v: validator.New(validator.WithRequiredStructEnabled())}

	flightsV1Group := tripsV1Group.Group("/:trip_id/flights")

	{
		flightsV1Group.Get("", r.getFlights)
		flightsV1Group.Post("", r.postFlight)
		flightsV1Group.Get("/:flight_id", r.getFlight)
	}
}

func NewActivityRoutes(tripsV1Group fiber.Router, uc usecase.Activities, log logger.Interface) {
	r := &ActivitiesV1{uc: uc, log: log, v: validator.New(validator.WithRequiredStructEnabled())}

	activitiesV1Group := tripsV1Group.Group("/:trip_id/activities")

	{
		activitiesV1Group.Get("", r.getActivities)
		activitiesV1Group.Post("", r.postActivity)
		activitiesV1Group.Get("/:activity_id", r.getActivity)
	}
}

func NewAccommodationRoutes(tripsV1Group fiber.Router, uc usecase.Accommodation, log logger.Interface) {
	r := &AccommodationV1{uc: uc, log: log, v: validator.New(validator.WithRequiredStructEnabled())}

	accommodationV1Group := tripsV1Group.Group("/:trip_id/accommodation")

	{
		accommodationV1Group.Get("", r.getAllAccommodation)
		accommodationV1Group.Post("", r.postAccommodation)
		accommodationV1Group.Get("/:accommodation_id", r.getAccommodationByID)
	}
}
