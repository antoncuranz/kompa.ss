import {Button} from "@/components/ui/button.tsx";
import {DialogFooter, DialogHeader, DialogTitle} from "@/components/ui/dialog.tsx";
import {useState} from "react";
import {Input} from "@/components/ui/input.tsx";
import {getTransportationTypeEmoji, Transportation, TransportationType, Trip} from "@/types.ts";
import {Form, FormField} from "@/components/ui/form"
import {z} from "zod"
import {useForm} from "react-hook-form";
import {zodResolver} from "@hookform/resolvers/zod";
import DateInput from "@/components/dialog/DateInput.tsx";
import AmountInput from "@/components/dialog/AmountInput.tsx";
import AddressInput from "@/components/dialog/AddressInput.tsx";
import {dateFromString, titleCase} from "@/components/util.ts";
import {RowContainer, useDialogContext} from "@/components/dialog/Dialog.tsx";
import {isoDateTime, location, optionalString} from "@/schemas";
import {toast} from "sonner";
import LocationInput from "@/components/dialog/LocationInput.tsx";
import {Spinner} from "@/components/ui/shadcn-io/spinner";
import {Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue} from "@/components/ui/select"
import {Separator} from "@/components/ui/separator.tsx";

const formSchema = z.object({
  name: z.string().nonempty("Required"),
  type: z.string().nonempty("Required"),
  price: z.number().optional(),
  departureDateTime: isoDateTime("Required"),
  arrivalDateTime: isoDateTime("Required"),
  origin: location("Required"),
  destination: location("Required"),
  originAddress: optionalString(),
  destinationAddress: optionalString()
})

export default function TransportationDialogContent({
  trip, transportation
}: {
  trip: Trip
  transportation?: Transportation | null
}) {
  const [edit, setEdit] = useState<boolean>(transportation == null)
  const {onClose} = useDialogContext()

  const form = useForm<z.input<typeof formSchema>, unknown, z.output<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "", //transportation?.name ?? "",
      type: transportation?.type ?? "",
      price: transportation?.price ?? undefined,
      departureDateTime: transportation?.departureDateTime ? dateFromString(transportation.departureDateTime) : undefined,
      arrivalDateTime: transportation?.arrivalDateTime ? dateFromString(transportation.arrivalDateTime) : undefined,
      origin: transportation?.origin ?? undefined,
      destination: transportation?.destination ?? undefined,
      originAddress: "", //transportation?.address ?? "",
      destinationAddress: "", //transportation?.address ?? "",
    },
    disabled: !edit
  })
  const { isSubmitting } = form.formState;

  async function onSubmit(values: z.infer<typeof formSchema>) {
    console.log(values)
    console.log(JSON.stringify(values))

    let response
    if (transportation != null) {
      response = await fetch("/api/v1/trips/" + trip.id + "/transportation/" + transportation?.id, {
        method: "PUT",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(values)
      })
    } else {
      response = await fetch("/api/v1/trips/" + trip.id + "/transportation", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(values)
      })
    }

    if (response.ok)
      onClose(true)
    else toast("Error upserting Transportation", {
      description: await response.text()
    })
  }

  async function onDeleteButtonClick() {
    const response = await fetch("/api/v1/trips/" + trip.id + "/activities/" + transportation!.id, {method: "DELETE"})

    if (response.ok)
      onClose(true)
    else toast("Error deleting Transportation", {
      description: await response.text()
    })
  }

  return (
    <>
      <DialogHeader>
        <DialogTitle>
          {edit ? (transportation != null ? "Edit" : "New") : "View"} Transportation
        </DialogTitle>
      </DialogHeader>
      <Form id="transportation-form" form={form} onSubmit={form.handleSubmit(onSubmit)}>
        <FormField control={form.control} name="name" label="Name"
                   render={({field}) =>
                       <Input data-1p-ignore placeholder="My new Transportation" {...field} />
                   }
        />
        <RowContainer>
          <FormField control={form.control} name="type" label="Type"
                     render={({field}) =>
                         <Select name={field.name} onValueChange={field.onChange} defaultValue={field.value}>
                           <SelectTrigger>
                             <SelectValue placeholder="Select type"/>
                           </SelectTrigger>
                           <SelectContent>
                             <SelectGroup>
                               {Object.values(TransportationType).map(type =>
                                 <SelectItem key={type} value={type}>
                                   {getTransportationTypeEmoji(type)} {titleCase(type)}
                                 </SelectItem>
                               )}
                             </SelectGroup>
                           </SelectContent>
                         </Select>
                     }
          />
          <FormField control={form.control} name="price" label="Price"
                     render={({field}) =>
                         <AmountInput {...field}/>
                     }
          />
        </RowContainer>
        <Separator className="mt-4 mb-2"/>
        <FormField control={form.control} name="departureDateTime" label="Departure Time"
                   render={({field}) =>
                       <DateInput startDate={trip.startDate} endDate={trip.endDate} {...field}/>
                   }
        />
        <FormField control={form.control} name="originAddress" label="Origin Address"
                   render={({field}) =>
                       <AddressInput
                           {...field}
                           updateCoordinates={coords => form.setValue("origin", coords)}
                       />
                   }
        />
        <FormField control={form.control} name="origin" label="Origin Coordinates"
                   render={({field}) =>
                       <LocationInput {...field}/>
                   }
        />
        <Separator className="mt-4 mb-2"/>
        <FormField control={form.control} name="arrivalDateTime" label="Arrival Date"
                   render={({field}) =>
                       <DateInput startDate={trip.startDate} endDate={trip.endDate} {...field}/>
                   }
        />
        <FormField control={form.control} name="destinationAddress" label="Destination Address"
                   render={({field}) =>
                       <AddressInput
                           {...field}
                           updateCoordinates={coords => form.setValue("destination", coords)}
                       />
                   }
        />
        <FormField control={form.control} name="destination" label="Destination Coordinates"
                   render={({field}) =>
                       <LocationInput {...field}/>
                   }
        />
      </Form>
      <DialogFooter>
        {edit ?
          <Button form="transportation-form" type="submit" className="w-full text-base" disabled={isSubmitting}>
            {isSubmitting ? <Spinner variant="pinwheel"/> : "Save"}
          </Button>
        :
          <>
            <Button variant="destructive" className="w-full text-base" onClick={onDeleteButtonClick}>
              Delete
            </Button>
            <Button variant="secondary" className="w-full text-base" onClick={() => setEdit(true)}>
              Edit
            </Button>
          </>
        }
      </DialogFooter>
    </>
  )
}