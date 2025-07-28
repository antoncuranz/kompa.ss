
export default function ActivityEnty({
  name
}: {
  name: string,
}){

  return (
    <div className="rounded-lg border border-dashed my-2 mx-6 py-2 px-4">
      {name}
    </div>
  )
}