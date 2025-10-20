import Link from "next/link";
import {ModeToggle} from "@/components/buttons/ModeToggle.tsx";
import {Button} from "@/components/ui/button.tsx";
import * as React from "react";
import {GalleryHorizontalEnd} from "lucide-react";
import {ButtonGroup} from "../ui/button-group";
import {Trip} from "@/types.ts";

export default function Navigation({
  trip
}: {
  trip?: Trip | null
}) {
  return (
    <header className="sticky top-0 flex h-14 items-center gap-4 pr-4 md:pr-6">
      <nav className="font-medium flex flex-row items-center gap-5 text-sm lg:gap-6 w-full">
        <div className="flex gap-4 lg:gap-6 overflow-x-auto w-full no-scrollbar h-10 items-center pl-4 md:pl-6 pr-10"
             style={{maskImage: "linear-gradient(to right, transparent .0em, black 1em calc(100% - 3em), transparent calc(100% - .0em))"}}
        >
          { trip ?
            <>
              <ButtonGroup>
                <Link href="/">
                  <Button size="sm" variant="outline" className="rounded-l-full">
                    <GalleryHorizontalEnd/>
                  </Button>
                </Link>
                <Button size="sm" variant="outline" className="rounded-r-full pointer-events-none">{trip.name}</Button>
              </ButtonGroup>
              <Link href={"/" + trip.id + "/itinerary"}>
                  <span className="text-muted-foreground transition-colors hover:text-foreground">Itinerary</span>
              </Link>
              <Link href={"/" + trip.id + "/map"}>
                  <span className="text-muted-foreground transition-colors hover:text-foreground">Map</span>
              </Link>
            </>
          :
            <Link href="/">
              <strong>🧭 kompa.ss</strong>
            </Link>
          }
        </div>
        <ModeToggle/>
      </nav>
    </header>
  )
}