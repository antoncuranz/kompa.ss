import {getTransportationTypeEmoji, Transportation} from "@/types.ts";
import React, {MouseEventHandler} from "react";
import {formatTime, titleCase} from "../util";

export default function TransportationEntry({
  transportation, onClick
}: {
  transportation: Transportation,
  onClick?: MouseEventHandler<HTMLDivElement> | undefined
}){

  return (
    <div className="cursor-pointer rounded-xl border mx-3 p-2 pl-4 pr-4 grid bg-background z-10 relative" onClick={onClick}>
      <div className="grid grid-cols-[1.5rem_1fr] gap-2">
        <span className="mt-0 m-auto text-2xl leading-[1.3rem] h-6">{getTransportationTypeEmoji(transportation.type)}</span>
        <div className="flex overflow-hidden whitespace-nowrap w-full">
          <span className="overflow-hidden text-ellipsis w-full">
            {formatTime(transportation.departureDateTime)}-{formatTime(transportation.arrivalDateTime)} {titleCase(transportation.type)}
          </span>
        </div>
      </div>
    </div>
 )
}