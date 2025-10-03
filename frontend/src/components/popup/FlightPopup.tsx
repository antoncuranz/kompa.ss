import React from "react";
import {formatTimePadded} from "@/components/util.ts";
import {GeoJsonPlane, GeoJsonFlightLeg} from "@/types.ts";

export default function FlightPopup({
  properties
}: {
  properties: GeoJsonPlane
}) {

  const legs = JSON.parse(properties.legs) as GeoJsonFlightLeg[]

  return (
    <div className="text-sm">
      <strong>✈️ Flight from {properties.fromMunicipality} to {properties.toMunicipality}</strong>
      <div className="iconsolata grid grid-cols-[auto_auto_1fr] gap-x-2">
        {legs.map((leg, idx) =>
          <div key={idx} className="grid col-span-3 grid-cols-subgrid">
            <span>
              {formatTimePadded(leg.departureDateTime)}-{formatTimePadded(leg.arrivalDateTime)}
            </span>
            <span className="text-right">
              {leg.flightNumber}
            </span>
            <span>
              {leg.fromIata}-{leg.toIata}
            </span>
          </div>
        )}
      </div>
    </div>
  )
}