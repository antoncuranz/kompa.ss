"use client"

import {Button} from "@/components/ui/button.tsx";
import {Collapsible} from "@/components/ui/collapsible.tsx";
import {ChevronDown, ChevronUp, SquarePen} from "lucide-react";
import {CollapsibleContent, CollapsibleTrigger} from "@radix-ui/react-collapsible";
import {cn} from "@/lib/utils.ts";
import {TrainLeg} from "@/types.ts";
import {formatDurationMinutes, formatTime} from "@/components/util.ts";
import React, {MouseEventHandler, useState} from "react";

export default function TrainEntry({
  trainLeg, className, onInfoBtnClick
}: {
  trainLeg: TrainLeg,
  className?: string
  onInfoBtnClick?: MouseEventHandler<HTMLButtonElement> | undefined
}){
  const [open, setOpen] = useState<boolean>(false)

  function logoFromOperatorName(operatorName: string): string {
    const lowerOperatorName = operatorName.toLowerCase()
    if (lowerOperatorName.startsWith("db")) {
      return "https://assets.static-bahn.de/dam/jcr:47b6ca20-95d9-4102-bc5a-6ebb5634f009/db-logo.svg"
    } else if (lowerOperatorName.startsWith("schweiz")) {
      return "https://digital.sbb.ch/assets/images/brand/signet.svg"
    } else if (lowerOperatorName.startsWith("österreich")) {
      return "https://upload.wikimedia.org/wikipedia/commons/5/5e/Logo_%C3%96BB.svg"
    // } else if (lowerOperatorName.startsWith("trenitalia")) {
    }
    return ""
  }

  return (
    <Collapsible
      open={open}
      onOpenChange={setOpen}
      className={cn("rounded-xl border mx-3 p-2 pl-4 pr-4 grid bg-background z-10 relative", className)}
    >
      <CollapsibleTrigger asChild>
        <div className="grid grid-cols-[1.5rem_1fr] gap-2 cursor-pointer w-full">
          <span className="mt-0 m-auto text-2xl leading-[1.3rem] h-6">🚇</span>
          <div className="flex overflow-hidden whitespace-nowrap w-full">
            <span className="overflow-hidden text-ellipsis w-full">
              {open ?
                `Train from ${trainLeg.origin.name} to ${trainLeg.destination.name}`
              :
                `${formatTime(trainLeg.departureDateTime)}-${formatTime(trainLeg.arrivalDateTime)} ${trainLeg.lineName} from ${trainLeg.origin.name} to ${trainLeg.destination.name}`
              }
            </span>
            {open ?
              <ChevronUp className="float-right text-muted-foreground"/>
            :
              <ChevronDown className="float-right text-muted-foreground"/>
            }
          </div>
        </div>
      </CollapsibleTrigger>
      <CollapsibleContent>
        <div className="grid mt-1" style={{gridTemplateColumns: "1.5rem 1fr", columnGap: "0.5rem"}}>
          <div className="mt-0 m-auto flex flex-col items-center relative top-2">
            <div className="w-1.5 h-1.5 rounded-lg bg-gray-300"/>
            <div className="h-10 w-0.5 bg-gray-300"/>
            <div className="w-1.5 h-1.5 rounded-lg bg-gray-300"/>
          </div>
          <div>
            <p>{formatTime(trainLeg.departureDateTime)} {trainLeg.origin.name}</p>
            <p className="text-sm text-muted-foreground">Duration: {formatDurationMinutes(trainLeg.durationInMinutes)}</p>
            <p>{formatTime(trainLeg.arrivalDateTime)} {trainLeg.destination.name}</p>
          </div>
          <div className="flex items-center">
            <img src={logoFromOperatorName(trainLeg.operatorName)} className="h-auto w-full" alt=""/>
          </div>
          <div>
            <span className="text-sm text-muted-foreground">{trainLeg.operatorName} - {trainLeg.lineName}</span>
            <div className="flex float-right">
              <Button variant="secondary" className="ml-2 p-2 h-6" onClick={onInfoBtnClick}>
                <SquarePen className="w-3.5 h-3.5"/>
              </Button>
            </div>
          </div>
        </div>
      </CollapsibleContent>
    </Collapsible>
 )
}