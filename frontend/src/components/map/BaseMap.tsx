"use client"

import React from "react";
import "mapbox-gl/dist/mapbox-gl.css";
import {Map, MapProps} from "react-map-gl/mapbox";
import {useTheme} from "next-themes";
import {Coordinates} from "@/types.ts";
import RenderAfterMap from "@/components/card/RenderAfterMap.tsx";

export default function BaseMap({
  children, initialCoordinates, ...props
}: MapProps & {
  children: React.ReactNode | React.ReactNode[]
  initialCoordinates?: Coordinates|undefined,
}) {
  const mapboxToken = process.env.NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN;
  const fallbackLat = 52.520007
  const fallbackLon = 13.404954

  const {resolvedTheme} = useTheme()

  function getMapboxTheme() {
    return resolvedTheme == "dark" ? "night" : "day"
  }

  return (
    <Map
        mapboxAccessToken={mapboxToken}
        mapStyle="mapbox://styles/mapbox/standard"
        projection="globe"
        initialViewState={{latitude: initialCoordinates?.latitude ?? fallbackLat, longitude: initialCoordinates?.longitude ?? fallbackLon, zoom: 10}}
        config={{"basemap": {"lightPreset": getMapboxTheme()}}}
        {...props}
    >
      <RenderAfterMap theme={getMapboxTheme()}>
        {children}
      </RenderAfterMap>
    </Map>
  )
}