import {Moment} from "moment";

export default function DayLabel({
  date, location
}: {
  date: Moment,
  location?: string | null
}){

  return (
    <span className="ml-6 mb-2 text-sm text-muted-foreground">
      {date.format("dddd, D MMMM YY")}{location && ", "}{location}
    </span>
  )
}