"use client"

import * as React from "react"
import { MapProvider as MapGlMapProvider } from "react-map-gl/mapbox"

export function MapProvider({
  children,
  ...props
}: React.ComponentProps<typeof MapGlMapProvider>) {
  return <MapGlMapProvider {...props}>{children}</MapGlMapProvider>
}