"use client"

import {Button} from "@/components/ui/button.tsx";
import {PlaneTakeoff} from "lucide-react";
import {useState} from "react";
import {useRouter} from "next/navigation";
import AddFlightDialog from "@/components/dialog/AddFlightDialog.tsx";
import {Trip} from "@/types.ts";

export default function AddFlightButton({trip}: {trip: Trip}) {
  const [flightDialogOpen, setFlightDialogOpen] = useState(false)
  const router = useRouter()

  function onFlightDialogClose(needsUpdate: boolean) {
    setFlightDialogOpen(false)

    if (needsUpdate)
      router.refresh()
  }

  return (
    <>
      <Button size="sm" className="h-8 gap-1 mt-0 self-end" onClick={() => setFlightDialogOpen(true)}>
        <PlaneTakeoff className="h-3.5 w-3.5"/>
        <span className="sr-only sm:not-sr-only sm:whitespace-nowrap">
          Add Flight
        </span>
      </Button>
      <AddFlightDialog trip={trip} open={flightDialogOpen} onClose={onFlightDialogClose}/>
    </>
  )
}