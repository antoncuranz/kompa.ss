
export type OmitNever<T extends Record<string, unknown>> = {
  [K in keyof T as T[K] extends never ? never : K]: T[K];
};
export type SharedProperties<A, B> = OmitNever<Pick<A & B, keyof A & keyof B>>;

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

export type Coordinates = {
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
  Plane = "PLANE", // deprecated
  Flight = "FLIGHT",
  Train = "TRAIN",
  Bus   = "BUS",
  Ferry  = "FERRY",
  Boat  = "BOAT",
  Bike  = "BIKE",
  Car   = "CAR",
  Hike  = "HIKE",
  Other = "OTHER",
}

export function getTransportationTypeEmoji(type: TransportationType) {
  switch (type) {
    case TransportationType.Plane:
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
  genericDetail: GenericDetail | null;
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
  operatorName: string;
};

export type TrainStation = {
  id: string;
  name: string;
  location: Location;
}

export type GenericDetail = {
  name: string;
  originAddress: string|null;
  destinationAddress: string|null;
};

export type DayRenderData = {
  day: string;
  transportation: Transportation[];
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
