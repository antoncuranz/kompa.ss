import FlightEntry from "@/components/itinerary/FlightEntry.tsx";
import DaySeparator from "@/components/itinerary/DaySeparator.tsx";
import {formatDuration, getDaysBetween, isSameDay} from "@/components/util.ts";
import DayLabel from "@/components/itinerary/DayLabel.tsx";
import React from "react";
import ActivityEntry from "@/components/itinerary/ActivityEntry.tsx";
import {Accommodation, Activity, DayRenderData, Transportation, TransportationType} from "@/types.ts";
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
  onFlightClick?: (flight: Transportation) => void;
}) {

  const collapsedDays = nextDay ? getDaysBetween(dayData.day, nextDay).length-2 : 0

  const nightTransportation = dayData.transportation
    .find(pair => !isSameDay(pair.arrivalDateTime, dayData.day))

  const dayTransportation = nightTransportation ? dayData.transportation.slice(0, -1) : dayData.transportation

  function renderDayTransportation(transportation: Transportation) {
    if (transportation.type != TransportationType.Plane)
      return <></>

    const flight = transportation.flightDetail!
    return flight.legs.map((leg, idx) =>
      <div key={idx}>
        <FlightEntry flight={flight} flightLeg={leg} onInfoBtnClick={() => onFlightClick(transportation)}/>
          {flight.legs.length > idx + 1 &&
          <span className="mx-3 text-sm text-muted-foreground">
            {formatDuration(leg.arrivalDateTime, flight.legs[idx+1].departureDateTime)} Layover
          </span>
        }
      </div>
    )
  }

  function renderNightTransportation(transportation: Transportation) {
    if (transportation.type != TransportationType.Plane)
      return <></>

    const flight = transportation.flightDetail!
    const nightLeg = flight.legs.find(leg => !isSameDay(leg.arrivalDateTime, dayData.day))!

    return (
      <FlightEntry flight={flight} flightLeg={nightLeg} onInfoBtnClick={() => onFlightClick(transportation)}/>
    )
  }

  return (
    <div>
      <DayLabel date={dayData.day}/>

      {dayData.activities.map(act =>
          <ActivityEntry key={act.id} activity={act} onClick={() => onActivityClick(act)}/>
      )}

      {dayTransportation.map((transportation, idx) =>
        <div key={idx}>
          {renderDayTransportation(transportation)}
        </div>
      )}

      {nextDay && (nightTransportation ?
        <>
          {renderNightTransportation(nightTransportation)}
          <Separator className="relative bottom-5 z-0"/>
        </>
      :
        <DaySeparator accomodation={dayData.accommodation} collapsedDays={collapsedDays} onAccommodationClick={onAccommodationClick}/>
      )}
    </div>
  )
}