import {Button} from "@/components/ui/button.tsx";
import {Dialog, DialogFooter, DialogContent, DialogHeader, DialogTitle} from "@/components/ui/dialog.tsx";
import {Label} from "@/components/ui/label.tsx";
import {Popover, PopoverContent, PopoverTrigger} from "@/components/ui/popover.tsx";
import {Calendar} from "@/components/ui/calendar.tsx";
import {CalendarIcon, Minus, Plus} from "lucide-react";
import {useState} from "react";
import { format } from "date-fns"
import { cn } from "@/lib/utils"
import {useToast} from "@/components/ui/use-toast.ts";
import AmountInput from "@/components/dialog/AmountInput.tsx";
import {Input} from "@/components/ui/input.tsx";
import {AddFlightLeg, AddPNR, Trip} from "@/types.ts";
import {getDateString} from "@/components/util.ts";
import {LabelInputContainer, RowContainer} from "@/components/dialog/DialogUtil.tsx";

export default function AddFlightDialog({
  trip, open, onClose
}: {
  trip: Trip,
  open: boolean,
  onClose: (needsUpdate: boolean) => void,
}) {
  const [price, setPrice] = useState<number|null>(null)
  const [flightLegs, setFlightLegs] = useState<AddFlightLeg[]>([{
    date: getDateString(trip.startDate),
    flightNumber: "",
    originAirport: null
  }])
  const [pnrs, setPnrs] = useState<AddPNR[]>([])

  const { toast } = useToast();

  async function onSaveButtonClick() {
    const response = await fetch("/api/v1/flights", {
      method: "POST",
      headers: {"Content-Type": "application/json"},
      body: JSON.stringify({
        tripId: trip.id,
        legs: flightLegs,
        pnrs: pnrs,
        price: price,
      })
    })

    if (response.ok)
      onClose(true)
    else toast({
      title: "Error adding Flight",
      description: response.statusText
    })
  }

  function addLeg() {
    setFlightLegs(legs => [...legs, {
      date: legs[legs.length-1].date,
      flightNumber: "",
      originAirport: null
    }])
  }

  function deleteLeg() {
    setFlightLegs(legs => legs.slice(0, legs.length-1))
  }

  function addPnr() {
    setPnrs(pnrs => [...pnrs, {airline: "", pnr: ""}])
  }

  function deletePnr() {
    setPnrs(pnrs => pnrs.slice(0, pnrs.length-1))
  }

  function updateFlightLeg(updateIdx: number, updateFn: (leg: AddFlightLeg) => void) {
    setFlightLegs(legs => legs.map((leg, idx) => {
          if (idx == updateIdx)
            updateFn(leg)
          return leg
        })
    )
  }

  function updatePnr(updateIdx: number, updateFn: (pnr: AddPNR) => void) {
    setPnrs(pnrs => pnrs.map((pnr, idx) => {
          if (idx == updateIdx)
            updateFn(pnr)
          return pnr
        })
    )
  }

  return (
    <Dialog open={open} onOpenChange={open => !open ? onClose(false) : {}}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Add Flight</DialogTitle>
        </DialogHeader>
        <div className="py-4 overflow-y-auto">
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

          <div className="flex">
            <h3 className="font-semibold mb-2 flex-grow">Flight Legs</h3>
            <div className="">
              {flightLegs.length > 1 ?
                <Button variant="ghost" className="p-2 h-auto rounded-full" onClick={() => deleteLeg()}>
                  <Minus className="h-4 w-4"/>
                </Button>
                :
                <div/>
              }
              <Button variant="ghost" className="p-2 h-auto rounded-full" onClick={() => addLeg()}>
                <Plus className="w-4 h-4"/>
              </Button>
            </div>
          </div>

          {flightLegs.map((leg, idx) =>
              <div key={idx}>
                <RowContainer>
                  <LabelInputContainer>
                    <Label htmlFor={"date" + idx}>
                      Date {flightLegs.length > 1 ? (idx+1) : ""}
                    </Label>
                    <Popover>
                      <PopoverTrigger asChild>
                        <Button
                            variant="outline"
                            className={cn(
                                "col-span-3 justify-start text-left font-normal",
                                !leg.date && "text-muted-foreground"
                            )}
                        >
                          <CalendarIcon className="mr-2 h-4 w-4"/>
                          {leg.date ? format(leg.date, "PPP") : <span>Pick a date</span>}
                        </Button>
                      </PopoverTrigger>
                      <PopoverContent className="w-auto p-0">
                        <Calendar
                            mode="single"
                            startMonth={trip.startDate}
                            endMonth={trip.endDate}
                            disabled={{before: trip.startDate, after: trip.endDate}}
                            selected={new Date(leg.date)}
                            onSelect={date => updateFlightLeg(idx, leg => leg.date = getDateString(date!))}
                        />
                      </PopoverContent>
                    </Popover>
                  </LabelInputContainer>
                  <LabelInputContainer>
                    <Label htmlFor={"flightno" + idx}>
                      Flight #{flightLegs.length > 1 ? (idx+1) : ""}
                    </Label>
                    <Input id={"flightno" + idx} value={leg.flightNumber}
                           onChange={e => updateFlightLeg(idx, leg => leg.flightNumber = e.target.value)}
                           placeholder="LH717"/>
                    </LabelInputContainer>
                  </RowContainer>
              </div>
          )}

          {/*<Separator/>*/}
          <div className="flex">
            <h3 className="font-semibold mb-2 flex-grow">PNRs</h3>
            <div className="">
              {pnrs.length > 0 ?
                <Button variant="ghost" className="p-2 h-auto rounded-full" onClick={() => deletePnr()}>
                  <Minus className="h-4 w-4"/>
                </Button>
                :
                <div/>
              }
              <Button variant="ghost" className="p-2 h-auto rounded-full" onClick={() => addPnr()}>
                <Plus className="w-4 h-4"/>
              </Button>
            </div>
          </div>

          {pnrs.map((pnr, idx) =>
            <div key={idx}>
              <RowContainer>
                <LabelInputContainer>
                  <Label htmlFor={"flightno" + idx}>
                    Airline {pnrs.length > 1 ? (idx+1) : ""}
                  </Label>
                  <Input id={"airline" + idx} value={pnr.airline}
                         onChange={e => updatePnr(idx, pnr => pnr.airline = e.target.value)}
                         placeholder="LH"/>
                </LabelInputContainer>
                <LabelInputContainer>
                  <Label htmlFor={"pnr" + idx}>
                    PNR {pnrs.length > 1 ? (idx+1) : ""}
                  </Label>
                  <Input id={"pnr" + idx} value={pnr.pnr}
                         onChange={e => updatePnr(idx, pnr => pnr.pnr = e.target.value)}
                         placeholder="123ABC"/>
                </LabelInputContainer>
              </RowContainer>
            </div>
          )}
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