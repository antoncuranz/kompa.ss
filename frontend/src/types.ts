
export type Trip = {
  id: number;
  name: string;
  startDate: Date;
  endDate: Date;
  description: string;
};

export type Location = {
  id: number;
  latitude: number;
  longitude: number;
}

export type Activity = {
  id: number;
  tripId: number;
  name: string;
  date: Date;
  time: string | null;
  description: string;
  // location: Location | null;
  // price: number | null;
};

export type Accommodation = {
  id: number;
  tripId: number;
  name: string;
  arrivalDate: Date;
  departureDate: Date;
  checkInTime: string | null;
  checkOutTime: string | null;
  description: string;
  location: Location | null;
  price: number | null;
};

export type Flight = {
  id: number;
  tripId: number;
  legs: FlightLeg[];
  pnrs: PNR[];
  price: number | null;
};

export type PNR = {
  id: number;
  airline: string;
  pnr: string;
};

export type FlightLeg = {
  id: number;
  origin: Airport;
  destination: Airport;
  airline: string;
  flightNumber: string;
  departureDateTime: Date;
  arrivalDateTime: Date;
  durationInMinutes: number;
  aircraft: string | null;
};

export type Airport = {
  iata: string;
  name: string;
  municipality: string;
  location: Location | null;
}
