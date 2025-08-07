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

export default function AddActivityDialog({
  trip, open, onClose
}: {
  trip: Trip,
  open: boolean,
  onClose: (needsUpdate: boolean) => void,
}) {
  const [name, setName] = useState<string>("")
  const [description, setDescription] = useState<string>("")
  const [date, setDate] = useState<Date>()
  const [price, setPrice] = useState<number|null>(null)
  const [address, setAddress] = useState<string>("")
  const [location, setLocation] = useState<Location|null>(null)

  const { toast } = useToast();

  async function onSaveButtonClick() {
    const response = await fetch("/api/v1/activities", {
      method: "POST",
      headers: {"Content-Type": "application/json"},
      body: JSON.stringify({
        tripId: trip.id,
        name: name,
        date: date ? getDateString(date) : null,
        time: null,
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
          <DialogTitle>Add Activity</DialogTitle>
        </DialogHeader>
        <div className="py-4 overflow-y-auto">
          <div className="mb-4 flex flex-col space-y-2 md:flex-row md:space-y-0 md:space-x-2">
            <LabelInputContainer>
              <Label htmlFor="act_name">Name</Label>
              <Input id="act_name" placeholder="My new Activity" type="text" value={name}
                     onChange={e => setName(e.target.value)}/>
            </LabelInputContainer>
          </div>
          <div className="mb-4 flex flex-col space-y-2 md:flex-row md:space-y-0 md:space-x-2">
            <LabelInputContainer>
              <Label htmlFor="description">Description</Label>
              <Textarea id="description" value={description}
                        onChange={e => setDescription(e.target.value)}/>
            </LabelInputContainer>
          </div>
          <div className="mb-4 flex flex-col space-y-2 md:flex-row md:space-y-0 md:space-x-2">
            <LabelInputContainer>
              <Label htmlFor="date">Date</Label>
              <Popover>
                <PopoverTrigger asChild>
                  <Button
                      variant="outline"
                      className={cn(
                          "col-span-3 justify-start text-left font-normal",
                          !date && "text-muted-foreground"
                      )}
                  >
                    <CalendarIcon className="mr-2 h-4 w-4"/>
                    {date ? format(date, "PPP") : <span>Pick a date</span>}
                  </Button>
                </PopoverTrigger>
                <PopoverContent className="w-auto p-0">
                  <Calendar
                      mode="single"
                      fromDate={trip.startDate}
                      toDate={trip.endDate}
                      selected={date}
                      onSelect={setDate}
                  />
                </PopoverContent>
              </Popover>
            </LabelInputContainer>
            <LabelInputContainer>
              <Label htmlFor="price">Price</Label>
              <AmountInput
                  id="price"
                  className="col-span-3"
                  amount={price}
                  updateAmount={setPrice}
              />
            </LabelInputContainer>
          </div>
          <div className="mb-4 flex flex-col space-y-2 md:flex-row md:space-y-0 md:space-x-2">
            <LabelInputContainer>
              <Label htmlFor="address">Address</Label>
              <AddressInput address={address} updateAddress={setAddress} updateLocation={setLocation}/>
            </LabelInputContainer>
          </div>
          <div className="mb-4 flex flex-col space-y-2 md:flex-row md:space-y-0 md:space-x-2">
            <LabelInputContainer>
              <Label htmlFor="lat">Latitude</Label>
              <Input id="lat" type="text" readOnly/>
            </LabelInputContainer>
            <LabelInputContainer>
              <Label htmlFor="lon">Longitude</Label>
              <Input id="lon" type="text" readOnly/>
            </LabelInputContainer>
          </div>

        </div>
        <DialogFooter>
          <button
              className="group/btn relative block h-10 w-full rounded-md bg-gradient-to-br from-black to-neutral-600 font-medium text-white shadow-[0px_1px_0px_0px_#ffffff40_inset,0px_-1px_0px_0px_#ffffff40_inset] dark:bg-zinc-800 dark:from-zinc-900 dark:to-zinc-900 dark:shadow-[0px_1px_0px_0px_#27272a_inset,0px_-1px_0px_0px_#27272a_inset]"
              type="submit"
              onClick={onSaveButtonClick}>
            Save
            <BottomGradient />
          </button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

const BottomGradient = () => {
  return (
      <>
        <span className="absolute inset-x-0 -bottom-px block h-px w-full bg-gradient-to-r from-transparent via-cyan-500 to-transparent opacity-0 transition duration-500 group-hover/btn:opacity-100" />
        <span className="absolute inset-x-10 -bottom-px mx-auto block h-px w-1/2 bg-gradient-to-r from-transparent via-indigo-500 to-transparent opacity-0 blur-sm transition duration-500 group-hover/btn:opacity-100" />
      </>
  );
};

const LabelInputContainer = ({
                               children,
                               className,
                             }: {
  children: React.ReactNode;
  className?: string;
}) => {
  return (
      <div className={cn("flex w-full flex-col space-y-2", className)}>
        {children}
      </div>
  );
};