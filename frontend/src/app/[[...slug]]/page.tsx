import React from "react";
import {fetchTrips} from "@/requests.ts";
import Card from "@/components/card/Card.tsx";
import Link from "next/link";

export default async function Page() {

  const trips = await fetchTrips()

  return (
    <div className="flex h-full gap-4">
      {trips.map(trip =>
        <Link key={trip.id} href={"/" + trip.id + "/itinerary"}>
          <Card title={trip.name}>
            <div className="m-3">
              {trip.id}
            </div>
          </Card>
        </Link>
      )}
    </div>
  )
}