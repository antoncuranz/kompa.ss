"use client"

import React from "react";
import "maplibre-gl/dist/maplibre-gl.css";
import {Map, MapProps} from "react-map-gl/maplibre";
import {useTheme} from "next-themes";
import {Coordinates} from "@/types.ts";
import RenderAfterMap from "@/components/card/RenderAfterMap.tsx";

export default function BaseMap({
  children, initialCoordinates, ...props
}: MapProps & {
  children: React.ReactNode | React.ReactNode[]
  initialCoordinates?: Coordinates|undefined,
}) {
  const fallbackLat = 52.520007
  const fallbackLon = 13.404954

  const {resolvedTheme} = useTheme()

  function getMapboxTheme() {
    return resolvedTheme == "dark" ? "night" : "day"
  }

  return (
    <Map
        mapStyle="https://api.maptiler.com/maps/basic-v2/style.json?key=SNIP"
        projection="globe"
        initialViewState={{latitude: initialCoordinates?.latitude ?? fallbackLat, longitude: initialCoordinates?.longitude ?? fallbackLon, zoom: 10}}
        // config={{"basemap": {"lightPreset": getMapboxTheme()}}}
        style={{background: "#04162a"}}
        {...props}
    >
      <RenderAfterMap theme={getMapboxTheme()}>
        {children}
      </RenderAfterMap>
    </Map>
  )
}