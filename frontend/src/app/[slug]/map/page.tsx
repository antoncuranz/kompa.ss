import React, {Suspense} from "react";
import SkeletonCard from "@/components/card/SkeletonCard.tsx";
import {ErrorBoundary} from "react-error-boundary";
import MapCard from "@/components/card/MapCard.tsx";

export default async function Page( {
  params
}: {
  params: Promise<{ slug: string }>
}) {
  const tripId = parseInt((await params).slug)

  return (
    <div className="flex h-full gap-4">
      <Suspense fallback={<SkeletonCard/>}>
        <ErrorBoundary fallback={<SkeletonCard title="Error loading Map"/>}>
          <MapCard tripId={tripId}/>
        </ErrorBoundary>
      </Suspense>
    </div>
  )
}