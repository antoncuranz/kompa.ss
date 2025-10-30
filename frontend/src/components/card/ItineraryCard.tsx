"use client"

import {DayRenderData, LoadedTransportation, RESOLVE_FLIGHT, RESOLVE_GENERIC_TRANSPORTATION, RESOLVE_TRAIN, RESOLVE_TRIP} from "@/schema.ts";
import {dayIsBetween, getDaysBetween, isSameDay} from "@/components/util.ts";
import React, { useEffect, useState } from "react";
import AddSomethingDropdown from "@/components/buttons/AddSomethingDropdown.tsx";
import Card from "@/components/card/Card.tsx";
import Itinerary from "@/components/itinerary/Itinerary.tsx";
import {useCoState} from "jazz-tools/react-core";
import {Trip} from "@/schema.ts";
import SkeletonCard from "@/components/card/SkeletonCard.tsx";

export default function ItineraryCard({
  tripId, className
}: {
  tripId: string
  className?: string
}) {
  const trip = useCoState(Trip, tripId, {resolve: RESOLVE_TRIP});
  const [dataByDays, setDataByDays] = useState<DayRenderData[]>([])

  function getDepartureDateTime(transportation: LoadedTransportation): string {
    switch (transportation.type) {
      case "flight":
        return transportation.legs[0]!.departureDateTime

      case "train":
        return transportation.legs[0]!.departureDateTime
      
      case "generic":
        return transportation.departureDateTime
    }
  }

  function getArrivalDateTime(transportation: LoadedTransportation): string {
    switch (transportation.type) {
      case "flight":
        return transportation.legs[0]!.arrivalDateTime

      case "train":
        return transportation.legs[0]!.arrivalDateTime
      
      case "generic":
        return transportation.arrivalDateTime
    }
  }

  function processDataAndGroupByDays(transportation: LoadedTransportation[]) {
    const grouped: DayRenderData[] = []

    for (const day of getDaysBetween(trip!.startDate, trip!.endDate)) {
      const filteredActivities = trip!.activities
          .filter(act => isSameDay(day, act.date))
        
      const relevantTransportation = transportation
          .filter(t => dayIsBetween(day, getDepartureDateTime(t), getArrivalDateTime(t)))

      const relevantAccommodation = trip!.accommodation.find(acc =>
          acc.arrivalDate <= day && acc.departureDate > day
      )

      // TODO: also push if day is today!
      if (isSameDay(day, trip!.endDate) || grouped.length == 0 || relevantTransportation.length != 0 ||
          filteredActivities.length != 0 || relevantAccommodation != grouped[grouped.length-1].accommodation ||
          grouped[grouped.length-1].transportation.find(t => isSameDay(getArrivalDateTime(t), day))
      ) {
        grouped.push({
          day: day,
          transportation: relevantTransportation,
          activities: filteredActivities,
          accommodation: relevantAccommodation,
        })
      }
    }

    setDataByDays(grouped)
  }

  useEffect(() => {
    async function todo() {
      if (!trip) {
        return
      }

      const loaded = await Promise.all(trip.transportation.map(async transportation => {
        switch (transportation.type) {
          case "flight":
            return await transportation.$jazz.ensureLoaded({resolve: RESOLVE_FLIGHT})

          case "train":
            return await transportation.$jazz.ensureLoaded({resolve: RESOLVE_TRAIN})

          case "generic":
            return await transportation.$jazz.ensureLoaded({resolve: RESOLVE_GENERIC_TRANSPORTATION})
        }
      }))

      processDataAndGroupByDays(loaded)
    }

    todo();
  }, [trip]);

  if (!trip) {
    return <SkeletonCard className={className} title={trip === null ? "Error loading Itinerary" : undefined}/>;
  }

  return (
    <Card title="Itinerary" headerSlot={<AddSomethingDropdown trip={trip}/>} className={className}>
      <Itinerary trip={trip} dataByDays={dataByDays}/>
    </Card>
  )
}