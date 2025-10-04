package v1

import (
	"kompass/internal/usecase"
	"kompass/pkg/logger"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func NewUserRoutes(apiV1Group fiber.Router, uc usecase.Users, log logger.Interface) {
	r := &UsersV1{uc: uc, log: log, v: validator.New(validator.WithRequiredStructEnabled())}

	usersV1Group := apiV1Group.Group("/users")

	{
		usersV1Group.Get("", r.getUsers)
		usersV1Group.Get("/:user_id", r.getUser)
	}
}

func NewGeocodingRoutes(apiV1Group fiber.Router, uc usecase.Geocoding, log logger.Interface) {
	r := &GeocodingV1{uc: uc, log: log, v: validator.New(validator.WithRequiredStructEnabled())}

	serviceV1Group := apiV1Group.Group("/geocoding")

	{
		serviceV1Group.Get("/location", r.getLocation)
		serviceV1Group.Get("/station", r.getTrainStation)
	}
}

func NewTripRoutes(apiV1Group fiber.Router, uc usecase.Trips, log logger.Interface, authorization func(c *fiber.Ctx) error) fiber.Router {
	r := &TripsV1{
		uc:  uc,
		log: log,
		v:   validator.New(validator.WithRequiredStructEnabled()),
	}

	apiV1Group.Get("/trips", r.getTrips)
	apiV1Group.Post("/trips", r.postTrip)

	tripsV1Group := apiV1Group.Group("/trips/:trip_id")
	tripsV1Group.Use(authorization)
	{
		tripsV1Group.Get("", r.getTrip)
		tripsV1Group.Put("", r.putTrip)
		tripsV1Group.Delete("", r.deleteTrip)
	}

	return tripsV1Group
}

func NewTransportationRoutes(tripsV1Group fiber.Router, uc usecase.Transportation, log logger.Interface) {
	r := &TransportationV1{
		uc:  uc,
		log: log,
		v:   validator.New(validator.WithRequiredStructEnabled()),
	}

	transportationV1Group := tripsV1Group.Group("/transportation")

	{
		transportationV1Group.Post("", r.postTransportation)
		transportationV1Group.Get("", r.getAllTransportation)
		transportationV1Group.Get("/geojson", r.getGeoJson)
		transportationV1Group.Get("/:transportation_id", r.getTransportation)
		transportationV1Group.Delete("/:transportation_id", r.deleteTransportation)
	}
}

func NewFlightRoutes(tripsV1Group fiber.Router, uc usecase.Flights, log logger.Interface) {
	r := &FlightsV1{uc: uc, log: log, v: validator.New(validator.WithRequiredStructEnabled())}

	flightsV1Group := tripsV1Group.Group("/flights")

	{
		flightsV1Group.Post("", r.postFlight)
		flightsV1Group.Put("/:flight_id", r.putFlight)
	}
}

func NewTrainRoutes(apiV1Group fiber.Router, uc usecase.Trains, log logger.Interface) {
	r := &TrainsV1{uc: uc, log: log, v: validator.New(validator.WithRequiredStructEnabled())}

	trainsV1Group := apiV1Group.Group("/trains")

	{
		trainsV1Group.Post("", r.postTrainJourney)
	}
}

func NewActivityRoutes(tripsV1Group fiber.Router, uc usecase.Activities, log logger.Interface) {
	r := &ActivitiesV1{uc: uc, log: log, v: validator.New(validator.WithRequiredStructEnabled())}

	activitiesV1Group := tripsV1Group.Group("/activities")

	{
		activitiesV1Group.Get("", r.getActivities)
		activitiesV1Group.Post("", r.postActivity)
		activitiesV1Group.Get("/:activity_id", r.getActivity)
		activitiesV1Group.Put("/:activity_id", r.putActivity)
		activitiesV1Group.Delete("/:activity_id", r.deleteActivity)
	}
}

func NewAccommodationRoutes(tripsV1Group fiber.Router, uc usecase.Accommodation, log logger.Interface) {
	r := &AccommodationV1{uc: uc, log: log, v: validator.New(validator.WithRequiredStructEnabled())}

	accommodationV1Group := tripsV1Group.Group("/accommodation")

	{
		accommodationV1Group.Get("", r.getAllAccommodation)
		accommodationV1Group.Post("", r.postAccommodation)
		accommodationV1Group.Get("/:accommodation_id", r.getAccommodationByID)
		accommodationV1Group.Put("/:accommodation_id", r.putAccommodation)
		accommodationV1Group.Delete("/:accommodation_id", r.deleteAccommodation)
	}
}

func NewAttachmentRoutes(tripsV1Group fiber.Router, uc usecase.Attachments, log logger.Interface) {
	r := &AttachmentsV1{uc: uc, log: log, v: validator.New(validator.WithRequiredStructEnabled())}

	attachmentsV1Group := tripsV1Group.Group("/attachments")

	{
		attachmentsV1Group.Get("", r.getAttachments)
		attachmentsV1Group.Post("", r.postAttachment)
		attachmentsV1Group.Get("/:attachment_id", r.downloadAttachment)
		attachmentsV1Group.Delete("/:attachment_id", r.deleteAttachment)
	}
}
