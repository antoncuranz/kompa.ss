import {Button} from "@/components/ui/button.tsx";
import {DialogFooter, DialogHeader, DialogTitle} from "@/components/ui/dialog.tsx";
import {Minus, Plus} from "lucide-react";
import {useState} from "react";
import {Input} from "@/components/ui/input.tsx";
import {TrainLeg, Transportation, Trip} from "@/types.ts";
import {RowContainer, useDialogContext} from "@/components/dialog/Dialog.tsx";
import {toast} from "sonner";
import {z} from "zod";
import {isoDate, trainStation} from "@/schemas.ts";
import {useFieldArray, useForm} from "react-hook-form";
import {zodResolver} from "@hookform/resolvers/zod";
import {Form, FormField} from "@/components/ui/form";
import DateInput from "@/components/dialog/DateInput.tsx";
import AmountInput from "@/components/dialog/AmountInput.tsx";
import TrainStationInput from "@/components/dialog/TrainStationInput.tsx";
import {dateFromString} from "@/components/util.ts";
import {Spinner} from "@/components/ui/shadcn-io/spinner";

const formSchema = z.object({
  departureDate: isoDate("Required"),
  fromStationId: trainStation("Required"),
  toStationId: trainStation("Required"),
  viaStationId: trainStation().optional(),
  trainNumbers: z.array(z.object({
    value: z.string().nonempty("Required")
  })).transform(x => x.map(y => y.value)),
  price: z.number().optional()
})

export default function TrainDialogContent({
  trip, train
}: {
  trip: Trip
  train?: Transportation | null
}) {
  const [edit, setEdit] = useState<boolean>(train == null)
  const {onClose} = useDialogContext()

  const form = useForm<z.input<typeof formSchema>, unknown, z.output<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      departureDate: train ? dateFromString(train.departureDateTime) : undefined,
      fromStationId: getDefaultFromStation(train?.trainDetail?.legs),
      toStationId: getDefaultToStation(train?.trainDetail?.legs),
      viaStationId: undefined,
      trainNumbers: mapLegsOrDefault(train?.trainDetail?.legs),
      price: train?.price ?? undefined
    },
    disabled: !edit
  })
  const { isSubmitting } = form.formState;

  function getDefaultFromStation(trainLegs: TrainLeg[]|undefined) {
    if (trainLegs != null) {
      return trainLegs[0].origin
    }
    return undefined
  }

  function getDefaultToStation(trainLegs: TrainLeg[]|undefined) {
    if (trainLegs != null) {
      return trainLegs[trainLegs.length-1].destination
    }
    return undefined
  }

  function mapLegsOrDefault(trainLegs: TrainLeg[]|undefined) {
    if (trainLegs != null) {
      return trainLegs.map(leg => ({
        value: leg.lineName,
      }))
    }
    return [{value: ""}]
  }

  const trainNumbersArray = useFieldArray({
    control: form.control,
    name: "trainNumbers",
    rules: {
      minLength: 1
    }
  });

  async function onSubmit(values: z.infer<typeof formSchema>) {
    let response
    if (train != null) {
      response = await fetch("/api/v1/trips/" + trip.id + "/trains/" + train?.id, {
        method: "PUT",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(values)
      })
    } else {
      response = await fetch("/api/v1/trips/" + trip.id + "/trains", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(values)
      })
    }

    if (response.ok)
      onClose(true)
    else toast("Error upserting Train", {
      description: await response.text()
    })
  }

  async function onDeleteButtonClick() {
    const response = await fetch("/api/v1/trips/" + trip.id + "/transportation/" + train!.id, {method: "DELETE"})

    if (response.ok)
      onClose(true)
    else toast("Error deleting Train", {
      description: await response.text()
    })
  }

  function addLeg() {
    trainNumbersArray.append({value: ""})
  }

  function deleteLeg() {
    trainNumbersArray.remove(trainNumbersArray.fields.length-1)
  }

  return (
    <>
      <DialogHeader>
        <DialogTitle>
          {edit ? (train != null ? "Edit" : "New") : "View"} Train Journey
        </DialogTitle>
      </DialogHeader>
      <Form id="train-form" form={form} onSubmit={form.handleSubmit(onSubmit)}>
        <FormField control={form.control} name="fromStationId" label="From Station"
                   render={({field}) =>
                       <TrainStationInput placeholder="Berlin Südkreuz" {...field}/>
                   }
        />
        <FormField control={form.control} name="viaStationId" label="Via Station"
                   render={({field}) =>
                       <TrainStationInput placeholder="Optional" {...field}/>
                   }
        />
        <FormField control={form.control} name="toStationId" label="To Station"
                   render={({field}) =>
                       <TrainStationInput placeholder="München" {...field}/>
                   }
        />
        <RowContainer>
          <FormField control={form.control} name="departureDate" label="Date"
                     render={({field}) =>
                         <DateInput startDate={trip.startDate} endDate={trip.endDate} {...field}/>
                     }
          />
          <FormField control={form.control} name="price" label="Price"
                     render={({field}) =>
                         <AmountInput {...field}/>
                     }
          />
        </RowContainer>

        <div className="flex">
          <h3 className="font-semibold mb-2 flex-grow">Train Legs</h3>
          {edit &&
            <div>
              {trainNumbersArray.fields.length > 1 ?
                <Button type="button" variant="ghost" className="p-2 h-auto rounded-full" onClick={() => deleteLeg()}>
                  <Minus className="h-4 w-4"/>
                </Button>
                :
                <div/>
              }
              <Button type="button" variant="ghost" className="p-2 h-auto rounded-full" onClick={() => addLeg()}>
                <Plus className="w-4 h-4"/>
              </Button>
            </div>
          }
        </div>

        {trainNumbersArray.fields.map((field, index) =>
            <FormField key={field.id} control={form.control} name={`trainNumbers.${index}.value`}
                       label={`Train #${trainNumbersArray.fields.length > 1 ? (index+1) : ""}`}
                       render={({field}) =>
                           <Input placeholder="ICE707" {...field}/>
                       }
            />
        )}
      </Form>
      <DialogFooter>
        {edit ?
          <Button form="train-form" type="submit" className="w-full text-base" disabled={isSubmitting}>
            {isSubmitting ? <Spinner variant="pinwheel"/> : "Save"}
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
    </>
  )
}