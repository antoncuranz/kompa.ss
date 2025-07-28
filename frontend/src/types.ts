import {Moment} from "moment";

export type Trip = {
  id: number;
  name: string;
  description: string;
  startDate: Moment;
  endDate: Moment;
};

export type Activity = {
  id: number;
  name: string;
  description: string;
  date: Moment;
  // startTime: Date | null;
  // endTime: Date | null;
  // price: number | null;
  // location: string | null;
};

export type Accommodation = {
  id: number;
  name: string;
  description: string;
  arrivalDate: Moment;
  departureDate: Moment;
  price: number | null;
  location: string | null;
};

export type Flight = {
  id: number;
  type: string;
  // origin: string; // should be locations for Mapbox!
  // destination: string;
  price: number | null;
  geojson: string | null;
  pnrs: PNR[];
  legs: FlightLeg[];
};

export type PNR = {
  id: number;
  airline: string;
  pnr: string;
};

export type FlightLeg = {
  id: number;
  airline: string;
  flightNumber: string;
  origin: Airport;
  destination: Airport;
  departureTime: Moment;
  arrivalTime: Moment;
  aircraft: string | null;
};

export type Airport = {
  name: string;
  iata: string;
  municipality: string;
}
