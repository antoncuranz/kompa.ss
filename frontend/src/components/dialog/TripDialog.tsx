import {Button} from "@/components/ui/button.tsx";
import {Dialog, DialogFooter, DialogContent, DialogHeader, DialogTitle} from "@/components/ui/dialog.tsx";
import {Label} from "@/components/ui/label.tsx";
import {useState} from "react";
import {useToast} from "@/components/ui/use-toast.ts";
import {Input} from "@/components/ui/input.tsx";
import {Trip} from "@/types.ts";
import {Textarea} from "@/components/ui/textarea.tsx";
import {nullIfEmpty} from "@/components/util.ts";
import {LabelInputContainer, RowContainer} from "@/components/dialog/DialogUtil.tsx";
import DateInput from "@/components/dialog/DateInput.tsx";

export default function TripDialog({
  trip, open, onClose
}: {
  open: boolean,
  onClose: (needsUpdate: boolean) => void,
  trip?: Trip | null
}) {
  const [edit, setEdit] = useState<boolean>(trip == null)

  const [name, setName] = useState<string>(trip?.name ?? "")
  const [startDate, setStartDate] = useState<string|null>(trip?.startDate ?? null)
  const [endDate, setEndDate] = useState<string|null>(trip?.endDate ?? null)
  const [description, setDescription] = useState<string>(trip?.description ?? "")
  const [imageUrl, setImageUrl] = useState<string>(trip?.imageUrl ?? "")

  const { toast } = useToast();

  async function onSaveButtonClick() {
    const body = JSON.stringify({
      name: name,
      startDate: startDate ?? null,
      endDate: endDate ?? null,
      description: nullIfEmpty(description),
      imageUrl: nullIfEmpty(imageUrl)
    })

    let response
    if (trip != null) {
      response = await fetch("/api/v1/trips/" + trip.id, {
        method: "PUT",
        headers: {"Content-Type": "application/json"},
        body: body
      })
    } else {
      response = await fetch("/api/v1/trips", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: body
      })
    }

    if (response.ok)
      onClose(true)
    else toast({
      title: "Error upserting Trip",
      description: await response.text()
    })
  }

  async function onDeleteButtonClick() {
    const response = await fetch("/api/v1/trips/" + trip!.id, {method: "DELETE"})

    if (response.ok)
      onClose(true)
    else toast({
      title: "Error deleting Trip",
      description: await response.text()
    })
  }

  return (
    <Dialog open={open} onOpenChange={open => !open ? onClose(false) : {}}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>
            {edit ? (trip != null ? "Edit" : "New") : "View"} Trip
          </DialogTitle>
        </DialogHeader>
        <div className="py-4 overflow-y-auto">
          <RowContainer>
            <LabelInputContainer>
              <Label htmlFor="act_name">Name</Label>
              <Input id="act_name" placeholder="My new Activity" type="text" value={name}
                     onChange={e => setName(e.target.value)} readOnly={!edit}/>
            </LabelInputContainer>
          </RowContainer>

          <RowContainer>
            <LabelInputContainer>
              <Label htmlFor="arrival_date">Arrival Date</Label>
              <DateInput date={startDate} updateDate={setStartDate} readOnly={!edit}/>
            </LabelInputContainer>
            <LabelInputContainer>
              <Label htmlFor="departure_date">Departure Date</Label>
              <DateInput date={endDate} updateDate={setEndDate} startDate={startDate} readOnly={!edit}/>
            </LabelInputContainer>
          </RowContainer>

          <RowContainer>
            <LabelInputContainer>
              <Label htmlFor="description">Description</Label>
              <Textarea id="description" value={description}
                        onChange={e => setDescription(e.target.value)}
                        readOnly={!edit}/>
            </LabelInputContainer>
          </RowContainer>

          <RowContainer>
            <LabelInputContainer>
              <Label htmlFor="image_url">Image URL</Label>
              <Input id="image_url" type="text" value={imageUrl}
                     onChange={e => setImageUrl(e.target.value)} readOnly={!edit}/>
            </LabelInputContainer>
          </RowContainer>
        </div>
        <DialogFooter>
          {edit ?
            <Button className="w-full text-base" onClick={onSaveButtonClick}>
              Save
            </Button>
          :
            <>
              <Button variant="destructive" className="w-full text-base" onClick={onDeleteButtonClick}>
                Delete
              </Button>
              <Button variant="secondary" className="w-full text-base" onClick={() => setEdit(true)}>
                Edit
              </Button>
            </>
          }
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}