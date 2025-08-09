import {Button} from "@/components/ui/button.tsx";
import {CalendarIcon} from "lucide-react";
import {cn} from "@/lib/utils.ts";
import {Popover, PopoverContent, PopoverTrigger} from "@/components/ui/popover.tsx";
import {format} from "date-fns";
import {Calendar} from "@/components/ui/calendar.tsx";

export default function DateInput({
  date, updateDate, startDate, endDate, readOnly
}: {
  date: Date | null
  updateDate: (newDate: Date) => void
  startDate: Date
  endDate: Date
  readOnly?: boolean
}) {

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
            selected={date ?? undefined}
            onSelect={updateDate}
            startMonth={startDate}
            endMonth={endDate}
            disabled={{before: startDate, after: endDate}}
            required={true}
        />
      </PopoverContent>
    </Popover>
  )
}