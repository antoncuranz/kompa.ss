import {Activity} from "@/types.ts";

export default function ActivityEnty({
  activity
}: {
  activity: Activity,
}){

  return (
    <div className="rounded-lg border border-dashed my-2 mx-6 py-2 px-4">
      {activity.name}
      <span className="float-right">
        {activity.time?.substr(0,5)}
      </span>
    </div>
  )
}