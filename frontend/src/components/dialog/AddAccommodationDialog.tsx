import {Button} from "@/components/ui/button.tsx";
import {Dialog, DialogFooter, DialogContent, DialogHeader, DialogTitle} from "@/components/ui/dialog.tsx";
import {Label} from "@/components/ui/label.tsx";
import {Popover, PopoverContent, PopoverTrigger} from "@/components/ui/popover.tsx";
import {Calendar} from "@/components/ui/calendar.tsx";
import {CalendarIcon} from "lucide-react";
import {useState} from "react";
import { format } from "date-fns"
import { cn } from "@/lib/utils"
import {useToast} from "@/components/ui/use-toast.ts";
import AmountInput from "@/components/dialog/AmountInput.tsx";
import {Input} from "@/components/ui/input.tsx";
import {Location, Trip} from "@/types.ts";
import {Textarea} from "@/components/ui/textarea.tsx";
import AddressInput from "@/components/dialog/AddressInput.tsx";
import {getDateString, nullIfEmpty} from "@/components/util.ts";
import {LabelInputContainer, RowContainer} from "@/components/dialog/DialogUtil.tsx";

export default function AddAccommodationDialog({
  trip, open, onClose
}: {
  trip: Trip,
  open: boolean,
  onClose: (needsUpdate: boolean) => void,
}) {
  const [name, setName] = useState<string>("")
  const [description, setDescription] = useState<string>("")
  const [arrivalDate, setArrivalDate] = useState<Date>()
  const [departureDate, setDepartureDate] = useState<Date>()
  const [price, setPrice] = useState<number|null>(null)
  const [address, setAddress] = useState<string>("")
  const [location, setLocation] = useState<Location|null>(null)

  const { toast } = useToast();

  async function onSaveButtonClick() {
    const response = await fetch("/api/v1/accommodation", {
      method: "POST",
      headers: {"Content-Type": "application/json"},
      body: JSON.stringify({
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
    })

    if (response.ok)
      onClose(true)
    else toast({
      title: "Error adding Activity",
      description: response.statusText
    })
  }

  return (
    <Dialog open={open} onOpenChange={open => !open ? onClose(false) : {}}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Add Accommodation</DialogTitle>
        </DialogHeader>
        <div className="py-4 overflow-y-auto">
          <RowContainer>
            <LabelInputContainer>
              <Label htmlFor="act_name">Name</Label>
              <Input id="act_name" placeholder="My new Activity" type="text" value={name}
                     onChange={e => setName(e.target.value)}/>
            </LabelInputContainer>
          </RowContainer>

          <RowContainer>
            <LabelInputContainer>
              <Label htmlFor="arrival_date">Arrival Date</Label>
              <Popover>
                <PopoverTrigger asChild>
                  <Button
                      variant="secondary"
                      className={cn(
                          "col-span-3 justify-start text-left font-normal",
                          !arrivalDate && "text-muted-foreground"
                      )}
                  >
                    <CalendarIcon className="mr-2 h-4 w-4"/>
                    {arrivalDate ? format(arrivalDate, "PPP") : <span>Pick a date</span>}
                  </Button>
                </PopoverTrigger>
                <PopoverContent className="w-auto p-0 rounded-2xl overflow-hidden shadow-lg">
                  <Calendar
                      mode="single"
                      selected={arrivalDate}
                      onSelect={setArrivalDate}
                      startMonth={trip.startDate}
                      endMonth={trip.endDate}
                      disabled={{before: trip.startDate, after: trip.endDate}}
                  />
                </PopoverContent>
              </Popover>
            </LabelInputContainer>
            <LabelInputContainer>
              <Label htmlFor="departure_date">Departure Date</Label>
              <Popover>
                <PopoverTrigger asChild>
                  <Button
                      variant="secondary"
                      className={cn(
                          "col-span-3 justify-start text-left font-normal",
                          !departureDate && "text-muted-foreground"
                      )}
                  >
                    <CalendarIcon className="mr-2 h-4 w-4"/>
                    {departureDate ? format(departureDate, "PPP") : <span>Pick a date</span>}
                  </Button>
                </PopoverTrigger>
                <PopoverContent className="w-auto p-0 rounded-2xl overflow-hidden shadow-lg">
                  <Calendar
                      mode="single"
                      selected={departureDate}
                      onSelect={setDepartureDate}
                      startMonth={trip.startDate}
                      endMonth={trip.endDate}
                      disabled={{before: trip.startDate, after: trip.endDate}}
                  />
                </PopoverContent>
              </Popover>
            </LabelInputContainer>
          </RowContainer>

          <RowContainer>
            <LabelInputContainer>
              <Label htmlFor="description">Description</Label>
              <Textarea id="description" value={description}
                        onChange={e => setDescription(e.target.value)}/>
            </LabelInputContainer>
          </RowContainer>

          <RowContainer>
            <LabelInputContainer>
              <Label htmlFor="price">Price</Label>
              <AmountInput
                  id="price"
                  amount={price}
                  updateAmount={setPrice}
              />
            </LabelInputContainer>
          </RowContainer>

          <RowContainer>
            <LabelInputContainer>
              <Label htmlFor="address">Address</Label>
              <AddressInput address={address} updateAddress={setAddress} updateLocation={setLocation}/>
            </LabelInputContainer>
          </RowContainer>
          <RowContainer>
            <LabelInputContainer>
              <Label htmlFor="lat">Latitude</Label>
              <Input id="lat" type="text" disabled/>
            </LabelInputContainer>
            <LabelInputContainer>
              <Label htmlFor="lon">Longitude</Label>
              <Input id="lon" type="text" disabled/>
            </LabelInputContainer>
          </RowContainer>
        </div>
        <DialogFooter>
          <Button className="w-full text-base" onClick={onSaveButtonClick}>
            Save
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}