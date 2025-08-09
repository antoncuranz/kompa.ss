"use client"

import React, {useState} from "react";
import {Accommodation, Activity, DayRenderData, Flight, Trip} from "@/types.ts";
import Day from "@/components/itinerary/Day.tsx";
import ActivityDialog from "@/components/dialog/ActivityDialog.tsx";
import AccommodationDialog from "@/components/dialog/AccommodationDialog.tsx";
import FlightDialog from "@/components/dialog/FlightDialog.tsx";
import {useRouter} from "next/navigation";

export default function Itinerary({
  trip, dataByDays
}: {
  trip: Trip
  dataByDays: DayRenderData[]
}) {
  const [activityDialogOpen, setActivityDialogOpen] = useState(false)
  const [dialogActivity, setDialogActivity] = useState<Activity|null>()

  const [accommodationDialogOpen, setAccommodationDialogOpen] = useState(false)
  const [dialogAccommodation, setDialogAccommodation] = useState<Accommodation|null>()

  const [flightDialogOpen, setFlightDialogOpen] = useState(false)
  const [dialogFlight, setDialogFlight] = useState<Flight|null>()

  const router = useRouter()

  function onActivityClick(activity: Activity) {
    setDialogActivity(activity)
    setActivityDialogOpen(true)
  }

  function onActivityDialogClose(needsUpdate: boolean) {
    setActivityDialogOpen(false)
    if (needsUpdate)
      router.refresh()
  }

  function onAccommodationClick(accommodation: Accommodation | undefined) {
    setDialogAccommodation(accommodation)
    setAccommodationDialogOpen(true)
  }

  function onAccommodationDialogClose(needsUpdate: boolean) {
    setAccommodationDialogOpen(false)
    if (needsUpdate)
      router.refresh()
  }

  function onFlightClick(flight: Flight) {
    setDialogFlight(flight)
    setFlightDialogOpen(true)
  }

  function onFlightDialogClose(needsUpdate: boolean) {
    setFlightDialogOpen(false)
    if (needsUpdate)
      router.refresh()
  }

  return (
    <>
      {dataByDays.map((dayData, idx) =>
        <Day key={dayData.day.toISOString()} dayData={dayData} nextDay={dataByDays[idx+1]?.day}
             onActivityClick={onActivityClick}
             onAccommodationClick={onAccommodationClick}
             onFlightClick={onFlightClick}
        />
      )}
      <ActivityDialog trip={trip} key={"act-" + dialogActivity?.id} activity={dialogActivity} open={activityDialogOpen} onClose={onActivityDialogClose}/>
      <AccommodationDialog trip={trip} key={"acc-" + dialogAccommodation?.id} accommodation={dialogAccommodation} open={accommodationDialogOpen} onClose={onAccommodationDialogClose}/>
      <FlightDialog trip={trip} key={"flight-" + dialogFlight?.id} flight={dialogFlight} open={flightDialogOpen} onClose={onFlightDialogClose}/>
    </>
  )
}