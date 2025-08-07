import {Activity} from "@/types.ts";
import {formatTime} from "@/components/util.ts";

export default function ActivityEnty({
  activity
}: {
  activity: Activity,
}){

  return (
    <div className="rounded-lg border border-dashed my-4 mx-3 py-2 px-4">
      {activity.name}
      <span className="float-right">
        {activity.time && formatTime(activity.time)}
      </span>
    </div>
  )
}