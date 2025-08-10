import {Button} from "@/components/ui/button.tsx";
import {Dialog, DialogFooter, DialogContent, DialogHeader, DialogTitle} from "@/components/ui/dialog.tsx";
import {Label} from "@/components/ui/label.tsx";
import {useState} from "react";
import {useToast} from "@/components/ui/use-toast.ts";
import AmountInput from "@/components/dialog/AmountInput.tsx";
import {Input} from "@/components/ui/input.tsx";
import {Accommodation, Location, Trip} from "@/types.ts";
import {Textarea} from "@/components/ui/textarea.tsx";
import AddressInput from "@/components/dialog/AddressInput.tsx";
import {getDateString, nullIfEmpty} from "@/components/util.ts";
import {LabelInputContainer, RowContainer} from "@/components/dialog/DialogUtil.tsx";
import DateInput from "@/components/dialog/DateInput.tsx";

export default function AccommodationDialog({
  trip, accommodation, open, onClose
}: {
  trip: Trip,
  open: boolean,
  onClose: (needsUpdate: boolean) => void,
  accommodation?: Accommodation | null
}) {
  const [edit, setEdit] = useState<boolean>(accommodation == null)

  const [name, setName] = useState<string>(accommodation?.name ?? "")
  const [description, setDescription] = useState<string>(accommodation?.description ?? "")
  const [arrivalDate, setArrivalDate] = useState<Date|null>(accommodation?.arrivalDate ?? null)
  const [departureDate, setDepartureDate] = useState<Date|null>(accommodation?.departureDate ?? null)
  const [price, setPrice] = useState<number|null>(accommodation?.price ?? null)
  const [address, setAddress] = useState<string>(accommodation?.address ?? "")
  const [location, setLocation] = useState<Location|null>(accommodation?.location ?? null)

  const { toast } = useToast();

  async function onSaveButtonClick() {
    const body = JSON.stringify({
      tripId: trip.id,
      name: name,
      arrivalDate: arrivalDate ? getDateString(arrivalDate) : null,
      departureDate: departureDate ? getDateString(departureDate) : null,
      checkInTime: null,
      checkOutTime: null,
      description: nullIfEmpty(description),
      address: nullIfEmpty(address),
      location: location,
      price: price,
    })

    let response
    if (accommodation != null) {
      response = await fetch("/api/v1/trips/" + trip.id + "/accommodation/" + accommodation?.id, {
        method: "PUT",
        headers: {"Content-Type": "application/json"},
        body: body
      })
    } else {
      response = await fetch("/api/v1/trips/" + trip.id + "/accommodation", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: body
      })
    }

    if (response.ok)
      onClose(true)
    else toast({
      title: "Error upserting Activity",
      description: response.statusText
    })
  }

  async function onDeleteButtonClick() {
    const response = await fetch("/api/v1/trips/" + trip.id + "/accommodation/" + accommodation!.id, {method: "DELETE"})

    if (response.ok)
      onClose(true)
    else toast({
      title: "Error deleting Accommodation",
      description: response.statusText
    })
  }

  return (
    <Dialog open={open} onOpenChange={open => !open ? onClose(false) : {}}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>
            {edit ? (accommodation != null ? "Edit" : "New") : "View"} Accommodation
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
              <DateInput date={arrivalDate} updateDate={setArrivalDate} startDate={trip.startDate} endDate={trip.endDate} readOnly={!edit}/>
            </LabelInputContainer>
            <LabelInputContainer>
              <Label htmlFor="departure_date">Departure Date</Label>
              <DateInput date={departureDate} updateDate={setDepartureDate} startDate={arrivalDate ?? trip.startDate} endDate={trip.endDate} readOnly={!edit}/>
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