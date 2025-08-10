import {Button} from "@/components/ui/button.tsx";
import {Dialog, DialogFooter, DialogContent, DialogHeader, DialogTitle} from "@/components/ui/dialog.tsx";
import {Label} from "@/components/ui/label.tsx";
import {Minus, Plus} from "lucide-react";
import {useState} from "react";
import {useToast} from "@/components/ui/use-toast.ts";
import AmountInput from "@/components/dialog/AmountInput.tsx";
import {Input} from "@/components/ui/input.tsx";
import {AddFlightLeg, AddPNR, Flight, FlightLeg, Trip} from "@/types.ts";
import {LabelInputContainer, RowContainer} from "@/components/dialog/DialogUtil.tsx";
import DateInput from "@/components/dialog/DateInput.tsx";

export default function FlightDialog({
  trip, flight, open, onClose
}: {
  trip: Trip
  open: boolean
  onClose: (needsUpdate: boolean) => void
  flight?: Flight | null
}) {
  const [edit, setEdit] = useState<boolean>(flight == null)

  const [price, setPrice] = useState<number|null>(flight?.price ?? null)
  const [flightLegs, setFlightLegs] = useState<AddFlightLeg[]>(mapLegsOrDefault(flight?.legs))

  function mapLegsOrDefault(flightLegs: FlightLeg[]|undefined): AddFlightLeg[] {
    if (flightLegs != null) {
      return flightLegs.map(leg => ({
        date: trip.startDate,
        flightNumber: leg.flightNumber,
        originAirport: leg.origin.iata
      }))
    }

    return [{
      date: trip.startDate,
      flightNumber: "",
      originAirport: null
    }]
  }

  const [pnrs, setPnrs] = useState<AddPNR[]>(flight?.pnrs ?? [])

  const { toast } = useToast();

  async function onSaveButtonClick() {
    const body = JSON.stringify({
      tripId: trip.id,
      legs: flightLegs,
      pnrs: pnrs,
      price: price,
    })

    let response
    if (flight != null) {
      response = await fetch("/api/v1/trips/" + trip.id + "/flights/" + flight?.id, {
        method: "PUT",
        headers: {"Content-Type": "application/json"},
        body: body
      })
    } else {
      response = await fetch("/api/v1/trips/" + trip.id + "/flights", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: body
      })
    }

    if (response.ok)
      onClose(true)
    else toast({
      title: "Error upserting Flight",
      description: await response.text()
    })
  }

  async function onDeleteButtonClick() {
    const response = await fetch("/api/v1/trips/" + trip.id + "/flights/" + flight!.id, {method: "DELETE"})

    if (response.ok)
      onClose(true)
    else toast({
      title: "Error deleting Flight",
      description: await response.text()
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
          <DialogTitle>
            {edit ? (flight != null ? "Edit" : "New") : "View"} Flight
          </DialogTitle>
        </DialogHeader>
        <div className="py-4 overflow-y-auto">
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

          <div className="flex">
            <h3 className="font-semibold mb-2 flex-grow">Flight Legs</h3>
            {edit &&
              <div>
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
            }
          </div>

          {flightLegs.map((leg, idx) =>
              <div key={idx}>
                <RowContainer>
                  <LabelInputContainer>
                    <Label htmlFor={"date" + idx}>
                      Date {flightLegs.length > 1 ? (idx+1) : ""}
                    </Label>
                    <DateInput date={leg.date} updateDate={date => updateFlightLeg(idx, leg => leg.date = date)}
                               startDate={trip.startDate} endDate={trip.endDate} readOnly={!edit}/>
                  </LabelInputContainer>
                  <LabelInputContainer>
                    <Label htmlFor={"flightno" + idx}>
                      Flight #{flightLegs.length > 1 ? (idx+1) : ""}
                    </Label>
                    <Input id={"flightno" + idx} value={leg.flightNumber}
                           onChange={e => updateFlightLeg(idx, leg => leg.flightNumber = e.target.value)}
                           placeholder="LH717" readOnly={!edit}/>
                    </LabelInputContainer>
                  </RowContainer>
              </div>
          )}

          {/*<Separator/>*/}
          <div className="flex">
            <h3 className="font-semibold mb-2 flex-grow">PNRs</h3>
            {edit &&
              <div>
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
            }
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
                         placeholder="LH" readOnly={!edit}/>
                </LabelInputContainer>
                <LabelInputContainer>
                  <Label htmlFor={"pnr" + idx}>
                    PNR {pnrs.length > 1 ? (idx+1) : ""}
                  </Label>
                  <Input id={"pnr" + idx} value={pnr.pnr}
                         onChange={e => updatePnr(idx, pnr => pnr.pnr = e.target.value)}
                         placeholder="123ABC" readOnly={!edit}/>
                </LabelInputContainer>
              </RowContainer>
            </div>
          )}
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
              <Button variant="secondary" className="w-full text-base" disabled onClick={() => setEdit(true)}>
                Edit
              </Button>
            </>
          }
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}