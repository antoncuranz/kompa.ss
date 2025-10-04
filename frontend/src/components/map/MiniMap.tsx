"use client"

import React from "react";
import {Marker} from "react-map-gl/mapbox";
import {MapMouseEvent} from "mapbox-gl";
import {Coordinates} from "@/types.ts";
import BaseMap from "@/components/map/BaseMap.tsx";

export default function MiniMap({
  children, value, onChange
}: {
  children: React.ReactNode | React.ReactNode[]
  value: Coordinates|undefined,
  onChange: (newLocation: Coordinates) => void,
}) {

  function onClick(event: MapMouseEvent) {
    onChange({
      latitude: event.lngLat.lat,
      longitude: event.lngLat.lng,
    })
  }

  return (
    <BaseMap initialCoordinates={value} onClick={onClick}>
      {value && <Marker longitude={value.longitude} latitude={value.latitude}/>}
      {children}
    </BaseMap>
  )
}