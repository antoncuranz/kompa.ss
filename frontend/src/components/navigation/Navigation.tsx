import Link from "next/link";
import {ModeToggle} from "@/components/buttons/ModeToggle.tsx";

export default function Navigation({
  tripId
}: {
  tripId?: number | null
}) {
  return (
    <header className="sticky top-0 flex h-14 sm:h-16 items-center gap-4 pr-4 md:pr-6">
      <nav className="font-medium flex flex-row items-center gap-5 text-sm lg:gap-6 w-full">
        <div className="flex gap-4 lg:gap-6 overflow-x-auto w-full no-scrollbar h-10 items-center pl-4 md:pl-6 pr-10"
             style={{maskImage: "linear-gradient(to right, transparent .0em, black 1em calc(100% - 3em), transparent calc(100% - .0em))"}}
        >
          { tripId &&
            <>
              <Link href="/">
                <span className="text-muted-foreground transition-colors hover:text-foreground">Trips</span>
              </Link>
              <Link href={"/" + tripId + "/itinerary"}>
                  <span className="text-muted-foreground transition-colors hover:text-foreground">Itinerary</span>
              </Link>
              <Link href={"/" + tripId + "/cost"}>
                  <span className="text-muted-foreground transition-colors hover:text-foreground">Cost</span>
              </Link>
              <Link href={"/" + tripId + "/map"}>
                  <span className="text-muted-foreground transition-colors hover:text-foreground">Map</span>
              </Link>
            </>
          }
        </div>
        <ModeToggle/>
      </nav>
    </header>
  )
}