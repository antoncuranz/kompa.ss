import {formatDateLong} from "@/components/util.ts";

export default function DayLabel({
  date, location
}: {
  date: Date,
  location?: string | null
}){

  return (
    <span className="ml-6 mb-2 text-sm text-muted-foreground">
      {formatDateLong(date)}{location && ", "}{location}
    </span>
  )
}