import React, {Suspense} from "react";
import SkeletonCard from "@/components/card/SkeletonCard.tsx";
import {ErrorBoundary} from "react-error-boundary";
import ItineraryCard from "@/components/card/ItineraryCard.tsx";

export default async function Page( {
  params
}: {
  params: Promise<{ slug: string }>
}) {
  const tripId = parseInt((await params).slug)

  return (
    <div className="flex h-full gap-4">
      <Suspense fallback={<SkeletonCard/>}>
        <ErrorBoundary fallback={<SkeletonCard title="Error loading Itinerary"/>}>
          <ItineraryCard tripId={tripId}/>
        </ErrorBoundary>
      </Suspense>
    </div>
  )
}