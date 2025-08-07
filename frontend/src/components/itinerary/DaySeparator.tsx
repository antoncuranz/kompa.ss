import {Separator} from "@/components/ui/separator.tsx";

export default function DaySeparator({
  collapsedDays, accomodation
}: {
  collapsedDays?: number | null,
  accomodation?: string | null
}){

  return (
    <>
      <div className="mx-3 mt-2 text-sm text-muted-foreground">
        {accomodation ? "🛏️ " + accomodation : "⚠️ missing accomodation"}
        {collapsedDays ? ` (${collapsedDays} day${collapsedDays != 1 ? "s" : ""} collapsed)` : ""}
      </div>
      <Separator/>
    </>
  )
}