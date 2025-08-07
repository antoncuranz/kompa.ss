import React, {Suspense} from "react";
import SkeletonCard from "@/components/card/SkeletonCard.tsx";
import {ErrorBoundary} from "react-error-boundary";
import ItineraryCard from "@/components/card/ItineraryCard.tsx";
import MapCard from "@/components/card/MapCard.tsx";
import {fetchAccommodation, fetchActivities, fetchFlights, fetchTrip} from "@/requests.ts";

export default async function Page() {

  const tripId = 1
  const trip = await fetchTrip(tripId)
  const activities = await fetchActivities()
  const accomodation = await fetchAccommodation()
  const flights = await fetchFlights()

  return (
    <div className="flex h-full gap-4">
      <Suspense fallback={<SkeletonCard/>}>
        <ErrorBoundary fallback={<SkeletonCard title="Error loading Itinerary"/>}>
          <ItineraryCard trip={trip} activities={activities} accomodation={accomodation} flights={flights}/>
        </ErrorBoundary>
      </Suspense>
      <Suspense fallback={<SkeletonCard/>}>
        <ErrorBoundary fallback={<SkeletonCard title="Error loading Map"/>}>
          <MapCard activities={activities} accommodation={accomodation} flights={flights}/>
        </ErrorBoundary>
      </Suspense>
    </div>
  )
}