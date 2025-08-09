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
  // await new Promise( resolve => setTimeout(resolve, 5 * 1000) )
  const trip = await fetchData("/api/v1/trips/" + tripId) as Trip
  return {
    ...trip,
    startDate: new Date(trip.startDate),
    endDate: new Date(trip.endDate),
  }
}

export async function fetchTrips() {
  const trips = await fetchData("/api/v1/trips") as Trip[]
  return trips.map(trip => ({
    ...trip,
    startDate: new Date(trip.startDate),
    endDate: new Date(trip.endDate),
  }))
}

export async function fetchAccommodation(tripId: number) {
  const accommodation = await fetchData("/api/v1/trips/" + tripId + "/accommodation") as Accommodation[]
  return accommodation.map(acc => ({
    ...acc,
    arrivalDate: new Date(acc.arrivalDate),
    departureDate: new Date(acc.departureDate)
  }))
}

export async function fetchFlights(tripId: number) {
  const flights = await fetchData("/api/v1/trips/" + tripId + "/flights") as Flight[]
  return flights.map(flight => ({
    ...flight,
    legs: flight.legs.map(leg => ({
      ...leg,
      departureDateTime: new Date(leg.departureDateTime),
      arrivalDateTime: new Date(leg.arrivalDateTime)
    }))
  }))
}

export async function fetchActivities(tripId: number) {
  const activities = await fetchData("/api/v1/trips/" + tripId + "/activities") as Activity[];
  return activities.map(activity => ({
    ...activity,
    date: new Date(activity.date),
    time: activity.time ? new Date("1970-01-01T" + activity.time) : null
  }))
}
