import {Input} from "@/components/ui/input.tsx";
import {useToast} from "@/components/ui/use-toast.ts";
import {Button} from "@/components/ui/button.tsx";
import {Pencil, Search} from "lucide-react";
import {cn} from "@/lib/utils.ts";
import {useState} from "react";
import {TrainStation} from "@/types.ts";

export default function TrainStationInput({
  station, updateStation, placeholder, className, readOnly
}: {
  station: TrainStation | null,
  updateStation: (newStation: TrainStation|null) => void,
  placeholder?: string,
  className?: string,
  readOnly?: boolean,
}) {
  const { toast } = useToast();
  const [edit, setEdit] = useState<boolean>(station == null)
  const [text, setText] = useState<string>(station?.name ?? "")

  async function searchForStationUsingApi() {
    const query = text.replace(" ", "-")
    const url = encodeURI("/api/v1/trips/1/trains/stations?query=" + query)
    const response = await fetch(url)

    if (response.ok) {
      const json = await response.json()
      updateStation(json as TrainStation)
      setEdit(false)
    } else toast({
      title: "No stations found",
      description: await response.text()
    })
  }

  async function onButtonClick() {
    if (edit) {
      await searchForStationUsingApi()
    } else {
      setText(station?.name ?? "")
      setEdit(true)
      updateStation(null)
    }
  }

  return (
      <div className={cn("", className)}>
        <div className="flex gap-2">
          <Input value={edit ? text : station?.name ?? ""}
                 onChange={e => edit && setText(e.target.value)}
                 placeholder={placeholder}
                 readOnly={readOnly}/>
          {!readOnly &&
            <Button variant="secondary" onClick={onButtonClick}>
              {edit ?
                <Search className="h-4 w-4"/>
              :
                <Pencil className="h-4 w-4"/>
              }
            </Button>
          }
        </div>
      </div>
  )
}