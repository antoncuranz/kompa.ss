"use client"

import {useState} from "react";
import {useRouter} from "next/navigation";
import {Trip} from "@/types.ts";
import AddActivityDialog from "@/components/dialog/AddActivityDialog.tsx";
import {DropdownMenuItem} from "@radix-ui/react-dropdown-menu";

export default function AddActivityButton({trip}: {trip: Trip}) {
  const [flightDialogOpen, setFlightDialogOpen] = useState(false)
  const router = useRouter()

  function onFlightDialogClose(needsUpdate: boolean) {
    setFlightDialogOpen(false)

    if (needsUpdate)
      router.refresh()
  }

  return (
    <>
      {/*<Button size="sm" className="h-8 gap-1 mt-0 ml-1 self-end" onClick={() => setFlightDialogOpen(true)}>*/}
      {/*  <Flame className="h-3.5 w-3.5"/>*/}
      {/*  <span className="sr-only sm:not-sr-only sm:whitespace-nowrap">*/}
      {/*    Add Activity*/}
      {/*  </span>*/}
      {/*</Button>*/}
      <DropdownMenuItem onClick={() => setFlightDialogOpen(true)}>
        Add Activity
      </DropdownMenuItem>
      <AddActivityDialog trip={trip} open={flightDialogOpen} onClose={onFlightDialogClose}/>
    </>
  )
}