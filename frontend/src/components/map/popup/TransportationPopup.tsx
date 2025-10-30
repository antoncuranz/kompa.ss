import React from "react";
import {formatDateShort, formatTime, isSameDay} from "@/components/util.ts";
import {GeoJsonTransportation, getTransportationTypeEmoji} from "@/schema.ts";

export default function TransportationPopup({
  properties
}: {
  properties: GeoJsonTransportation
}) {

  return (
    <div className="text-sm">
      <strong>
        {getTransportationTypeEmoji(properties.type)} {properties.name}
      </strong>
      <div className="iconsolata">
        {isSameDay(properties.departureDateTime, properties.arrivalDateTime) ?
          `${formatDateShort(properties.departureDateTime)} ${formatTime(properties.departureDateTime)}-${formatTime(properties.arrivalDateTime)}`
        :
          `${formatDateShort(properties.departureDateTime)} ${formatTime(properties.departureDateTime)} - ${formatDateShort(properties.arrivalDateTime)} ${formatTime(properties.arrivalDateTime)}`
        }
      </div>
    </div>
  )
}