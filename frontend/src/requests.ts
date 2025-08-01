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
  const trip = await fetchData("/api/v1/trips/" + tripId) as Trip
  return {
    ...trip,
    startDate: new Date(trip.startDate),
    endDate: new Date(trip.endDate),
  }
}

export async function fetchAccommodation() {
  const accommodation = await fetchData("/api/v1/accommodation") as Accommodation[]
  return accommodation.map(acc => ({
    ...acc,
    arrivalDate: new Date(acc.arrivalDate),
    departureDate: new Date(acc.departureDate)
  }))
}

export async function fetchFlights() {
  const flights = await fetchData("/api/v1/flights") as Flight[]
  return flights.map(flight => ({
    ...flight,
    legs: flight.legs.map(leg => ({
      ...leg,
      departureDateTime: new Date(leg.departureDateTime),
      arrivalDateTime: new Date(leg.arrivalDateTime)
    }))
  }))
}

export async function fetchActivities() {
  const activities = await fetchData("/api/v1/activities") as Activity[];
  return activities.map(activity => ({
    ...activity,
    date: new Date(activity.date)
  }))
}
