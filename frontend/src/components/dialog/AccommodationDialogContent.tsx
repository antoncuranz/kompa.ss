import {Button} from "@/components/ui/button.tsx";
import {DialogFooter, DialogHeader, DialogTitle} from "@/components/ui/dialog.tsx";
import {useState} from "react";
import {Input} from "@/components/ui/input.tsx";
import {Accommodation, Trip} from "@/types.ts";
import {Textarea} from "@/components/ui/textarea.tsx";
import {RowContainer, useDialogContext} from "@/components/dialog/Dialog.tsx";
import {toast} from "sonner";
import {z} from "zod"
import {isoDate, optionalLocation, optionalString} from "@/schemas.ts";
import {Form, FormField} from "@/components/ui/form"
import {useForm} from "react-hook-form";
import {zodResolver} from "@hookform/resolvers/zod";
import {dateFromString} from "@/components/util.ts";
import AmountInput from "@/components/dialog/AmountInput.tsx";
import AddressInput from "@/components/dialog/AddressInput.tsx";
import DateInput from "@/components/dialog/DateInput.tsx";
import LocationInput from "@/components/dialog/LocationInput.tsx";
import {Spinner} from "@/components/ui/shadcn-io/spinner";

const formSchema = z.object({
  name: z.string().nonempty("Required"),
  description: optionalString(),
  arrivalDate: isoDate("Required"),
  departureDate: isoDate("Required"),
  price: z.number().optional(),
  address: optionalString(),
  location: optionalLocation()
})

export default function AccommodationDialogContent({
  trip, accommodation
}: {
  trip: Trip,
  accommodation?: Accommodation | null
}) {
  const [edit, setEdit] = useState<boolean>(accommodation == null)
  const {onClose} = useDialogContext()

  const form = useForm<z.input<typeof formSchema>, unknown, z.output<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: accommodation?.name ?? "",
      description: accommodation?.description ?? "",
      arrivalDate: accommodation?.arrivalDate ? dateFromString(accommodation.arrivalDate) : undefined,
      departureDate: accommodation?.departureDate ? dateFromString(accommodation.departureDate) : undefined,
      price: accommodation?.price ?? undefined,
      address: accommodation?.address ?? "",
      location: accommodation?.location ?? undefined
    },
    disabled: !edit
  })
  const { isSubmitting } = form.formState;

  async function onSubmit(values: z.infer<typeof formSchema>) {
    let response
    if (accommodation != null) {
      response = await fetch("/api/v1/trips/" + trip.id + "/accommodation/" + accommodation?.id, {
        method: "PUT",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(values)
      })
    } else {
      response = await fetch("/api/v1/trips/" + trip.id + "/accommodation", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(values)
      })
    }

    if (response.ok)
      onClose(true)
    else toast("Error upserting Accommodation", {
      description: await response.text()
    })
  }

  async function onDeleteButtonClick() {
    const response = await fetch("/api/v1/trips/" + trip.id + "/accommodation/" + accommodation!.id, {method: "DELETE"})

    if (response.ok)
      onClose(true)
    else toast("Error deleting Accommodation", {
      description: await response.text()
    })
  }

  return (
    <>
      <DialogHeader>
        <DialogTitle>
          {edit ? (accommodation != null ? "Edit" : "New") : "View"} Accommodation
        </DialogTitle>
      </DialogHeader>
      <Form id="accommodation-form" form={form} onSubmit={form.handleSubmit(onSubmit)}>
        <FormField control={form.control} name="name" label="Name"
                   render={({field}) =>
                       <Input data-1p-ignore {...field} />
                   }
        />
        <RowContainer>
          <FormField control={form.control} name="arrivalDate" label="Arrival Date"
                     render={({field}) =>
                         <DateInput startDate={trip.startDate} endDate={trip.endDate} {...field}/>
                     }
          />
          <FormField control={form.control} name="departureDate" label="Departure Date"
                     render={({field}) =>
                         <DateInput startDate={trip.startDate} endDate={trip.endDate} {...field}/>
                     }
          />
        </RowContainer>
        <FormField control={form.control} name="description" label="Description"
                   render={({field}) =>
                       <Textarea id="description" {...field}/>
                   }
        />
        <FormField control={form.control} name="price" label="Price"
                   render={({field}) =>
                       <AmountInput {...field}/>
                   }
        />
        <FormField control={form.control} name="address" label="Address"
                   render={({field}) =>
                       <AddressInput
                           {...field}
                           updateCoordinates={coords => form.setValue("location", coords)}
                       />
                   }
        />
        <FormField control={form.control} name="location" label="Coordinates"
                   render={({field}) =>
                       <LocationInput {...field}/>
                   }
        />
      </Form>
      <DialogFooter>
        {edit ?
          <Button form="accommodation-form" type="submit" className="w-full text-base" disabled={isSubmitting}>
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