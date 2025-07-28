"use client"

import React from "react";
import Map from "react-map-gl/mapbox";
import "mapbox-gl/dist/mapbox-gl.css";

export default function MapCard() {
  const mapboxToken = process.env.NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN;

  return (
    <div className="flex-grow rounded-lg overflow-hidden border shadow-sm">
      <Map
        mapboxAccessToken={mapboxToken}
        mapStyle="mapbox://styles/mapbox/standard"
        projection="globe"
        initialViewState={{ latitude: 35.668641, longitude: 139.750567, zoom: 10 }}
      />
    </div>
  )
}