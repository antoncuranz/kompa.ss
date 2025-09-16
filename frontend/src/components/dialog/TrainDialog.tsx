import {Button} from "@/components/ui/button.tsx";
import {Dialog, DialogFooter, DialogContent, DialogHeader, DialogTitle} from "@/components/ui/dialog.tsx";
import {Label} from "@/components/ui/label.tsx";
import {Minus, Plus} from "lucide-react";
import {useState} from "react";
import {useToast} from "@/components/ui/use-toast.ts";
import AmountInput from "@/components/dialog/AmountInput.tsx";
import {Input} from "@/components/ui/input.tsx";
import {TrainStation, Transportation, Trip} from "@/types.ts";
import {LabelInputContainer, RowContainer} from "@/components/dialog/DialogUtil.tsx";
import DateInput from "@/components/dialog/DateInput.tsx";
import TrainStationInput from "@/components/dialog/TrainStationInput.tsx";

export default function TrainDialog({
  trip, train, open, onClose
}: {
  trip: Trip
  open: boolean
  onClose: (needsUpdate: boolean) => void
  train?: Transportation | null
}) {
  const [edit, setEdit] = useState<boolean>(train == null)

  const [date, setDate] = useState<string|null>(null)
  const [price, setPrice] = useState<number|null>(train?.price ?? null)
  const [trainLegs, setTrainLegs] = useState<string[]>(mapLegsOrDefault(train))

  const [fromStation, setFromStation] = useState<TrainStation|null>(null)
  const [toStation, setToStation] = useState<TrainStation|null>(null)

  function mapLegsOrDefault(train: Transportation|null|undefined): string[] {
    return train?.trainDetail?.legs.map(leg => leg.lineName) ?? [""];
  }

  const { toast } = useToast();

  async function onSaveButtonClick() {
    const body = JSON.stringify({
      departureDate: date,
      fromStationId: fromStation?.id,
      toStationId: toStation?.id,
      trainNumbers: trainLegs,
    })

    let response
    if (train != null) {
      response = await fetch("/api/v1/trips/" + trip.id + "/trains/" + train?.id, {
        method: "PUT",
        headers: {"Content-Type": "application/json"},
        body: body
      })
    } else {
      response = await fetch("/api/v1/trips/" + trip.id + "/trains", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: body
      })
    }

    if (response.ok)
      onClose(true)
    else toast({
      title: "Error upserting Train",
      description: await response.text()
    })
  }

  async function onDeleteButtonClick() {
    const response = await fetch("/api/v1/trips/" + trip.id + "/transportation/" + train!.id, {method: "DELETE"})

    if (response.ok)
      onClose(true)
    else toast({
      title: "Error deleting Train",
      description: await response.text()
    })
  }

  function addLeg() {
    setTrainLegs(legs => [...legs, ""])
  }

  function deleteLeg() {
    setTrainLegs(legs => legs.slice(0, legs.length-1))
  }

  function updateTrainLeg(updateIdx: number, newLineName: string) {
    setTrainLegs(legs => {
      const updatedLegs = [...legs];
      updatedLegs[updateIdx] = newLineName
      return updatedLegs
    })
  }

  function validateForm(): boolean {
    return fromStation != null &&
        toStation != null &&
        date != null &&
        trainLegs.find(leg => leg.length == 0) == undefined
  }

  return (
    <Dialog open={open} onOpenChange={open => !open ? onClose(false) : {}}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>
            {edit ? (train != null ? "Edit" : "New") : "View"} Train Journey
          </DialogTitle>
        </DialogHeader>
        <div className="py-4 overflow-y-auto">
          <RowContainer>
            <LabelInputContainer>
              <Label htmlFor="from">From Station</Label>
              <TrainStationInput
                  station={fromStation}
                  updateStation={setFromStation}
                  placeholder="Berlin Südkreuz"
                  readOnly={!edit}/>
            </LabelInputContainer>
          </RowContainer>

          <RowContainer>
            <LabelInputContainer>
              <Label htmlFor="to">To Station</Label>
              <TrainStationInput
                  station={toStation}
                  updateStation={setToStation}
                  placeholder="München"
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

          <div className="flex">
            <h3 className="font-semibold mb-2 flex-grow">Train Legs</h3>
            {edit &&
              <div>
                {trainLegs.length > 1 ?
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

          {trainLegs.map((lineName, idx) =>
              <div key={idx}>
                <RowContainer>
                  <LabelInputContainer>
                    <Label htmlFor={"trainno" + idx}>
                      Train #{trainLegs.length > 1 ? (idx+1) : ""}
                    </Label>
                    <Input id={"trainno" + idx} value={lineName}
                           onChange={e => updateTrainLeg(idx, e.target.value)}
                           placeholder="ICE707" readOnly={!edit}/>
                    </LabelInputContainer>
                  </RowContainer>
              </div>
          )}
        </div>
        <DialogFooter>
          {edit ?
            <Button className="w-full text-base" onClick={onSaveButtonClick} disabled={!validateForm()}>
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