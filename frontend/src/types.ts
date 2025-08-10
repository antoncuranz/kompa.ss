
export type Trip = {
  id: number;
  name: string;
  startDate: Date;
  endDate: Date;
  description: string;
  imageUrl: string;
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
  time: Date | null;
  description: string;
  address: string | null;
  location: Location | null;
  price: number | null;
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
  address: string | null;
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

export type AddPNR = {
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

export type AddFlightLeg = {
  date: string;
  flightNumber: string;
  originAirport: string | null;
};

export type Airport = {
  iata: string;
  name: string;
  municipality: string;
  location: Location | null;
}

export type DayRenderData = {
  day: Date;
  flights: {flight: Flight; leg: FlightLeg}[];
  activities: Activity[];
  accommodation: Accommodation | undefined;
};
