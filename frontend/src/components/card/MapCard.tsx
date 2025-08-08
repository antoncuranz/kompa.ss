"use client"

import React, {useState} from "react";
import "mapbox-gl/dist/mapbox-gl.css";
import {Layer, Map, Source} from "react-map-gl/mapbox";
import type {FeatureCollection, Feature} from 'geojson';
import RenderAfterMap from "@/components/card/RenderAfterMap.tsx";
import {Accommodation, Activity, Flight, FlightLeg} from "@/types.ts";
import {Popup} from "react-map-gl/mapbox";
import {LngLat, MapMouseEvent} from "mapbox-gl";
import {formatDateShort, formatTime, isSameDay} from "@/components/util.ts";
import {GlowContainer} from "@/components/ui/glow-container.tsx";
import {useTheme} from "next-themes";

export default function MapCard({
  activities, accommodation, flights
}: {
  activities: Activity[],
  accommodation: Accommodation[],
  flights: Flight[]
}) {
  const mapboxToken = process.env.NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN;

  const {resolvedTheme} = useTheme()

  function getMapboxTheme() {
    return resolvedTheme == "dark" ? "night" : "day"
  }

  type PopupInfo = {
    lngLat: LngLat;
    children: React.ReactNode;
  }
  const [popupInfo, setPopupInfo] = useState<PopupInfo|null>(null)

  function getActivityGeoJson(): FeatureCollection {
    const features = activities
        .filter(activity => activity.location)
        .map(mapActivityToFeature)

    return {type: "FeatureCollection", features: features};
  }

  function mapActivityToFeature(activity: Activity): Feature {
    return {
      type: 'Feature',
      geometry: {
        type: 'Point',
        coordinates: [activity.location!.longitude, activity.location!.latitude],
      },
      properties: {
        "popupTitle": activity.name,
        "popupBody": formatDateShort(activity.date) + (activity.time ? " " + formatTime(activity.time) : "")
      }
    }
  }

  function getAccommodationGeoJson(): FeatureCollection {
    const features = accommodation
        .filter(activity => activity.location)
        .map(mapAccommodationToFeature)

    return {type: "FeatureCollection", features: features};
  }

  function mapAccommodationToFeature(accommodation: Accommodation): Feature {
    return {
      type: 'Feature',
      geometry: {
        type: 'Point',
        coordinates: [accommodation.location!.longitude, accommodation.location!.latitude],
      },
      properties: {
        "popupTitle": `🛏️ ${accommodation.name}`,
        "popupBody": `${formatDateShort(accommodation.arrivalDate)} - ${formatDateShort(accommodation.departureDate)}`,
      }
    }
  }

  function getFlightGeoJson(): FeatureCollection {
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
      properties: {
        "popupTitle": `✈️ Flight from ${leg.origin.municipality} to ${leg.destination.municipality}`,
        "popupBody": leg.flightNumber,
        "popupBodyRight": `${formatDateShort(leg.departureDateTime)} ${formatTime(leg.departureDateTime)} - ${formatTime(leg.arrivalDateTime)}` + (!isSameDay(leg.departureDateTime, leg.arrivalDateTime) ? " (+1)" : "")
      }
    }
  }

  function onMouseEnter(event: MapMouseEvent) {
    if (!event.features || event.features.length == 0)
      return

    const featureProperties = event.features.filter(feature =>
        feature.properties && feature.properties["popupTitle"] && feature.properties["popupBody"]
    ).map(feature => feature.properties)

    setPopupInfo({
      lngLat: event.lngLat,
      children: featureProperties.map((properties, idx) =>
          <div key={idx}>
            <strong>{properties!["popupTitle"]}</strong>
            {properties!["popupBody"] &&
              <div className="flex">
                <p className="flex-grow">{properties!["popupBody"]}</p>
                {properties!["popupBodyRight"] &&
                  <p className="ml-2">{properties!["popupBodyRight"]}</p>
                }
              </div>
            }
          </div>
      )
    })
  }

  return (
    <div className="flex-grow rounded-3xl shadow-xl shadow-black/[0.1] dark:shadow-white/[0.05] boder border-black/10 dark:border-white/20">
      <GlowContainer className="rounded-3xl h-full p-2">
        <div className="rounded-2xl h-full overflow-hidden">
          <Map
              mapboxAccessToken={mapboxToken}
              mapStyle="mapbox://styles/mapbox/standard"
              projection="globe"
              initialViewState={{latitude: 52.520007, longitude: 13.404954, zoom: 10}}
              config={{"basemap": {"lightPreset": getMapboxTheme()}}}
              interactiveLayerIds={["activity", "accommodation", "flight"]}
              onMouseEnter={onMouseEnter}
              onMouseMove={onMouseEnter}
              onMouseLeave={() => setPopupInfo(null)}
          >
            <RenderAfterMap theme={getMapboxTheme()}>
              <Source type="geojson" data={getAccommodationGeoJson()}>
                <Layer id="accommodation"
                       type="circle"
                       paint={{"circle-color": "#f1b216", "circle-radius": 5, "circle-stroke-color": "white", "circle-stroke-width": 3, "circle-emissive-strength": 1}}
                />
              </Source>
              <Source type="geojson" data={getActivityGeoJson()}>
                <Layer id="activity"
                       type="circle"
                       paint={{"circle-color": "#36bf00", "circle-radius": 5, "circle-stroke-color": "white", "circle-stroke-width": 3, "circle-emissive-strength": 1}}
                />
              </Source>
              <Source type="geojson" data={getFlightGeoJson()}>
                <Layer type="line"
                       paint={{"line-color": "#007cbf", "line-width": 5, "line-emissive-strength": 1}}
                       layout={{"line-cap": "round"}}
                />
                <Layer id="flight"
                       type="circle"
                       paint={{"circle-color": "#007cbf", "circle-radius": 5, "circle-stroke-color": "white", "circle-stroke-width": 3, "circle-emissive-strength": 1}}
                />
              </Source>
              {popupInfo && (
                  <Popup offset={10}
                         closeButton={false}
                         closeOnClick={false}
                         longitude={popupInfo.lngLat.lng}
                         latitude={popupInfo.lngLat.lat}
                         className="shadow-xl bg-background"
                  >
                    {popupInfo.children}
                  </Popup>
              )}
            </RenderAfterMap>
          </Map>
        </div>
      </GlowContainer>
    </div>
  )
}