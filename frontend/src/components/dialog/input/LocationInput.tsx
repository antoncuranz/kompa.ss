import {Input} from "@/components/ui/input.tsx";
import {Button} from "@/components/ui/button.tsx";
import {Map} from "lucide-react";
import {cn} from "@/lib/utils.ts";
import {ControllerRenderProps, FieldValues} from "react-hook-form";
import {Dialog} from "@/components/dialog/Dialog.tsx";
import React, {useState} from "react";
import MapDialogContent from "@/components/dialog/MapDialogContent.tsx";
import {Coordinates} from "@/types.ts";

export default function LocationInput({
  onChange, onBlur, value, disabled,
  className
}: ControllerRenderProps<FieldValues, string> & {
  className?: string,
}) {
  const [mapDialogOpen, setMapDialogOpen] = useState(false)

  function onMapDialogClose(newLocation: Coordinates) {
    onChange(newLocation)
    onBlur()
  }

  return (
      <div className={cn("", className)}>
        <div className="flex gap-2">
          <Input value={value?.longitude ?? ""}
                 disabled={true}
                 placeholder="Longitude"
          />
          <Input value={value?.latitude ?? ""}
                 disabled={true}
                 placeholder="Latitude"
          />
          {!disabled &&
            <Button variant="secondary" onClick={() => setMapDialogOpen(true)}>
              <Map className="h-4 w-4"/>
            </Button>
          }
        </div>
        <Dialog open={mapDialogOpen} setOpen={setMapDialogOpen}>
          <MapDialogContent value={value} onChange={onMapDialogClose}/>
        </Dialog>
      </div>
  )
}