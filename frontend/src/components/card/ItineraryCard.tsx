import FlightEntry from "@/components/itinerary/FlightEntry.tsx";
import DaySeparator from "@/components/itinerary/DaySeparator.tsx";
import {Accommodation, Flight, Activity, Trip} from "@/types.ts";
import moment, {Moment} from "moment";
import FlightSeparator from "@/components/itinerary/FlightSeparator.tsx";
import {durationString, getDaysBetween, isSameLocalDay} from "@/components/util.ts";
import DayLabel from "@/components/itinerary/DayLabel.tsx";
import AddActivityButton from "@/components/itinerary/AddActivityButton.tsx";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card.tsx";
import React from "react";
import ActivityEntry from "@/components/itinerary/ActivityEntry.tsx";

export default async function ItineraryCard({
  trip, flights, accomodation, activities
}: {
  trip: Trip,
  accomodation: Accommodation[],
  flights: Flight[],
  activities: Activity[]
}) {

  function renderDay(day: Moment) {
    const filteredActivities = activities
        .filter(act => day.isSame(act.date))

    const filteredFlights = flights
      .flatMap(flight => flight.legs.map(leg => ({flight, leg})))
      .filter(pair => isSameLocalDay(pair.leg.departureTime, day))

    const nightFlight = filteredFlights
      .find(pair => !isSameLocalDay(pair.leg.arrivalTime, day))

    const dayFlights = nightFlight ? filteredFlights.slice(0, -1) : filteredFlights

    const todaysAccomodation = accomodation.find(acc =>
      acc.arrivalDate.isSameOrBefore(day) && acc.departureDate.isAfter(day)
    )

    const nextDay = moment(day).add(1, "day")

    return (
      <div key={day.toISOString()}>
        {filteredActivities.map(act => <ActivityEntry key={act.id} name={act.name}/>)}
        <AddActivityButton/>
        {dayFlights.map((pair, idx) =>
          <div key={idx}>
            <FlightEntry flight={pair.flight} flightLeg={pair.leg} data-superjson/>
            {filteredFlights.length > idx + 1 &&
              <span className="ml-6 text-sm text-muted-foreground">
                {durationString(filteredFlights[idx+1].leg.departureTime, pair.leg.arrivalTime)} Layover
              </span>
            }
          </div>
        )}
        {!trip.endDate.isSame(day) && (nightFlight ?
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