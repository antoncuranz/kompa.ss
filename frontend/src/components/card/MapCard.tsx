import React from "react";
import "mapbox-gl/dist/mapbox-gl.css";
import Card from "@/components/card/Card.tsx";
import {fetchAccommodation, fetchActivities, fetchGeoJson} from "@/requests.ts";
import TheMap from "@/components/card/TheMap.tsx";

export default async function MapCard({
  tripId, className
}: {
  tripId: number
  className?: string
}) {
  const activities = await fetchActivities(tripId)
  const accommodation = await fetchAccommodation(tripId)
  const geojson = await fetchGeoJson(tripId)

  return (
    <Card className={className}>
      <TheMap activities={activities} accommodation={accommodation} geojson={geojson}/>
    </Card>
  )
}