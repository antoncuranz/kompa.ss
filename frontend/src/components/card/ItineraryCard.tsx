import FlightEntry from "@/components/itinerary/FlightEntry.tsx";
import DaySeparator from "@/components/itinerary/DaySeparator.tsx";
import {Accommodation, Flight, Activity, FlightLeg} from "@/types.ts";
import {formatDuration, getDaysBetween, isSameDay} from "@/components/util.ts";
import DayLabel from "@/components/itinerary/DayLabel.tsx";
import React from "react";
import ActivityEntry from "@/components/itinerary/ActivityEntry.tsx";
import FlightSeparator from "@/components/itinerary/FlightSeparator.tsx";
import AddSomethingDropdown from "@/components/dialog/AddSomethingDropdown.tsx";
import Card from "@/components/card/Card.tsx";
import {fetchAccommodation, fetchActivities, fetchFlights, fetchTrip} from "@/requests.ts";

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

  type DayRenderData = {
    day: Date;
    flights: {flight: Flight; leg: FlightLeg}[];
    activities: Activity[];
    accommodation: Accommodation | undefined;
  };

  function renderAllDays() {
    const dataByDays = processDataAndGroupByDays()
    return dataByDays.map((dayData, idx) => renderDay(dayData, dataByDays[idx+1]?.day))
  }

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

  function renderDay(dayData: DayRenderData, nextDay: Date) {
    const collapsedDays = nextDay ? getDaysBetween(dayData.day, nextDay).length-2 : 0

    const nightFlight = dayData.flights
      .find(pair => !isSameDay(pair.leg.arrivalDateTime, dayData.day))

    const dayFlights = nightFlight ? dayData.flights.slice(0, -1) : dayData.flights

    return (
      <div key={dayData.day.toISOString()}>
        <DayLabel date={dayData.day}/>

        {dayData.activities.map(act => <ActivityEntry key={act.id} activity={act}/>)}
        {dayFlights.map((pair, idx) =>
          <div key={idx}>
            <FlightEntry flight={pair.flight} flightLeg={pair.leg}/>
            {dayData.flights.length > idx + 1 &&
              <span className="mx-3 text-sm text-muted-foreground">
                {formatDuration(pair.leg.arrivalDateTime, dayData.flights[idx+1].leg.departureDateTime)} Layover
              </span>
            }
          </div>
        )}
        {nextDay && (nightFlight ?
          <FlightSeparator flight={nightFlight.flight} flightLeg={nightFlight.leg}/>
        :
          <DaySeparator accomodation={dayData.accommodation?.name} collapsedDays={collapsedDays}/>
        )}
      </div>
    )
  }

  return (
    <Card title="Itinerary" headerSlot={<AddSomethingDropdown trip={trip}/>} className={className}>
      {renderAllDays()}
    </Card>
  )
}