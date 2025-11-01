import {
  Layer as MaplibreLayer,
  MapMouseEvent as MaplibreMouseEvent,
  MapProps as MaplibreMapProps,
  MapRef as MaplibreRef,
  Marker as MaplibreMarker,
  MarkerProps as MaplibreMarkerProps,
  Popup as MaplibrePopup,
  PopupProps as MaplibrePopupProps,
  Source as MaplibreSource,
  useMap as useMaplibreMap
} from "react-map-gl/maplibre";
import {
  Layer as MapboxLayer,
  MapMouseEvent as MapboxMouseEvent,
  MapProps as MapboxMapProps,
  MapRef as MapboxRef,
  Marker as MapboxMarker,
  MarkerProps as MapboxMarkerProps,
  Popup as MapboxPopup,
  PopupProps as MapboxPopupProps,
  Source as MapboxSource,
  useMap as useMapboxMap
} from "react-map-gl/mapbox";
import {SharedProperties} from "@/types.ts";
import type {FeatureCollection} from "geojson";
import React, {ReactNode} from "react";
import {LngLat as MapboxLngLat} from "mapbox-gl";
import {LngLat as MaplibreLngLat} from "maplibre-gl";

export const mapboxToken = process.env.NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN;
export const mapStyle = process.env.NEXT_PUBLIC_MAP_STYLE;
export const isMapbox: boolean = !!mapboxToken

export type MapMouseEvent = MapboxMouseEvent | MaplibreMouseEvent
export type MapProps = SharedProperties<MaplibreMapProps, MapboxMapProps>
export type LngLat = MapboxLngLat|MaplibreLngLat

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
  layout?: any
  filter?: any
}) {
  return isMapbox ? <MapboxLayer {...props}/> : <MaplibreLayer {...props}/>
}

export function Popup(props: SharedProperties<MaplibrePopupProps, MapboxPopupProps>) {
  return isMapbox ? <MapboxPopup {...props}/> : <MaplibrePopup {...props}/>
}

export type MapRef = SharedProperties<MaplibreRef, MapboxRef>

export type MapCollection = {
  [id: string]: MapRef | undefined;
  current?: MapRef;
};

export function useMap(): MapCollection {
  // eslint-disable-next-line react-hooks/rules-of-hooks
  return isMapbox ? useMapboxMap() as MapCollection : useMaplibreMap() as MapCollection
}
