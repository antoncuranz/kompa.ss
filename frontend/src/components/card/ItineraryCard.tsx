import {DayRenderData} from "@/types.ts";
import {getDaysBetween, isSameDay} from "@/components/util.ts";
import React from "react";
import AddSomethingDropdown from "@/components/dialog/AddSomethingDropdown.tsx";
import Card from "@/components/card/Card.tsx";
import {fetchAccommodation, fetchActivities, fetchFlights, fetchTrip} from "@/requests.ts";
import Itinerary from "@/components/itinerary/Itinerary.tsx";

export default async function ItineraryCard({
  tripId, className
}: {
  tripId: number
  className?: string
}) {
  const trip = await fetchTrip(tripId)
  const activities = await fetchActivities(tripId)
  const accommodation = await fetchAccommodation(tripId)
  const flights = await fetchFlights(tripId)

  function processDataAndGroupByDays() {
    const grouped: DayRenderData[] = []

    for (const day of getDaysBetween(trip.startDate, trip.endDate)) {
      const filteredActivities = activities
          .filter(act => isSameDay(day, act.date))

      const filteredFlights = flights
          .flatMap(flight => flight.legs.map(leg => ({flight, leg})))
          .filter(pair => isSameDay(pair.leg.departureDateTime, day))

      const filteredAccommodation = accommodation.find(acc =>
          acc.arrivalDate <= day && acc.departureDate > day
      )

      // TODO: also push if day is today!
      if (isSameDay(day, trip.endDate) || grouped.length == 0 || filteredFlights.length != 0 ||
          filteredActivities.length != 0 || filteredAccommodation != grouped[grouped.length-1].accommodation ||
          grouped[grouped.length-1].flights.find(pair => isSameDay(pair.leg.arrivalDateTime, day))
      ) {
        grouped.push({
          day: day,
          flights: filteredFlights,
          activities: filteredActivities,
          accommodation: filteredAccommodation,
        })
      }
    }

    return grouped
  }

  return (
    <Card title="Itinerary" headerSlot={<AddSomethingDropdown trip={trip}/>} className={className}>
      <Itinerary trip={trip} dataByDays={processDataAndGroupByDays()}/>
    </Card>
  )
}