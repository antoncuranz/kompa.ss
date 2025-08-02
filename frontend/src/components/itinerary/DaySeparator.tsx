import {Separator} from "@/components/ui/separator.tsx";
import DayLabel from "@/components/itinerary/DayLabel.tsx";

export default function DaySeparator({
  date, collapsedDays, accomodation, location
}: {
  date: Date,
  collapsedDays?: number | null
  accomodation?: string | null,
  location?: string | null
}){

  return (
    <>
      <div className="ml-6 mt-2 text-sm text-muted-foreground">
        {accomodation ? "🛏️ " + accomodation : "⚠️ missing accomodation"}
        {collapsedDays && ` (${collapsedDays} day${collapsedDays != 1 ? "s" : ""} collapsed)`}
      </div>
      <Separator/>
      <DayLabel date={date} location={location}/>
    </>
  )
}