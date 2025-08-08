"use client"

import React, {useState} from "react";
import {useRouter} from "next/navigation";
import {Trip} from "@/types.ts";
import AddActivityDialog from "@/components/dialog/AddActivityDialog.tsx";
import {DropdownMenu, DropdownMenuItem, DropdownMenuContent, DropdownMenuTrigger} from "@/components/ui/dropdown-menu.tsx";
import {Button} from "@/components/ui/button.tsx";
import {PlaneTakeoff} from "lucide-react";
import AddAccommodationDialog from "@/components/dialog/AddAccommodationDialog.tsx";
import AddFlightDialog from "@/components/dialog/AddFlightDialog.tsx";

export default function AddSomethingDropdown({trip}: {trip: Trip}) {
  const [activityDialogOpen, setActivityDialogOpen] = useState(false)
  const [accommodationDialogOpen, setAccommodationDialogOpen] = useState(false)
  const [flightDialogOpen, setFlightDialogOpen] = useState(false)

  const router = useRouter()

  function onActivityDialogClose(needsUpdate: boolean) {
    console.log("closing activity dialog")
    setActivityDialogOpen(false)
    if (needsUpdate)
      router.refresh()
  }

  function onAccommodationDialogClose(needsUpdate: boolean) {
    setAccommodationDialogOpen(false)
    if (needsUpdate)
      router.refresh()
  }

  function onFlightDialogClose(needsUpdate: boolean) {
    setFlightDialogOpen(false)
    if (needsUpdate)
      router.refresh()
  }

  return (
    <div>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button size="sm" className="h-8 gap-1 mt-0 ml-1 self-end">
            <PlaneTakeoff className="h-3.5 w-3.5"/>
            <span className="sr-only sm:not-sr-only sm:whitespace-nowrap">
              Add Something
            </span>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuItem onClick={() => setActivityDialogOpen(true)}>
            Add Activity
          </DropdownMenuItem>
          <DropdownMenuItem onClick={() => setAccommodationDialogOpen(true)}>
            Add Accommodation
          </DropdownMenuItem>
          <DropdownMenuItem onClick={() => setFlightDialogOpen(true)}>
            Add Flight
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
      <AddActivityDialog trip={trip} open={activityDialogOpen} onClose={onActivityDialogClose}/>
      <AddAccommodationDialog trip={trip} open={accommodationDialogOpen} onClose={onAccommodationDialogClose}/>
      <AddFlightDialog trip={trip} open={flightDialogOpen} onClose={onFlightDialogClose}/>
    </div>
  )
}