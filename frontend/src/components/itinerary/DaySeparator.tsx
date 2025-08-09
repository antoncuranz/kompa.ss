import {Separator} from "@/components/ui/separator.tsx";
import {Accommodation} from "@/types.ts";

export default function DaySeparator({
  collapsedDays, accomodation, onAccommodationClick = () => {}
}: {
  collapsedDays: number,
  accomodation: Accommodation | undefined
  onAccommodationClick?: (accommodation: Accommodation | undefined) => void;
}){

  return (
    <>
      <div className="mx-3 mt-2 text-sm text-muted-foreground">
        <span className="hover:underline hover:cursor-pointer" onClick={() => onAccommodationClick(accomodation)}>
          {accomodation ? `🛏️ ${accomodation.name}` : "⚠️ missing accomodation"}
        </span>
        {collapsedDays > 0 &&
          <span> ({collapsedDays} {collapsedDays != 1 ? "days" : "day"} collapsed)</span>
        }
      </div>
      <Separator/>
    </>
  )
}