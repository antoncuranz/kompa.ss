import {
    Accommodation,
    Activity,
    Transportation,
    Trip,
} from "@/types.ts";
import {headers} from "next/headers";

class UnauthorizedError extends Error {}


async function fetchData(url: string) {
  const authorizationHeader = "Authorization"
  const response = await fetch(process.env.BACKEND_URL + url, {
    headers: {[authorizationHeader]: (await headers()).get(authorizationHeader)!},
    cache: "no-cache"
  })
  if (response.status == 401) {
    throw new UnauthorizedError()
  } else if (!response.ok) {
    throw new Error("Failed to fetch data");
  }
  return await response.json()
}

export async function fetchTrip(tripId: number) {
  return await fetchData("/api/v1/trips/" + tripId) as Trip
}

export async function fetchTrips() {
  return await fetchData("/api/v1/trips") as Trip[]
}

export async function fetchAccommodation(tripId: number) {
  return await fetchData("/api/v1/trips/" + tripId + "/accommodation") as Accommodation[]
}

export async function fetchTransportation(tripId: number) {
  return await fetchData("/api/v1/trips/" + tripId + "/transportation") as Transportation[]
}

export async function fetchActivities(tripId: number) {
  return await fetchData("/api/v1/trips/" + tripId + "/activities") as Activity[];
}

export async function fetchGeoJson(tripId: number) {
  return await fetchData("/api/v1/trips/" + tripId + "/transportation/geojson") as GeoJSON.FeatureCollection[]
}
