import {Activity} from "@/types.ts";
import {formatTime} from "@/components/util.ts";
import {MouseEventHandler} from "react";

export default function ActivityEnty({
  activity, onClick
}: {
  activity: Activity
  onClick?: MouseEventHandler<HTMLDivElement> | undefined
}){

  return (
    <div className="rounded-lg border border-dashed my-4 mx-3 py-2 px-4 hover:border-solid hover:cursor-pointer" onClick={onClick}>
      {activity.name}
      <span className="float-right">
        {activity.time && formatTime(activity.time, true)}
      </span>
    </div>
  )
}