"use client"

import React from "react"
import { Coordinates } from "@/types.ts"
import BaseMap from "@/components/map/BaseMap.tsx"
import {
  LngLat,
  MapMouseEvent,
  Marker,
  useMap,
} from "@/components/map/common.tsx"

export default function MiniMap({
  children,
  value,
  onChange,
}: {
  children: React.ReactNode | React.ReactNode[]
  value: Coordinates | undefined
  onChange: (newLocation: Coordinates) => void
}) {
  const { heroMap } = useMap()

  function getInitialCoordinates() {
    if (value) {
      return value
    } else if (heroMap) {
      return lngLatToCoordinates(heroMap.getCenter())
    } else {
      return undefined
    }
  }

  function onClick(event: MapMouseEvent) {
    onChange(lngLatToCoordinates(event.lngLat))
  }

  function lngLatToCoordinates(lngLat: LngLat) {
    return {
      latitude: lngLat.lat,
      longitude: lngLat.lng,
    }
  }

  return (
    <BaseMap initialCoordinates={getInitialCoordinates()} onClick={onClick}>
      {value && (
        <Marker longitude={value.longitude} latitude={value.latitude} />
      )}
      {children}
    </BaseMap>
  )
}
