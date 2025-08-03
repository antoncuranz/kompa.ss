"use client"

import {Button} from "@/components/ui/button.tsx";
import {Bed} from "lucide-react";
import {useState} from "react";
import {useRouter} from "next/navigation";
import {Trip} from "@/types.ts";
import AddAccommodationDialog from "@/components/dialog/AddAccommodationDialog.tsx";

export default function AddAccommodationButton({trip}: {trip: Trip}) {
  const [flightDialogOpen, setFlightDialogOpen] = useState(false)
  const router = useRouter()

  function onFlightDialogClose(needsUpdate: boolean) {
    setFlightDialogOpen(false)

    if (needsUpdate)
      router.refresh()
  }

  return (
    <>
      <Button size="sm" className="h-8 gap-1 mt-0 ml-1 self-end" onClick={() => setFlightDialogOpen(true)}>
        <Bed className="h-3.5 w-3.5"/>
        <span className="sr-only sm:not-sr-only sm:whitespace-nowrap">
          Add Accommodation
        </span>
      </Button>
      <AddAccommodationDialog trip={trip} open={flightDialogOpen} onClose={onFlightDialogClose}/>
    </>
  )
}