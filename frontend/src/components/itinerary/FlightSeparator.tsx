import {Separator} from "@/components/ui/separator.tsx";
import FlightEntry from "@/components/itinerary/FlightEntry.tsx";
import {Flight, FlightLeg} from "@/types.ts";
import DayLabel from "@/components/itinerary/DayLabel.tsx";

export default function FlightSeparator({
  date, flight, flightLeg, location
}: {
  date: Date,
  flight: Flight,
  flightLeg: FlightLeg,
  location?: string | null
}){

  return (
    <>
      <FlightEntry flight={flight} flightLeg={flightLeg}/>
      <Separator className="relative bottom-5 z-0"/>
      <DayLabel date={date} location={location}/>
    </>
  )
}