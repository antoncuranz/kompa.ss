"use client"

import React, {useState} from "react";
import {Accommodation, Activity, DayRenderData, Transportation, Trip} from "@/types.ts";
import Day from "@/components/itinerary/Day.tsx";
import AccommodationDialogContent from "@/components/dialog/AccommodationDialogContent.tsx";
import FlightDialogContent from "@/components/dialog/FlightDialogContent.tsx";
import ActivityDialogContent from "@/components/dialog/ActivityDialogContent.tsx";
import {Dialog} from "@/components/dialog/Dialog.tsx";
import TrainDialogContent from "@/components/dialog/TrainDialogContent.tsx";
import TransportationDialogContent from "@/components/dialog/TransportationDialogContent.tsx";

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
  const [dialogFlight, setDialogFlight] = useState<Transportation|null>()

  const [trainDialogOpen, setTrainDialogOpen] = useState(false)
  const [dialogTrain, setDialogTrain] = useState<Transportation|null>()

  const [transportationDialogOpen, setTransportationDialogOpen] = useState(false)
  const [dialogTransportation, setDialogTransportation] = useState<Transportation|null>()

  function onActivityClick(activity: Activity) {
    setDialogActivity(activity)
    setActivityDialogOpen(true)
  }

  function onAccommodationClick(accommodation: Accommodation | undefined) {
    setDialogAccommodation(accommodation)
    setAccommodationDialogOpen(true)
  }

  function onFlightClick(flight: Transportation) {
    setDialogFlight(flight)
    setFlightDialogOpen(true)
  }

  function onTrainClick(train: Transportation) {
    setDialogTrain(train)
    setTrainDialogOpen(true)
  }

  function onTransportationClick(transportation: Transportation) {
    setDialogTransportation(transportation)
    setTransportationDialogOpen(true)
  }

  return (
    <div className="mb-4">
      {dataByDays.map((dayData, idx) =>
        <Day key={dayData.day} dayData={dayData} nextDay={dataByDays[idx+1]?.day}
             onActivityClick={onActivityClick}
             onAccommodationClick={onAccommodationClick}
             onFlightClick={onFlightClick}
             onTrainClick={onTrainClick}
             onTransportationClick={onTransportationClick}
        />
      )}
      <Dialog open={activityDialogOpen} setOpen={setActivityDialogOpen}>
        <ActivityDialogContent trip={trip} activity={dialogActivity}/>
      </Dialog>
      <Dialog open={accommodationDialogOpen} setOpen={setAccommodationDialogOpen}>
        <AccommodationDialogContent trip={trip} accommodation={dialogAccommodation}/>
      </Dialog>
      <Dialog open={flightDialogOpen} setOpen={setFlightDialogOpen}>
        <FlightDialogContent trip={trip} flight={dialogFlight}/>
      </Dialog>
      <Dialog open={trainDialogOpen} setOpen={setTrainDialogOpen}>
        <TrainDialogContent trip={trip} train={dialogTrain}/>
      </Dialog>
      <Dialog open={transportationDialogOpen} setOpen={setTransportationDialogOpen}>
        <TransportationDialogContent trip={trip} transportation={dialogTransportation}/>
      </Dialog>
    </div>
  )
}