import {Button} from "@/components/ui/button.tsx";
import {Dialog, DialogFooter, DialogContent, DialogHeader, DialogTitle} from "@/components/ui/dialog.tsx";
import {Label} from "@/components/ui/label.tsx";
import {useState} from "react";
import {useToast} from "@/components/ui/use-toast.ts";
import AmountInput from "@/components/dialog/AmountInput.tsx";
import {Input} from "@/components/ui/input.tsx";
import {Activity, Location, Trip} from "@/types.ts";
import {Textarea} from "@/components/ui/textarea.tsx";
import AddressInput from "@/components/dialog/AddressInput.tsx";
import {getDateString, nullIfEmpty} from "@/components/util.ts";
import {LabelInputContainer, RowContainer } from "./DialogUtil";
import DateInput from "@/components/dialog/DateInput.tsx";

export default function ActivityDialog({
  trip, open, onClose, activity
}: {
  trip: Trip
  open: boolean
  onClose: (needsUpdate: boolean) => void
  activity?: Activity | null
}) {
  const [edit, setEdit] = useState<boolean>(activity == null)

  const [name, setName] = useState<string>(activity?.name ?? "")
  const [description, setDescription] = useState<string>(activity?.description ?? "")
  const [date, setDate] = useState<Date|null>(activity?.date ?? null)
  const [price, setPrice] = useState<number|null>(activity?.price ?? null)
  const [address, setAddress] = useState<string>(activity?.address ?? "")
  const [location, setLocation] = useState<Location|null>(activity?.location ?? null)

  const { toast } = useToast();

  async function onSaveButtonClick() {
    const body = JSON.stringify({
      tripId: trip.id,
      name: name,
      date: date ? getDateString(date) : null,
      time: null,
      description: nullIfEmpty(description),
      address: nullIfEmpty(address),
      location: location,
      price: price,
    })

    let response
    if (activity != null) {
      response = await fetch("/api/v1/trips/" + trip.id + "/activities/" + activity?.id, {
        method: "PUT",
        headers: {"Content-Type": "application/json"},
        body: body
      })
    } else {
      response = await fetch("/api/v1/trips/" + trip.id + "/activities", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: body
      })
    }

    if (response.ok)
      onClose(true)
    else toast({
      title: "Error upserting Activity",
      description: await response.text()
    })
  }

  async function onDeleteButtonClick() {
    const response = await fetch("/api/v1/trips/" + trip.id + "/activities/" + activity!.id, {method: "DELETE"})

    if (response.ok)
      onClose(true)
    else toast({
      title: "Error deleting Activity",
      description: await response.text()
    })
  }

  return (
    <Dialog open={open} onOpenChange={open => !open ? onClose(false) : {}}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>
            {edit ? (activity != null ? "Edit" : "New") : "View"} Activity
          </DialogTitle>
        </DialogHeader>
        <div className="py-4 overflow-y-auto">
          <RowContainer>
            <LabelInputContainer>
              <Label htmlFor="act_name">Name</Label>
              <Input id="act_name" placeholder="My new Activity" type="text" value={name}
                     onChange={e => setName(e.target.value)}
                     readOnly={!edit}/>
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
              <Label htmlFor="date">Date</Label>
              <DateInput date={date} updateDate={setDate} startDate={trip.startDate} endDate={trip.endDate} readOnly={!edit}/>
            </LabelInputContainer>
            <LabelInputContainer>
              <Label htmlFor="price">Price</Label>
              <AmountInput
                  id="price"
                  amount={price}
                  updateAmount={setPrice}
                  readOnly={!edit}
              />
            </LabelInputContainer>
          </RowContainer>
          <RowContainer>
            <LabelInputContainer>
              <Label htmlFor="address">Address</Label>
              <AddressInput address={address} updateAddress={setAddress} updateLocation={setLocation} readOnly={!edit}/>
            </LabelInputContainer>
          </RowContainer>
          <RowContainer>
            <LabelInputContainer>
              <Label htmlFor="lat">Latitude</Label>
              <Input id="lat" type="text" value={location?.latitude ?? ""} readOnly/>
            </LabelInputContainer>
            <LabelInputContainer>
              <Label htmlFor="lon">Longitude</Label>
              <Input id="lon" type="text" value={location?.longitude ?? ""} readOnly/>
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