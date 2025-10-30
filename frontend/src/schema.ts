import * as core from "zod/v4/core";
import {dateTimeToString, dateToString} from "@/components/util.ts";
import {co, z} from "jazz-tools";

export type OmitNever<T extends Record<string, unknown>> = {
  [K in keyof T as T[K] extends never ? never : K]: T[K];
};
export type SharedProperties<A, B> = OmitNever<Pick<A & B, keyof A & keyof B>>;

export function optionalString(params?: string | core.$ZodStringParams) {
  return z.string(params).trim().transform(x => x || undefined).optional()
}

export function isoDate(params?: string | core.$ZodDateParams) {
  return z.date(params).transform(dateToString)
}

export function isoDateTime(params?: string | core.$ZodDateParams) {
  return z.date(params).transform(dateTimeToString)
}

export function optionalNumberString(params?: string | core.$ZodStringParams) {
  return z.string(params).transform(x => x ? Number(x) : undefined)
}

export function location(params?: string | core.$ZodObjectParams) {
  return z.object({
    latitude: z.number(),
    longitude: z.number()
  }, params)
}

export function optionalLocation(params?: string | core.$ZodObjectParams) {
  return z.object({
    latitude: z.number(),
    longitude: z.number()
  }, params).transform(x => (x.latitude && x.longitude) ? x : undefined).optional()
}

export function trainStation(params?: string | core.$ZodObjectParams) {
  return z.object({
    id: z.string().nonempty(),
    name: z.string().nonempty(),
  }, params).transform(station => station.id)
}

// Jazz types

export const Location = co.map({
  latitude: z.number(),
  longitude: z.number()
})
export type Location = co.loaded<typeof Location>;

export const Activity = co.map({
  name: z.string(),
  description: z.string().optional(),
  date: z.iso.date(),
  // time: z.iso.time(),
  price: z.number().optional(),
  address: z.string().optional(),
  location: Location.optional()
})
export const RESOLVE_ACTIVITY = {
  location: true
}
export type Activity = co.loaded<typeof Activity, typeof RESOLVE_ACTIVITY>;

export const Accommodation = co.map({
  name: z.string(),
  description: z.string().optional(),
  arrivalDate: z.iso.date(),
  departureDate: z.iso.date(),
  price: z.number().optional(),
  address: z.string().optional(),
  location: Location.optional()
})
export const RESOLVE_ACCOMMODATION = {
  location: true
}
export type Accommodation = co.loaded<typeof Accommodation, typeof RESOLVE_ACCOMMODATION>;

export const Airport = co.map({
  iata: z.string(),
  name: z.string(),
  municipality: z.string(),
  location: Location
})
export const RESOLVE_AIRPORT = {
  location: true
}
export type Airport = co.loaded<typeof Airport, typeof RESOLVE_AIRPORT>;

export const FlightLeg = co.map({
  origin: Airport,
  destination: Airport,
  airline: z.string(),
  flightNumber: z.string(),
  departureDateTime: z.iso.datetime(),
  arrivalDateTime: z.iso.datetime(),
  durationInMinutes: z.number(),
  aircraft: z.string().optional()
})
export const RESOLVE_FLIGHT_LEG = {
  origin: RESOLVE_AIRPORT,
  destination: RESOLVE_AIRPORT
}
export type FlightLeg = co.loaded<typeof FlightLeg, typeof RESOLVE_FLIGHT_LEG>;

export const PNR = co.map({
  airline: z.string(),
  pnr: z.string()
})
export type PNR = co.loaded<typeof PNR>;

export const FlightDetail = co.map({
  legs: co.list(FlightLeg),
  pnrs: co.list(PNR)
})
export const RESOLVE_FLIGHT_DETAIL = {
  legs: {$each: RESOLVE_FLIGHT_LEG},
  pnrs: {$each: true}
}
export type FlightDetail = co.loaded<typeof FlightDetail, typeof RESOLVE_FLIGHT_DETAIL>;

export const TrainStation = co.map({
  name: z.string(),
  location: Location
})
export const RESOLVE_TRAIN_STATION = {
  location: true
}
export type TrainStation = co.loaded<typeof TrainStation, typeof RESOLVE_TRAIN_STATION>;

export const TrainLeg = co.map({
  origin: TrainStation,
  destination: TrainStation,
  departureDateTime: z.iso.datetime(),
  arrivalDateTime: z.iso.datetime(),
  durationInMinutes: z.number(),
  lineName: z.string(),
  operatorName: z.string()
})
export const RESOLVE_TRAIN_LEG = {
  origin: RESOLVE_TRAIN_STATION,
  destination: RESOLVE_TRAIN_STATION
}
export type TrainLeg = co.loaded<typeof TrainLeg, typeof RESOLVE_TRAIN_LEG>;

