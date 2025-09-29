"use client"

import React, {useState} from "react";
import "mapbox-gl/dist/mapbox-gl.css";
import {Layer, Map, Popup, Source} from "react-map-gl/mapbox";
import type {Feature, FeatureCollection, GeoJsonProperties} from 'geojson';
import RenderAfterMap from "@/components/card/RenderAfterMap.tsx";
import {Accommodation, Activity, GeoJsonPlane, GeoJsonTrain} from "@/types.ts";
import {LngLat, MapMouseEvent} from "mapbox-gl";
import {formatDateShort, formatTime} from "@/components/util.ts";
import {useTheme} from "next-themes";
import TrainPopup from "@/components/popup/TrainPopup.tsx";
import FlightPopup from "@/components/popup/FlightPopup.tsx";

export default function TheMap({
  activities, accommodation, geojson
}: {
  activities: Activity[],
  accommodation: Accommodation[],
  geojson: GeoJSON.FeatureCollection[],
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

  function onMouseEnter(event: MapMouseEvent) {
    if (!event.features || event.features.length == 0)
      return

    const featureProperties = event.features.filter(feature =>
        feature.properties && feature.properties["type"]
    ).map(feature => feature.properties)

    setPopupInfo({
      lngLat: event.lngLat,
      children: featureProperties.map((properties, idx) =>
        <div key={idx}>
          {renderTransportationProperties(properties)}
        </div>
      )
    })
  }

  function renderTransportationProperties(props: GeoJsonProperties) {
    switch (props!["type"]) {
      case "PLANE":
        return <FlightPopup properties={props as GeoJsonPlane}/>
      case "TRAIN":
        return <TrainPopup properties={props as GeoJsonTrain}/>
      default:
        return <></>
    }
  }

  function getColorByType(fc: FeatureCollection) {
    // @ts-expect-error custom property
    const type = fc["transportationType"] as string

    switch (type) {
      case "PLANE": return "#007cbf"
      case "TRAIN": return "#ec0016"
      default: return "black"
    }
  }

  return (
    <Map
        mapboxAccessToken={mapboxToken}
        mapStyle="mapbox://styles/mapbox/standard"
        projection="globe"
        initialViewState={{latitude: 52.520007, longitude: 13.404954, zoom: 10}}
        config={{"basemap": {"lightPreset": getMapboxTheme()}}}
        interactiveLayerIds={geojson.map((_, idx) => "geojson" + idx).concat(["activity", "accommodation", "flight"])}
        onMouseEnter={onMouseEnter}
        onMouseMove={onMouseEnter}
        onMouseLeave={() => setPopupInfo(null)}
    >
      <RenderAfterMap theme={getMapboxTheme()}>
        {geojson.map((fc, idx) =>
            <Source key={idx} type="geojson" data={fc}>
              <Layer type="line"
                     paint={{"line-color": getColorByType(fc), "line-width": 5, "line-emissive-strength": 1}}
                     layout={{"line-cap": "round"}}
              />
              <Layer type="circle" id={"geojson" + idx}
                     filter={["==", ["geometry-type"], "Point"]}
                     paint={{"circle-color": getColorByType(fc), "circle-radius": 5, "circle-stroke-color": "white", "circle-stroke-width": 3, "circle-emissive-strength": 1}}
              />
            </Source>
        )}
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
        {popupInfo && (
            <Popup offset={10}
                   closeButton={false}
                   closeOnClick={false}
                   longitude={popupInfo.lngLat.lng}
                   latitude={popupInfo.lngLat.lat}
                   maxWidth={undefined}
                   className="shadow-xl bg-background"
            >
              {popupInfo.children}
            </Popup>
        )}
      </RenderAfterMap>
    </Map>
  )
}