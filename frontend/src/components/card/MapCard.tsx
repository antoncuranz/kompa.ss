import React from "react";
import "mapbox-gl/dist/mapbox-gl.css";
import Card from "@/components/card/Card.tsx";
import {fetchAccommodation, fetchActivities, fetchFlights} from "@/requests.ts";
import TheMap from "@/components/card/TheMap.tsx";

export default async function MapCard({
  tripId, className
}: {
  tripId: number
  className?: string
}) {
  const activities = await fetchActivities(tripId)
  const accommodation = await fetchAccommodation(tripId)
  const flights = await fetchFlights(tripId)

  return (
    <Card className={className}>
      <TheMap activities={activities} accommodation={accommodation} flights={flights}/>
    </Card>
  )
}