export const GenericTransportation = co.map({
  type: z.literal("generic"),
  name: z.string(),
  genericType: z.string(),
  price: z.number().optional(),
  departureDateTime: z.iso.datetime(),
  arrivalDateTime: z.iso.datetime(),
  origin: Location,
  destination: Location,
  originAddress: z.string().optional(),
  destinationAddress: z.string().optional(),
})
export const RESOLVE_GENERIC_TRANSPORTATION = {
  origin: true,
  destination: true
}
export type GenericTransportation = co.loaded<typeof GenericTransportation, typeof RESOLVE_GENERIC_TRANSPORTATION>;

export const Flight = co.map({
  type: z.literal("flight"),
  legs: co.list(FlightLeg),
  pnrs: co.list(PNR),
  price: z.number().optional()
})
export const RESOLVE_FLIGHT = {
  legs: {$each: RESOLVE_FLIGHT_LEG},
  pnrs: {$each: true}
}
export type Flight = co.loaded<typeof Flight, typeof RESOLVE_FLIGHT>;

export const Train = co.map({
  type: z.literal("train"),
  legs: co.list(TrainLeg),
  price: z.number().optional()
})
export const RESOLVE_TRAIN = {
  legs: {$each: RESOLVE_TRAIN_LEG},
}
export type Train = co.loaded<typeof Train, typeof RESOLVE_TRAIN>;

export const Transportation = co.discriminatedUnion("type", [GenericTransportation, Flight, Train]);
export type Transportation = co.loaded<typeof Transportation>;
export type LoadedTransportation = Flight|Train|GenericTransportation;

export const Trip = co.map({
  name: z.string(),
  startDate: z.iso.date(),
  endDate: z.iso.date(),
  description: z.string().optional(),
  imageUrl: z.string().optional(),
  activities: co.list(Activity),
  accommodation: co.list(Accommodation),
  transportation: co.list(Transportation),
})
export const RESOLVE_TRIP = {
  activities: {$each: RESOLVE_ACTIVITY},
  accommodation: {$each: RESOLVE_ACCOMMODATION},
  transportation: {$each: true},
}
export type Trip = co.loaded<typeof Trip, typeof RESOLVE_TRIP>;

export const AccountRoot = co.map({
  trips: co.list(Trip),
});
export const RESOLVE_ROOT = {
  trips: {$each: RESOLVE_TRIP}
}
export type AccountRoot = co.loaded<typeof AccountRoot, typeof RESOLVE_ROOT>;

export const JazzAccount = co.account({
  root: AccountRoot,
  profile: co.profile(),
}).withMigration(async (account) => {
  if (!account.$jazz.has("root")) {
    account.$jazz.set("root", {
      trips: [],
    });
  }
});
export const RESOLVE_ACCOUNT = {
  profile: true,
  root: RESOLVE_ROOT
}
export type JazzAccount = co.loaded<typeof JazzAccount, typeof RESOLVE_ACCOUNT>;


// Regular types

export enum TransportationType {
  Flight = "FLIGHT",
  Train = "TRAIN",
  Bus   = "BUS",
  Ferry = "FERRY",
  Boat  = "BOAT",
  Bike  = "BIKE",
  Car   = "CAR",
  Hike  = "HIKE",
  Other = "OTHER",
}

export function getTransportationTypeEmoji(type: string) {
  switch (type) {
    case TransportationType.Flight:
      return "✈️"
    case TransportationType.Train:
      return "🚇"
    case TransportationType.Bus:
      return "🚌"
    case TransportationType.Car:
      return "🚗"
    case TransportationType.Ferry:
      return "⛴️"
    case TransportationType.Boat:
      return "⛵️"
    case TransportationType.Bike:
      return "🚲"
    case TransportationType.Hike:
      return "🥾"
    case TransportationType.Other:
    default:
      return "🛸"
  }
}

export type AmbiguousFlightChoice = {
  departureDateTime: string;
  destinationIata: string;
  originIata: string;
}

export type DayRenderData = {
  day: string;
  transportation: LoadedTransportation[];
  activities: Activity[];
  accommodation: Accommodation | undefined;
};

export type GeoJsonFlight = {
  type: TransportationType;
  fromMunicipality: string;
  toMunicipality: string;
  legs: string;
}

export type GeoJsonFlightLeg = {
  flightNumber: string;
  departureDateTime: string;
  arrivalDateTime: string;
  fromIata: string;
  toIata: string;
}

export type GeoJsonTrain = {
  type: TransportationType;
  fromMunicipality: string;
  toMunicipality: string;
  legs: string;
}

export type GeoJsonTrainLeg = {
  lineName: string;
  departureDateTime: string;
  arrivalDateTime: string;
  fromStation: string;
  toStation: string;
}

export type GeoJsonTransportation = {
  type: TransportationType;
  name: string;
  departureDateTime: string;
  arrivalDateTime: string;
}

export type Coordinates = {
  latitude: number;
  longitude: number;
}
