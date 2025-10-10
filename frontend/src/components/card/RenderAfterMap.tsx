"use client"

import React, {useEffect, useState} from "react";
import {MapRef as MaplibreRef, useMap as useMaplibreMap} from "react-map-gl/maplibre";
import {MapRef as MapboxRef, useMap as useMapboxMap} from "react-map-gl/mapbox";
import {isMapbox} from "@/components/map/common.tsx";
import {SharedProperties} from "@/types.ts";

type MapRef = SharedProperties<MaplibreRef, MapboxRef>

export type MapCollection = {
  [id: string]: MapRef | undefined;
  current?: MapRef;
};

export default function RenderAfterMap({
  children, theme
}: {
  theme?: string,
  children: React.ReactNode,
}) {
  // eslint-disable-next-line react-hooks/rules-of-hooks
  const map: MapCollection = isMapbox ? useMapboxMap() as MapCollection : useMaplibreMap() as MapCollection
  const [canRender, setCanRender] = useState(false)

  if (isMapbox) {
    // eslint-disable-next-line react-hooks/rules-of-hooks
    useEffect(() => {
      if (!map.current || !theme)
        return

      // @ts-expect-error i know
      const mapboxRef = map.current as MapboxRef
      mapboxRef.setConfigProperty("basemap", "lightPreset", theme);
    }, [map, theme])
  }

  useEffect(() => {
    map.current?.on('load', () => setCanRender(true))
  }, [map])

  return <>{canRender && children}</>
}
