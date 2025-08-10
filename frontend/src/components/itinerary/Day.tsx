import FlightEntry from "@/components/itinerary/FlightEntry.tsx";
import DaySeparator from "@/components/itinerary/DaySeparator.tsx";
import {formatDuration, getDaysBetween, isSameDay} from "@/components/util.ts";
import DayLabel from "@/components/itinerary/DayLabel.tsx";
import React from "react";
import ActivityEntry from "@/components/itinerary/ActivityEntry.tsx";
import {Accommodation, Activity, DayRenderData, Flight} from "@/types.ts";
import {Separator} from "@/components/ui/separator.tsx";

export default function Day({
  dayData,
  nextDay,
  onActivityClick = () => {},
  onAccommodationClick = () => {},
  onFlightClick = () => {}
}: {
  dayData: DayRenderData
  nextDay: string
  onActivityClick?: (activity: Activity) => void;
  onAccommodationClick?: (accommodation: Accommodation | undefined) => void;
  onFlightClick?: (flight: Flight) => void;
}) {

  const collapsedDays = nextDay ? getDaysBetween(dayData.day, nextDay).length-2 : 0

  const nightFlight = dayData.flights
    .find(pair => !isSameDay(pair.leg.arrivalDateTime, dayData.day))

  const dayFlights = nightFlight ? dayData.flights.slice(0, -1) : dayData.flights

  return (
    <div>
      <DayLabel date={dayData.day}/>

      {dayData.activities.map(act =>
          <ActivityEntry key={act.id} activity={act} onClick={() => onActivityClick(act)}/>
      )}

      {dayFlights.map((pair, idx) =>
        <div key={idx}>
          <FlightEntry flight={pair.flight} flightLeg={pair.leg} onInfoBtnClick={() => onFlightClick(pair.flight)}/>
          {dayData.flights.length > idx + 1 &&
            <span className="mx-3 text-sm text-muted-foreground">
              {formatDuration(pair.leg.arrivalDateTime, dayData.flights[idx+1].leg.departureDateTime)} Layover
            </span>
          }
        </div>
      )}

      {nextDay && (nightFlight ?
        <>
          <FlightEntry flight={nightFlight.flight} flightLeg={nightFlight.leg} onInfoBtnClick={() => onFlightClick(nightFlight?.flight)}/>
          <Separator className="relative bottom-5 z-0"/>
        </>
      :
        <DaySeparator accomodation={dayData.accommodation} collapsedDays={collapsedDays} onAccommodationClick={onAccommodationClick}/>
      )}
    </div>
  )
}