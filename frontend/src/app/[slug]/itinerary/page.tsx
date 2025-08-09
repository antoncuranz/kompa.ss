import React, {Suspense} from "react";
import SkeletonCard from "@/components/card/SkeletonCard.tsx";
import {ErrorBoundary} from "react-error-boundary";
import ItineraryCard from "@/components/card/ItineraryCard.tsx";
import MapCard from "@/components/card/MapCard.tsx";

export default async function Page( {
  params
}: {
  params: Promise<{ slug: string }>
}) {
  const tripId = parseInt((await params).slug)

  const itineraryClasses = "lg:max-w-[48rem] lg:min-w-[38rem]"
  const mapClasses = "hidden lg:block"

  return (
    <div className="flex h-full gap-4">
      <Suspense fallback={<SkeletonCard className={itineraryClasses}/>}>
        <ErrorBoundary fallback={<SkeletonCard className={itineraryClasses} title="Error loading Itinerary"/>}>
          <ItineraryCard tripId={tripId} className={itineraryClasses}/>
        </ErrorBoundary>
      </Suspense>
      <Suspense fallback={<SkeletonCard className={mapClasses}/>}>
        <ErrorBoundary fallback={<SkeletonCard className={mapClasses} title="Error loading Map"/>}>
          <MapCard tripId={tripId} className={mapClasses}/>
        </ErrorBoundary>
      </Suspense>
    </div>
  )
}