import {Separator} from "@/components/ui/separator.tsx";
import {Moment} from "moment";
import DayLabel from "@/components/itinerary/DayLabel.tsx";

export default function DaySeparator({
  date, accomodation, location
}: {
  date: Moment,
  accomodation?: string | null,
  location?: string | null
}){

  return (
    <>
      <div className="ml-6 mt-2 text-sm text-muted-foreground">
        {accomodation ? "🛏️ " + accomodation : "⚠️ missing accomodation"}
      </div>
      <Separator/>
      <DayLabel date={date} location={location}/>
    </>
  )
}