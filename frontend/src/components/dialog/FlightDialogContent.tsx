import { Button } from "@/components/ui/button.tsx";
import {
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog.tsx";
import { Minus, Plus } from "lucide-react";
import { useState } from "react";
import { Input } from "@/components/ui/input.tsx";
import {
  FlightLeg,
  Trip,
  AmbiguousFlightChoice,
  Flight,
} from "@/schema.ts";
import {
  Dialog,
  RowContainer,
  useDialogContext,
} from "@/components/dialog/Dialog.tsx";
import { z } from "zod";
import { isoDate, optionalString } from "@/schema.ts";
import { useFieldArray, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { dateFromString } from "@/components/util.ts";
import { Form, FormField } from "@/components/ui/form.tsx";
import AmountInput from "@/components/dialog/input/AmountInput.tsx";
import DateInput from "@/components/dialog/input/DateInput.tsx";
import { Spinner } from "@/components/ui/shadcn-io/spinner";
import { AmbiguousFlightDialogContent } from "@/components/dialog/AmbiguousFlightDialogContent.tsx";

const formSchema = z.object({
  legs: z.array(z.object({
    date: isoDate("Required"),
    flightNumber: z.string().nonempty("Required"),
    originAirport: optionalString(),
  })),
  pnrs: z.array(z.object({
    airline: z.string().nonempty("Required"),
    pnr: z.string().nonempty("Required"),
  })),
  price: z.number().optional()
})

type AmbiguousDialogData = {
  legs: { date: string; flightNumber: string }[];
  choices: { [flightNumber: string]: AmbiguousFlightChoice[] };
}

export default function FlightDialogContent({
  trip,
  flight,
}: {
  trip: Trip
  flight?: Flight
}) {
  const [edit] = useState<boolean>(flight == undefined);
  const [ambiguousDialogOpen, setAmbiguousDialogOpen] = useState(false);
  const [ambiguousDialogData, setAmbiguousDialogData] = useState<AmbiguousDialogData>({ legs: [], choices: {} });
  const { onClose } = useDialogContext();

  function mapLegsOrDefault(flightLegs: FlightLeg[]|undefined) {
    if (flightLegs) {
      return flightLegs.map((leg) => ({
        date: dateFromString(trip.startDate),
        flightNumber: leg.flightNumber
      }))
    }

    return [{
      date: dateFromString(trip.startDate),
      flightNumber: ""
    }]
  }

  const form = useForm<
    z.input<typeof formSchema>,
    unknown,
    z.output<typeof formSchema>
  >({
    resolver: zodResolver(formSchema),
    defaultValues: {
      legs: mapLegsOrDefault(flight?.legs.filter(leg => leg !== null)),
      pnrs: flight?.pnrs.filter(pnr => pnr !== null) ?? [],
      price: flight?.price ?? undefined
    },
    disabled: !edit
  })
  const { isSubmitting } = form.formState

  const legsArray = useFieldArray({
    control: form.control,
    name: "legs",
    rules: {
      minLength: 1
    }
  })

  const pnrsArray = useFieldArray({
    control: form.control,
    name: "pnrs"
  })

  async function onSubmit(values: z.infer<typeof formSchema>) {
    console.log(values)
    // if (flight) {
    //   flight.$jazz.applyDiff({ price: values.price })

    //   for (let i = 0; i < values.legs.length; i++) {
    //     if (flight.legs.length > i) {
    //       flight.legs[i].$jazz.applyDiff(values.legs[i])
    //     } else {
    //       flight.legs.$jazz.push(values.legs[i])
    //     }
    //   }
    //   while (flight.legs.length > values.legs.length) {
    //     flight.legs.$jazz.pop()
    //   }

    //   while (flight.pnrs.length > values.pnrs.length) {
    //     flight.pnrs.$jazz.pop()
    //   }
    // } else {
    //   trip.transportation.$jazz.push({
    //     type: "flight",
    //     legs: [],
    //     pnrs: [],
    //     price: values.price
    //   })
    // }
    onClose()
  }

  function handleAmbiguousFlightSelection(selectedFlights: Map<number, AmbiguousFlightChoice>) {
    setAmbiguousDialogOpen(false)
    setAmbiguousDialogData({ legs: [], choices: {} })

    for (const [legId, flight] of selectedFlights) {
      form.setValue(`legs.${legId}.originAirport`, flight.originIata)
    }

    form.handleSubmit(onSubmit)()
  }

  async function onDeleteButtonClick() {
    if (flight === undefined) {
      return
    }

    trip.transportation.$jazz.remove(t => t && t.$jazz.id == flight.$jazz.id)
    onClose()
  }

  async function onUpdateButtonClick() {
    if (flight === undefined) {
      return
    }

    // TODO!
  }

  function addLeg() {
    legsArray.append({
      date: legsArray.fields[legsArray.fields.length - 1].date,
      flightNumber: ""
    })
  }

  function deleteLeg() {
    legsArray.remove(legsArray.fields.length - 1)
  }

  function addPnr() {
    pnrsArray.append({ airline: "", pnr: "" })
  }

  function deletePnr() {
    pnrsArray.remove(pnrsArray.fields.length - 1)
  }

  return (
    <>
      <DialogHeader>
        <DialogTitle>
          {edit ? (flight != null ? "Edit" : "New") : "View"} Flight
        </DialogTitle>
      </DialogHeader>
      <Form id="flight-form" form={form} onSubmit={form.handleSubmit(onSubmit)}>
        <div className="flex">
          <h3 className="font-semibold mb-2 grow">Flight Legs</h3>
          {edit && (
            <div>
              {legsArray.fields.length > 1 ? (
                <Button
                  variant="ghost"
                  className="p-2 h-auto rounded-full"
                  onClick={() => deleteLeg()}
                >
                  <Minus className="h-4 w-4" />
                </Button>
              ) : (
                <div />
              )}
              <Button
                variant="ghost"
                className="p-2 h-auto rounded-full"
                onClick={() => addLeg()}
              >
                <Plus className="w-4 h-4" />
              </Button>
            </div>
          )}
        </div>

        {legsArray.fields.map((field, index) => (
          <div key={field.id}>
            <RowContainer>
              <FormField
                control={form.control}
                name={`legs.${index}.date`}
                label={`Date ${legsArray.fields.length > 1 ? index + 1 : ""}`}
                render={({ field }) => (
                  <DateInput
                    startDate={trip.startDate}
                    endDate={trip.endDate}
                    {...field}
                  />
                )}
              />
              <FormField
                control={form.control}
                name={`legs.${index}.flightNumber`}
                label={`Flight #${legsArray.fields.length > 1 ? index + 1 : ""}`}
                render={({ field }) => <Input placeholder="LH717" {...field} />}
              />
            </RowContainer>
          </div>
        ))}

        <div className="flex">
          <h3 className="font-semibold mb-2 grow">PNRs</h3>
          {edit && (
            <div>
              {pnrsArray.fields.length > 0 ? (
                <Button
                  variant="ghost"
                  className="p-2 h-auto rounded-full"
                  onClick={() => deletePnr()}
                >
                  <Minus className="h-4 w-4" />
                </Button>
              ) : (
                <div />
              )}
              <Button
                variant="ghost"
                className="p-2 h-auto rounded-full"
                onClick={() => addPnr()}
              >
                <Plus className="w-4 h-4" />
              </Button>
            </div>
          )}
        </div>

        {pnrsArray.fields.map((field, index) => (
          <div key={field.id}>
            <RowContainer>
              <FormField
                control={form.control}
                name={`pnrs.${index}.airline`}
                label={`Airline ${pnrsArray.fields.length > 1 ? index + 1 : ""}`}
                render={({ field }) => <Input placeholder="LH" {...field} />}
              />
              <FormField
                control={form.control}
                name={`pnrs.${index}.pnr`}
                label={`PNR ${pnrsArray.fields.length > 1 ? index + 1 : ""}`}
                render={({ field }) => (
                  <Input placeholder="123ABC" {...field} />
                )}
              />
            </RowContainer>
          </div>
        ))}

        <FormField
          control={form.control}
          name="price"
          label="Price"
          render={({ field }) => <AmountInput {...field} />}
        />
      </Form>
      <Dialog open={ambiguousDialogOpen} setOpen={setAmbiguousDialogOpen}>
        <AmbiguousFlightDialogContent
          flightLegs={ambiguousDialogData.legs}
          flightChoices={ambiguousDialogData.choices}
          onSelection={handleAmbiguousFlightSelection}
        />
      </Dialog>
      <DialogFooter>
        {edit ?
          <Button form="flight-form" type="submit" className="w-full" disabled={isSubmitting}>
            {isSubmitting ? <Spinner variant="pinwheel"/> : "Save"}
          </Button>
        :
          <>
            <Button variant="destructive" className="w-full" onClick={onDeleteButtonClick}>
              Delete
            </Button>
            <Button variant="secondary" className="w-full" onClick={onUpdateButtonClick}>
              Refresh Data
            </Button>
          </>
        }
      </DialogFooter>
    </>
  )
}
