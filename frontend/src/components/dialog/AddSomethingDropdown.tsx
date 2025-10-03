"use client"

import React, {useState} from "react";
import {Trip} from "@/types.ts";
import {DropdownMenu, DropdownMenuItem, DropdownMenuContent, DropdownMenuTrigger} from "@/components/ui/dropdown-menu.tsx";
import {Button} from "@/components/ui/button.tsx";
import {PlaneTakeoff} from "lucide-react";
import AccommodationDialogContent from "@/components/dialog/AccommodationDialogContent.tsx";
import FlightDialogContent from "@/components/dialog/FlightDialogContent.tsx";
import TrainDialogContent from "@/components/dialog/TrainDialogContent.tsx";
import ActivityDialogContent from "@/components/dialog/ActivityDialogContent.tsx";
import {Dialog} from "@/components/dialog/Dialog.tsx";

export default function AddSomethingDropdown({trip}: {trip: Trip}) {
  const [activityDialogOpen, setActivityDialogOpen] = useState(false)
  const [accommodationDialogOpen, setAccommodationDialogOpen] = useState(false)
  const [flightDialogOpen, setFlightDialogOpen] = useState(false)
  const [trainDialogOpen, setTrainDialogOpen] = useState(false)

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
          <DropdownMenuItem onClick={() => setTrainDialogOpen(true)}>
            Add Train
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
      <Dialog open={activityDialogOpen} setOpen={setActivityDialogOpen}>
        <ActivityDialogContent trip={trip}/>
      </Dialog>
      <Dialog open={accommodationDialogOpen} setOpen={setAccommodationDialogOpen}>
        <AccommodationDialogContent trip={trip}/>
      </Dialog>
      <Dialog open={flightDialogOpen} setOpen={setFlightDialogOpen}>
        <FlightDialogContent trip={trip}/>
      </Dialog>
      <Dialog open={trainDialogOpen} setOpen={setTrainDialogOpen}>
        <TrainDialogContent trip={trip}/>
      </Dialog>
    </div>
  )
}