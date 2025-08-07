import {Separator} from "@/components/ui/separator.tsx";
import FlightEntry from "@/components/itinerary/FlightEntry.tsx";
import {Flight, FlightLeg} from "@/types.ts";

export default function FlightSeparator({
  flight, flightLeg
}: {
  flight: Flight,
  flightLeg: FlightLeg
}){

  return (
    <>
      <FlightEntry flight={flight} flightLeg={flightLeg}/>
      <Separator className="relative bottom-5 z-0"/>
    </>
  )
}