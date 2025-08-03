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
        <div className="grid gap-4 py-4 overflow-y-auto">
          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="act_name" className="text-right">
              Name
            </Label>
            <Input id="act_name" value={name}
                   onChange={e => setName(e.target.value)}
                   placeholder="My new Activity" className="col-span-3"/>
          </div>

          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="arr_date" className="text-right">
              Arr. Date
            </Label>
            <Popover>
              <PopoverTrigger asChild>
                <Button
                    variant="outline"
                    className={cn(
                        "col-span-3 justify-start text-left font-normal",
                        !arrivalDate && "text-muted-foreground"
                    )}
                >
                  <CalendarIcon className="mr-2 h-4 w-4"/>
                  {arrivalDate ? format(arrivalDate, "PPP") : <span>Pick a date</span>}
                </Button>
              </PopoverTrigger>
              <PopoverContent className="w-auto p-0">
                <Calendar
                    mode="single"
                    fromDate={trip.startDate}
                    toDate={trip.endDate}
                    selected={arrivalDate}
                    onSelect={setArrivalDate}
                />
              </PopoverContent>
            </Popover>
          </div>

          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="dep_date" className="text-right">
              Dep. Date
            </Label>
            <Popover>
              <PopoverTrigger asChild>
                <Button
                    variant="outline"
                    className={cn(
                        "col-span-3 justify-start text-left font-normal",
                        !departureDate && "text-muted-foreground"
                    )}
                >
                  <CalendarIcon className="mr-2 h-4 w-4"/>
                  {departureDate ? format(departureDate, "PPP") : <span>Pick a date</span>}
                </Button>
              </PopoverTrigger>
              <PopoverContent className="w-auto p-0">
                <Calendar
                    mode="single"
                    fromDate={trip.startDate}
                    toDate={trip.endDate}
                    selected={departureDate}
                    onSelect={setDepartureDate}
                />
              </PopoverContent>
            </Popover>
          </div>

          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="description" className="text-right">
              Description
            </Label>
            <Textarea id="description" value={description}
                      onChange={e => setDescription(e.target.value)}
                      className="col-span-3"/>
          </div>

          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="price" className="text-right">
              Price
            </Label>
            <AmountInput
                id="price"
                className="col-span-3"
                amount={price}
                updateAmount={setPrice}
            />
          </div>

          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="description" className="text-right">
              Address
            </Label>
            <AddressInput address={address} updateAddress={setAddress} updateLocation={setLocation}/>
          </div>

          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="coords" className="text-right">
              Coordinates
            </Label>
            <div className="col-span-3 flex gap-2">
              <Input id="lat" value={location?.latitude ?? ""}
                     disabled
                     className="cl-span-1"/>
              <Input id="lon" value={location?.longitude ?? ""}
                     disabled
                     className="cl-span-1"/>
            </div>
          </div>

        </div>
        <DialogFooter>
          <Button type="submit" onClick={onSaveButtonClick}>Save</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}