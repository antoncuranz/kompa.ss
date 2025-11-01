"use client"

import React from "react"
import Card from "@/components/card/Card.tsx"
import HeroMap from "@/components/map/HeroMap.tsx"
import { useCoState } from "jazz-tools/react-core"
import { RESOLVE_TRIP, Trip } from "@/schema.ts"
import SkeletonCard from "@/components/card/SkeletonCard.tsx"

export default function MapCard({
  tripId,
  className,
}: {
  tripId: string
  className?: string
}) {
  const trip = useCoState(Trip, tripId, { resolve: RESOLVE_TRIP })

  if (!trip) {
    return (
      <SkeletonCard
        className={className}
        title={trip === null ? "Error loading Map" : undefined}
      />
    )
  }

  return (
    <Card className={className}>
      <HeroMap
        activities={trip.activities.filter(act => act !== null)}
        accommodation={trip.accommodation.filter(acc => acc !== null)}
        geojson={[]}
      />
    </Card>
  )
}
