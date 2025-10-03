import React from "react";
import Card from "@/components/card/Card.tsx";
import {fetchAccommodation, fetchActivities, fetchGeoJson} from "@/requests.ts";
import HeroMap from "@/components/map/HeroMap.tsx";

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
      <HeroMap activities={activities} accommodation={accommodation} geojson={geojson}/>
    </Card>
  )
}