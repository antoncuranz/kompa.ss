import FlightEntry from "@/components/itinerary/FlightEntry.tsx";
import DaySeparator from "@/components/itinerary/DaySeparator.tsx";
import {formatDuration, getDaysBetween, isSameDay} from "@/components/util.ts";
import DayLabel from "@/components/itinerary/DayLabel.tsx";
import React from "react";
import ActivityEntry from "@/components/itinerary/ActivityEntry.tsx";
import {Accommodation, Activity, DayRenderData, Transportation, TransportationType} from "@/types.ts";
import {Separator} from "@/components/ui/separator.tsx";
import TrainEntry from "@/components/itinerary/TrainEntry.tsx";
import TransportationEntry from "@/components/itinerary/TransportationEntry.tsx";

export default function Day({
  dayData,
  nextDay,
  onActivityClick = () => {},
  onAccommodationClick = () => {},
  onFlightClick = () => {},
  onTrainClick = () => {},
  onTransportationClick = () => {}
}: {
  dayData: DayRenderData
  nextDay: string
  onActivityClick?: (activity: Activity) => void;
  onAccommodationClick?: (accommodation: Accommodation | undefined) => void;
  onFlightClick?: (flight: Transportation) => void;
  onTrainClick?: (train: Transportation) => void;
  onTransportationClick?: (transportation: Transportation) => void;
}) {

  const collapsedDays = nextDay ? getDaysBetween(dayData.day, nextDay).length-2 : 0

  const hasNightTransportation = dayData.transportation
    .find(t => !isSameDay(t.arrivalDateTime, dayData.day)) != undefined

  function renderTransportation(transportation: Transportation) {
    switch (transportation.type) {
      case TransportationType.Plane: // deprecated
      case TransportationType.Flight:
        return renderFlight(transportation)
      case TransportationType.Train:
        return renderTrain(transportation)
      default:
        return isSameDay(transportation.departureDateTime, dayData.day) &&
            <TransportationEntry transportation={transportation} onClick={() => onTransportationClick(transportation)}/>
    }
  }

  function renderFlight(transportation: Transportation) {
    const flight = transportation.flightDetail!
    const filteredLegs = flight.legs
        .filter(leg => isSameDay(leg.departureDateTime, dayData.day))

    return filteredLegs.map((leg, idx) =>
      <div key={idx}>
        <FlightEntry flight={flight} flightLeg={leg} onInfoBtnClick={() => onFlightClick(transportation)}/>
        {filteredLegs.length > idx + 1 &&
          <span className="mx-3 text-sm text-muted-foreground">
            {formatDuration(leg.arrivalDateTime, filteredLegs[idx+1].departureDateTime)} Layover
          </span>
        }
      </div>
    )
  }

  function renderTrain(transportation: Transportation) {
    const train = transportation.trainDetail!
    const filteredLegs = train.legs
        .filter(leg => isSameDay(leg.departureDateTime, dayData.day))

    return filteredLegs.map((leg, idx) =>
        <div key={idx}>
          <TrainEntry trainLeg={leg} onInfoBtnClick={() => onTrainClick(transportation)}/>
          {filteredLegs.length > idx + 1 &&
            <span className="mx-3 text-sm text-muted-foreground">
              {formatDuration(leg.arrivalDateTime, filteredLegs[idx+1].departureDateTime)} Layover
            </span>
          }
        </div>
    )
  }

  return (
    <div>
      <DayLabel date={dayData.day}/>

      {dayData.activities.map(act =>
          <ActivityEntry key={act.id} activity={act} onClick={() => onActivityClick(act)}/>
      )}

      {dayData.transportation.map((transportation, idx) =>
        <div key={idx} className="mt-4">
          {renderTransportation(transportation)}
        </div>
      )}

      {nextDay && (hasNightTransportation ?
        <Separator className="relative bottom-5 z-0"/>
      :
        <DaySeparator className="mt-4" accomodation={dayData.accommodation} collapsedDays={collapsedDays} onAccommodationClick={onAccommodationClick}/>
      )}
    </div>
  )
}