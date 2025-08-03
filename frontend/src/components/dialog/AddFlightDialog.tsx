import {Button} from "@/components/ui/button.tsx";
import {Dialog, DialogFooter, DialogContent, DialogHeader, DialogTitle} from "@/components/ui/dialog.tsx";
import {Label} from "@/components/ui/label.tsx";
import {Popover, PopoverContent, PopoverTrigger} from "@/components/ui/popover.tsx";
import {Calendar} from "@/components/ui/calendar.tsx";
import {CalendarIcon, Delete} from "lucide-react";
import {useState} from "react";
import { format } from "date-fns"
import { cn } from "@/lib/utils"
import {useToast} from "@/components/ui/use-toast.ts";
import AmountInput from "@/components/dialog/AmountInput.tsx";
import {Input} from "@/components/ui/input.tsx";
import {Separator} from "@/components/ui/separator.tsx";
import {AddFlightLeg, AddPNR, Trip} from "@/types.ts";
import {getDateString} from "@/components/util.ts";

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
        <div className="grid gap-4 py-4 overflow-y-auto">

          {flightLegs.map((leg, idx) =>
              <div key={idx}>
                <div className="grid grid-cols-4 items-center gap-4 mb-4">
                  <Label htmlFor={"date" + idx} className="text-right">
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
                          fromDate={trip.startDate}
                          toDate={trip.endDate}
                          selected={new Date(leg.date)}
                          onSelect={date => updateFlightLeg(idx, leg => leg.date = getDateString(date!))}
                      />
                    </PopoverContent>
                  </Popover>
                </div>
                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor={"flightno" + idx} className="text-right">
                    Flight #{flightLegs.length > 1 ? (idx+1) : ""}
                  </Label>
                  <Input id={"flightno" + idx} value={leg.flightNumber}
                         onChange={e => updateFlightLeg(idx, leg => leg.flightNumber = e.target.value)}
                         placeholder="LH717" className="col-span-3"/>
                </div>
              </div>
          )}
          <div className="grid grid-cols-4 items-center gap-4">
            {flightLegs.length > 1 ?
                <Button variant="outline" className="col-span-1 border-dashed hover:border-solid" onClick={() => deleteLeg()}>
                  <Delete className="h-5 w-5"/>
                </Button>
                :
                <div/>
            }
            <Button variant="outline" className="col-span-3 border-dashed hover:border-solid" onClick={() => addLeg()}>
              Add Flight leg
            </Button>
          </div>

          <Separator/>

          {pnrs.map((pnr, idx) =>
            <div key={idx}>
              <div className="grid grid-cols-4 items-center gap-4 mb-4">
                <Label htmlFor={"flightno" + idx} className="text-right">
                  Airline {pnrs.length > 1 ? (idx+1) : ""}
                </Label>
                <Input id={"airline" + idx} value={pnr.airline}
                       onChange={e => updatePnr(idx, pnr => pnr.airline = e.target.value)}
                       placeholder="LH" className="col-span-3"/>
              </div>
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor={"pnr" + idx} className="text-right">
                  PNR {pnrs.length > 1 ? (idx+1) : ""}
                </Label>
                <Input id={"pnr" + idx} value={pnr.pnr}
                       onChange={e => updatePnr(idx, pnr => pnr.pnr = e.target.value)}
                       placeholder="123ABC" className="col-span-3"/>
              </div>
            </div>
          )}
          <div className="grid grid-cols-4 items-center gap-4">
            {pnrs.length > 0 ?
              <Button variant="outline" className="col-span-1 border-dashed hover:border-solid" onClick={() => deletePnr()}>
                <Delete className="h-5 w-5"/>
              </Button>
            :
              <div/>
            }
            <Button variant="outline" className="col-span-3 border-dashed hover:border-solid" onClick={() => addPnr()}>
              Add PNR
            </Button>
          </div>

          <Separator/>

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

        </div>
        <DialogFooter>
          <Button type="submit" onClick={onSaveButtonClick}>Save</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}