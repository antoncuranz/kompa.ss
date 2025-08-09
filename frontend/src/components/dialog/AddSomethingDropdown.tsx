"use client"

import React, {useState} from "react";
import {useRouter} from "next/navigation";
import {Trip} from "@/types.ts";
import ActivityDialog from "@/components/dialog/ActivityDialog.tsx";
import {DropdownMenu, DropdownMenuItem, DropdownMenuContent, DropdownMenuTrigger} from "@/components/ui/dropdown-menu.tsx";
import {Button} from "@/components/ui/button.tsx";
import {PlaneTakeoff} from "lucide-react";
import AccommodationDialog from "@/components/dialog/AccommodationDialog.tsx";
import FlightDialog from "@/components/dialog/FlightDialog.tsx";

export default function AddSomethingDropdown({trip}: {trip: Trip}) {
  const [activityDialogOpen, setActivityDialogOpen] = useState(false)
  const [accommodationDialogOpen, setAccommodationDialogOpen] = useState(false)
  const [flightDialogOpen, setFlightDialogOpen] = useState(false)

  const router = useRouter()

  function onActivityDialogClose(needsUpdate: boolean) {
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
      <ActivityDialog key={"act-" + activityDialogOpen} trip={trip} open={activityDialogOpen} onClose={onActivityDialogClose}/>
      <AccommodationDialog key={"acc-" + accommodationDialogOpen} trip={trip} open={accommodationDialogOpen} onClose={onAccommodationDialogClose}/>
      <FlightDialog key={"flight-" + flightDialogOpen} trip={trip} open={flightDialogOpen} onClose={onFlightDialogClose}/>
    </div>
  )
}