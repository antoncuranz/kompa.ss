import {DayRenderData} from "@/types.ts";
import {dayIsBetween, getDaysBetween, isSameDay} from "@/components/util.ts";
import React from "react";
import AddSomethingDropdown from "@/components/dialog/AddSomethingDropdown.tsx";
import Card from "@/components/card/Card.tsx";
import {fetchAccommodation, fetchActivities, fetchTransportation, fetchTrip} from "@/requests.ts";
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
  const transportation = await fetchTransportation(tripId)

  function processDataAndGroupByDays() {
    const grouped: DayRenderData[] = []

    for (const day of getDaysBetween(trip.startDate, trip.endDate)) {
      const filteredActivities = activities
          .filter(act => isSameDay(day, act.date))

      const relevantTransportation = transportation
          .filter(t => dayIsBetween(day, t.departureDateTime, t.arrivalDateTime))

      const relevantAccommodation = accommodation.find(acc =>
          acc.arrivalDate <= day && acc.departureDate > day
      )

      // TODO: also push if day is today!
      if (isSameDay(day, trip.endDate) || grouped.length == 0 || relevantTransportation.length != 0 ||
          filteredActivities.length != 0 || relevantAccommodation != grouped[grouped.length-1].accommodation ||
          grouped[grouped.length-1].transportation.find(pair => isSameDay(pair.arrivalDateTime, day))
      ) {
        grouped.push({
          day: day,
          transportation: relevantTransportation,
          activities: filteredActivities,
          accommodation: relevantAccommodation,
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