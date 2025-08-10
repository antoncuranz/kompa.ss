import {Button} from "@/components/ui/button.tsx";
import {CalendarIcon} from "lucide-react";
import {cn} from "@/lib/utils.ts";
import {Popover, PopoverContent, PopoverTrigger} from "@/components/ui/popover.tsx";
import {format} from "date-fns";
import {Calendar} from "@/components/ui/calendar.tsx";
import {Matcher} from "react-day-picker";
import {dateFromString, dateToString, getNextDay} from "@/components/util.ts";

export default function DateInput({
  date, updateDate, startDate, excludeStartDate, endDate, readOnly
}: {
  date: string | null
  updateDate: (newDate: string) => void
  startDate?: string | null
  excludeStartDate?: boolean | null
  endDate?: string | null
  readOnly?: boolean
}) {

  function getMatcher(): Matcher | undefined {
    const adjustedStartDate = getAdjustedStartDate()

    if (adjustedStartDate && endDate) {
      return {before: dateFromString(adjustedStartDate), after: dateFromString(endDate)}
    } else if (adjustedStartDate) {
      return {before: dateFromString(adjustedStartDate)}
    } else if (endDate) {
      return {after: dateFromString(endDate)}
    }
    return undefined
  }

  function getAdjustedStartDate() {
    if (!startDate) {
      return null
    }

    return excludeStartDate ? getNextDay(startDate) : startDate
  }

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button
            variant="secondary"
            className={cn(
                "col-span-3 justify-start text-left font-normal disabled:opacity-1",
                !date && "text-muted-foreground"
            )}
            disabled={readOnly}
        >
          <CalendarIcon className="mr-2 h-4 w-4"/>
          {date ? format(date, "PPP") : <span>Pick a date</span>}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-auto p-0 rounded-2xl overflow-hidden shadow-lg">
        <Calendar
            mode="single"
            selected={date ? dateFromString(date) : undefined}
            onSelect={date => updateDate(dateToString(date))}
            startMonth={startDate ? dateFromString(getAdjustedStartDate()!) : undefined}
            endMonth={endDate ? dateFromString(endDate) : undefined}
            disabled={getMatcher()}
            required={true}
        />
      </PopoverContent>
    </Popover>
  )
}