import {Button} from "@/components/ui/button.tsx";
import {DialogFooter, DialogHeader, DialogTitle} from "@/components/ui/dialog.tsx";
import {Minus, Plus} from "lucide-react";
import {useState} from "react";
import {Input} from "@/components/ui/input.tsx";
import {Transportation, FlightLeg, Trip} from "@/types.ts";
import {RowContainer, useDialogContext} from "@/components/dialog/Dialog.tsx";
import { toast } from "sonner";
import {z} from "zod"
import {isoDate} from "@/schemas.ts";
import {useFieldArray, useForm} from "react-hook-form";
import {zodResolver} from "@hookform/resolvers/zod";
import {dateFromString} from "@/components/util.ts";
import {Form, FormField} from "@/components/ui/form.tsx";
import AmountInput from "@/components/dialog/AmountInput.tsx";
import DateInput from "@/components/dialog/DateInput.tsx";
import {Spinner} from "@/components/ui/shadcn-io/spinner";

const formSchema = z.object({
  legs: z.array(z.object({
    date: isoDate("Required"),
    flightNumber: z.string().nonempty("Required"),
  })),
  pnrs: z.array(z.object({
    airline: z.string().nonempty("Required"),
    pnr: z.string().nonempty("Required"),
  })),
  price: z.number().optional()
})

export default function FlightDialogContent({
  trip, flight
}: {
  trip: Trip
  flight?: Transportation | null
}) {
  const [edit] = useState<boolean>(flight == null)
  const {onClose} = useDialogContext()

  function mapLegsOrDefault(flightLegs: FlightLeg[]|undefined) {
    if (flightLegs != null) {
      return flightLegs.map(leg => ({
        date: dateFromString(trip.startDate),
        flightNumber: leg.flightNumber,
      }))
    }

    return [{
      date: dateFromString(trip.startDate),
      flightNumber: "",
    }]
  }

  const form = useForm<z.input<typeof formSchema>, unknown, z.output<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      legs: mapLegsOrDefault(flight?.flightDetail?.legs),
      pnrs: flight?.flightDetail?.pnrs ?? [],
      price: flight?.price ?? undefined
    },
    disabled: !edit
  })
  const { isSubmitting } = form.formState;

  const legsArray = useFieldArray({
    control: form.control,
    name: "legs",
    rules: {
      minLength: 1
    }
  });

  const pnrsArray = useFieldArray({
    control: form.control,
    name: "pnrs"
  });

  async function onSubmit(values: z.infer<typeof formSchema>) {
    let response
    if (flight != null) {
      response = await fetch("/api/v1/trips/" + trip.id + "/flights/" + flight?.id, {
        method: "PUT",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(values)
      })
    } else {
      response = await fetch("/api/v1/trips/" + trip.id + "/flights", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(values)
      })
    }

    if (response.ok)
      onClose(true)
    else toast("Error upserting Flight", {
      description: await response.text()
    })
  }

  async function onDeleteButtonClick() {
    const response = await fetch("/api/v1/trips/" + trip.id + "/transportation/" + flight!.id, {method: "DELETE"})

    if (response.ok)
      onClose(true)
    else toast("Error deleting Flight", {
      description: await response.text()
    })
  }

  async function onUpdateButtonClick() {
    const response = await fetch("/api/v1/trips/" + trip.id + "/flights/" + flight!.id, {method: "PUT"})

    if (response.ok)
      onClose(true)
    else toast("Error updating Flight", {
      description: await response.text()
    })
  }

  function addLeg() {
    legsArray.append({
      date: legsArray.fields[legsArray.fields.length-1].date,
      flightNumber: ""
    })
  }

  function deleteLeg() {
    legsArray.remove(legsArray.fields.length-1)
  }

  function addPnr() {
    pnrsArray.append({airline: "", pnr: ""})
  }

  function deletePnr() {
    pnrsArray.remove(pnrsArray.fields.length-1)
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
          {edit &&
            <div>
              {legsArray.fields.length > 1 ?
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

        {legsArray.fields.map((field, index) =>
            <div key={field.id}>
              <RowContainer>
                <FormField control={form.control} name={`legs.${index}.date`} label={`Date ${legsArray.fields.length > 1 ? (index+1) : ""}`}
                           render={({field}) =>
                               <DateInput startDate={trip.startDate} endDate={trip.endDate} {...field}/>
                           }
                />
                <FormField control={form.control} name={`legs.${index}.flightNumber`}
                           label={`Flight #${legsArray.fields.length > 1 ? (index+1) : ""}`}
                           render={({field}) =>
                               <Input placeholder="LH717" {...field}/>
                           }
                />
              </RowContainer>
            </div>
        )}

        <div className="flex">
          <h3 className="font-semibold mb-2 grow">PNRs</h3>
          {edit &&
            <div>
              {pnrsArray.fields.length > 0 ?
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

        {pnrsArray.fields.map((field, index) =>
          <div key={field.id}>
            <RowContainer>
              <FormField control={form.control} name={`pnrs.${index}.airline`}
                         label={`Airline ${pnrsArray.fields.length > 1 ? (index+1) : ""}`}
                         render={({field}) =>
                             <Input placeholder="LH" {...field}/>
                         }
              />
              <FormField control={form.control} name={`pnrs.${index}.pnr`}
                         label={`PNR ${pnrsArray.fields.length > 1 ? (index+1) : ""}`}
                         render={({field}) =>
                             <Input placeholder="123ABC" {...field}/>
                         }
              />
            </RowContainer>
          </div>
        )}

        <FormField control={form.control} name="price" label="Price"
                   render={({field}) =>
                       <AmountInput {...field}/>
                   }
        />
      </Form>
      <DialogFooter>
        {edit ?
          <Button form="flight-form" type="submit" className="w-full text-base" disabled={isSubmitting}>
            {isSubmitting ? <Spinner variant="pinwheel"/> : "Save"}
          </Button>
        :
          <>
            <Button variant="destructive" className="w-full text-base" onClick={onDeleteButtonClick}>
              Delete
            </Button>
            <Button variant="secondary" className="w-full text-base" onClick={onUpdateButtonClick}>
              Refresh Data
            </Button>
          </>
        }
      </DialogFooter>
    </>
  )
}