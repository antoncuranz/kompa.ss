import {
  Accommodation,
  Activity,
  Flight, Trip,
} from "@/types.ts";
import React from "react";

class UnauthorizedError extends Error {}

export async function hideIfUnauthorized(func: () => Promise<React.ReactElement>) {
  try {
    return await func()
  } catch (e: unknown) {
    if (e instanceof UnauthorizedError)
      return ""
    else
      throw e
  }
}

const usernameHeader = "X-Auth-Request-Preferred-Username"

async function fetchData(url: string) {
  const response = await fetch(process.env.BACKEND_URL + url, {
    headers: {[usernameHeader]: await getCurrentUser()},
    cache: "no-cache"
  })
  if (response.status == 401) {
    throw new UnauthorizedError()
  } else if (!response.ok) {
    throw new Error("Failed to fetch data");
  }
  return await response.json()
}

export async function getCurrentUser(): Promise<string> {
  return "ant0n"
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

export async function fetchFlights(tripId: number) {
  return await fetchData("/api/v1/trips/" + tripId + "/flights") as Flight[]
}

export async function fetchActivities(tripId: number) {
  return await fetchData("/api/v1/trips/" + tripId + "/activities") as Activity[];
}
