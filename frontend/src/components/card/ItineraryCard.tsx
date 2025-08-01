import FlightEntry from "@/components/itinerary/FlightEntry.tsx";
import DaySeparator from "@/components/itinerary/DaySeparator.tsx";
import {Accommodation, Flight, Activity, Trip} from "@/types.ts";
import {formatDuration, getDaysBetween, isSameDay} from "@/components/util.ts";
import DayLabel from "@/components/itinerary/DayLabel.tsx";
import AddActivityButton from "@/components/itinerary/AddActivityButton.tsx";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card.tsx";
import React from "react";
import ActivityEntry from "@/components/itinerary/ActivityEntry.tsx";
import FlightSeparator from "@/components/itinerary/FlightSeparator.tsx";

export default async function ItineraryCard({
  trip, flights, accomodation, activities
}: {
  trip: Trip,
  accomodation: Accommodation[],
  flights: Flight[],
  activities: Activity[]
}) {

  function renderDay(day: Date) {
    const filteredActivities = activities
        .filter(act => isSameDay(day, act.date))

    const filteredFlights = flights
      .flatMap(flight => flight.legs.map(leg => ({flight, leg})))
      .filter(pair => isSameDay(pair.leg.departureDateTime, day))

    const nightFlight = filteredFlights
      .find(pair => !isSameDay(pair.leg.arrivalDateTime, day))

    const dayFlights = nightFlight ? filteredFlights.slice(0, -1) : filteredFlights

    const todaysAccomodation = accomodation.find(acc =>
      acc.arrivalDate <= day && acc.departureDate > day
    )

    const nextDay = new Date(day)
    nextDay.setDate(nextDay.getDate() + 1)

    return (
      <div key={day.toISOString()}>
        {filteredActivities.map(act => <ActivityEntry key={act.id} activity={act}/>)}
        <AddActivityButton/>
        {dayFlights.map((pair, idx) =>
          <div key={idx}>
            <FlightEntry flight={pair.flight} flightLeg={pair.leg}/>
            {filteredFlights.length > idx + 1 &&
              <span className="ml-6 text-sm text-muted-foreground">
                {formatDuration(pair.leg.arrivalDateTime, filteredFlights[idx+1].leg.departureDateTime)} Layover
              </span>
            }
          </div>
        )}
        {!isSameDay(trip.endDate, day) && (nightFlight ?
          <FlightSeparator date={nextDay} flight={nightFlight.flight} flightLeg={nightFlight.leg}/>
        :
          <DaySeparator date={nextDay} accomodation={todaysAccomodation?.name}/>
        )}
      </div>
    )
  }

  return (
    <div className="w-1/2 h-full">
      <Card className="overflow-hidden card h-full">
        <CardHeader className="py-5 border-b">
          <CardTitle>Itinerary</CardTitle>
        </CardHeader>
        <CardContent className="p-0 pb-4 h-full overflow-y-scroll no-scrollbar" style={{height: "calc(100% - 4rem)"}}>
          <DayLabel date={trip.startDate}/>
          {getDaysBetween(trip.startDate, trip.endDate).map(renderDay)}
        </CardContent>
      </Card>
    </div>
  )
}