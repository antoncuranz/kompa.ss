import {
  Layer as MaplibreLayer,
  MapMouseEvent as MaplibreMouseEvent, MapProps as MaplibreMapProps,
  Marker as MaplibreMarker,
  MarkerProps as MaplibreMarkerProps, Popup as MaplibrePopup, PopupProps as MaplibrePopupProps, Source as MaplibreSource
} from "react-map-gl/maplibre";
import {
  Layer as MapboxLayer,
  MapMouseEvent as MapboxMouseEvent, MapProps as MapboxMapProps,
  Marker as MapboxMarker,
  MarkerProps as MapboxMarkerProps, Popup as MapboxPopup, PopupProps as MapboxPopupProps, Source as MapboxSource
} from "react-map-gl/mapbox";
import {SharedProperties} from "@/types.ts";
import type {FeatureCollection} from "geojson";
import React, {ReactNode} from "react";

export const mapboxToken = process.env.NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN;
export const mapStyle = process.env.NEXT_PUBLIC_MAP_STYLE;
export const isMapbox: boolean = !!mapboxToken

export type MapMouseEvent = MapboxMouseEvent | MaplibreMouseEvent
export type MapProps = SharedProperties<MaplibreMapProps, MapboxMapProps>

export function Marker(props: SharedProperties<MaplibreMarkerProps, MapboxMarkerProps>) {
  return isMapbox ? <MapboxMarker {...props}/> : <MaplibreMarker {...props}/>
}

export function Source(props: {
  type: "geojson" | "vector"
  data: FeatureCollection
  children: ReactNode | ReactNode[]
}) {
  return isMapbox ? <MapboxSource {...props}/> : <MaplibreSource {...props}/>
}

export function Layer(props: {
  id?: string
  type: "line" | "circle"
  paint: {
    "circle-color"?: string
    "circle-radius"?: number
    "circle-stroke-color"?: string
    "circle-stroke-width"?: number
    "line-color"?: string
    "line-width"?: number
  }
  /* eslint-disable @typescript-eslint/no-explicit-any */
  layout?: any
  filter?: any
}) {
  return isMapbox ? <MapboxLayer {...props}/> : <MaplibreLayer {...props}/>
}

export function Popup(props: SharedProperties<MaplibrePopupProps, MapboxPopupProps>) {
  return isMapbox ? <MapboxPopup {...props}/> : <MaplibrePopup {...props}/>
}
