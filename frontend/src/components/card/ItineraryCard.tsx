import FlightEntry from "@/components/itinerary/FlightEntry.tsx";
import DaySeparator from "@/components/itinerary/DaySeparator.tsx";
import {Accommodation, Flight, Activity, Trip, FlightLeg} from "@/types.ts";
import {formatDuration, getDaysBetween, isSameDay} from "@/components/util.ts";
import DayLabel from "@/components/itinerary/DayLabel.tsx";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card.tsx";
import React from "react";
import ActivityEntry from "@/components/itinerary/ActivityEntry.tsx";
import FlightSeparator from "@/components/itinerary/FlightSeparator.tsx";
import AddFlightButton from "@/components/dialog/AddFlightButton.tsx";

export default async function ItineraryCard({
  trip, flights, accomodation, activities
}: {
  trip: Trip,
  accomodation: Accommodation[],
  flights: Flight[],
  activities: Activity[]
}) {

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

      const accommodation = accomodation.find(acc =>
          acc.arrivalDate <= day && acc.departureDate > day
      )

      if (isSameDay(day, trip.endDate) || grouped.length == 0 || filteredFlights.length != 0 ||
          filteredActivities.length != 0 || accommodation != grouped[grouped.length-1].accommodation) {
        grouped.push({
          day: day,
          flights: filteredFlights,
          activities: filteredActivities,
          accommodation: accommodation,
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
        {dayData.activities.map(act => <ActivityEntry key={act.id} activity={act}/>)}
        {dayFlights.map((pair, idx) =>
          <div key={idx}>
            <FlightEntry flight={pair.flight} flightLeg={pair.leg}/>
            {dayData.flights.length > idx + 1 &&
              <span className="ml-6 text-sm text-muted-foreground">
                {formatDuration(pair.leg.arrivalDateTime, dayData.flights[idx+1].leg.departureDateTime)} Layover
              </span>
            }
          </div>
        )}
        {nextDay && (nightFlight ?
          <FlightSeparator date={nextDay} flight={nightFlight.flight} flightLeg={nightFlight.leg}/>
        :
          <DaySeparator date={nextDay} accomodation={dayData.accommodation?.name} collapsedDays={collapsedDays}/>
        )}
      </div>
    )
  }

  return (
    <div className="w-1/2 h-full">
      <Card className="overflow-hidden card h-full">
        <CardHeader className="py-5 border-b flex-row justify-between space-y-0">
          <CardTitle className="h-8 text-[1.5rem]">Itinerary</CardTitle>
          <AddFlightButton trip={trip}/>
        </CardHeader>
        <CardContent className="p-0 pb-4 h-full overflow-y-scroll no-scrollbar" style={{height: "calc(100% - 4rem)"}}>
          <DayLabel date={trip.startDate}/>
          {renderAllDays()}
        </CardContent>
      </Card>
    </div>
  )
}