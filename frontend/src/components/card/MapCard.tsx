"use client"

import React from "react";
import "mapbox-gl/dist/mapbox-gl.css";
import {Layer, Map, Source} from "react-map-gl/mapbox";
import type {FeatureCollection, Feature} from 'geojson';
import RenderAfterMap from "@/components/card/RenderAfterMap.tsx";
import {useTheme} from "next-themes";
import {Flight, FlightLeg} from "@/types.ts";

export default function MapCard({
  flights
}: {
  flights: Flight[]
}) {
  const mapboxToken = process.env.NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN;

  const {resolvedTheme} = useTheme()

  function getConfig() {
    return {
      "basemap": {"lightPreset": resolvedTheme == "dark" ? "night" : "day"}
    }
  }

  function getGeoJson(): FeatureCollection {
    const features = flights
        .flatMap(flight => flight.legs)
        .map(mapLegToFeature)

    return {type: "FeatureCollection", features: features};
  }

  function mapLegToFeature(leg: FlightLeg): Feature {
    return {
      type: 'Feature',
      geometry: {
        type: 'LineString',
        coordinates: [
          [leg.origin.location!.longitude, leg.origin.location!.latitude],
          [leg.destination.location!.longitude, leg.destination.location!.latitude],
        ]
      },
      properties: {}
    }
  }

  return (
    <div className="flex-grow rounded-lg overflow-hidden border shadow-sm">
      <Map
          mapboxAccessToken={mapboxToken}
          mapStyle="mapbox://styles/mapbox/standard"
          projection="globe"
          initialViewState={{latitude: 52.520007, longitude: 13.404954, zoom: 10}}
          config={getConfig()}
      >
        <RenderAfterMap>
          <Source type="geojson" data={getGeoJson()}>
            <Layer type="line"
                   paint={{"line-color": "#007cbf", "line-width": 5}}
                   layout={{"line-cap": "round"}}
            />
            <Layer type="circle"
                   paint={{"circle-color": "#007cbf", "circle-radius": 5, "circle-stroke-color": "white", "circle-stroke-width": 3}}
            />
          </Source>
        </RenderAfterMap>
      </Map>
    </div>
  )
}