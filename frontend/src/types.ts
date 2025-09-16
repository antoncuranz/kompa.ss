
export type Trip = {
  id: number;
  name: string;
  startDate: string;
  endDate: string;
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
  date: string;
  time: string | null;
  description: string;
  address: string | null;
  location: Location | null;
  price: number | null;
};

export type Accommodation = {
  id: number;
  tripId: number;
  name: string;
  arrivalDate: string;
  departureDate: string;
  checkInTime: string | null;
  checkOutTime: string | null;
  description: string;
  address: string | null;
  location: Location | null;
  price: number | null;
};

export enum TransportationType {
  Plane = "PLANE",
  Train = "TRAIN",
  Bus   = "BUS",
  Boat  = "BOAT",
  Bike  = "BIKE",
  Car   = "CAR",
  Foot  = "FOOT",
  Other = "OTHER",
}

export type Transportation = {
  id: number;
  tripId: number;
  type: TransportationType;
  origin: Location;
  destination: Location;
  departureDateTime: string;
  arrivalDateTime: string;
  geoJson: string | null;
  price: number | null;
  flightDetail: FlightDetail | null;
  trainDetail: TrainDetail | null;
};

export type FlightDetail = {
  legs: FlightLeg[];
  pnrs: PNR[];
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
  departureDateTime: string;
  arrivalDateTime: string;
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
  location: Location;
}

export type TrainDetail = {
  legs: TrainLeg[];
};

export type TrainLeg = {
  id: number;
  origin: TrainStation;
  destination: TrainStation;
  departureDateTime: string;
  arrivalDateTime: string;
  durationInMinutes: number;
  lineName: string;
};

export type TrainStation = {
  id: string;
  name: string;
}

export type DayRenderData = {
  day: string;
  transportation: Transportation[];
  activities: Activity[];
  accommodation: Accommodation | undefined;
};
