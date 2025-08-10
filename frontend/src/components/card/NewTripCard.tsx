"use client"

import React, {useState} from "react";
import Card from "@/components/card/Card.tsx";
import {CirclePlus} from "lucide-react";
import TripDialog from "@/components/dialog/TripDialog.tsx";
import {useRouter} from "next/navigation";

export default function NewTripCard() {
  const [tripDialogOpen, setTripDialogOpen] = useState(false)

  const router = useRouter()

  function onTripDialogClose(needsUpdate: boolean) {
    setTripDialogOpen(false)
    if (needsUpdate)
      router.refresh()
  }

  return (
    <>
      <Card key="new-trip" className="h-80 w-56 md:h-[40rem] md:w-96 shadow-none hover:shadow-xl hover:cursor-pointer" onClick={() => setTripDialogOpen(true)}>
        <div className="h-full w-full rounded-2xl no-scrollbar overflow-hidden overflow-y-scroll flex items-center justify-center">
          <CirclePlus className="w-14 h-14 text-gray-400"/>
        </div>
      </Card>
      <TripDialog open={tripDialogOpen} onClose={onTripDialogClose}/>
    </>
  )
}