import {Button} from "@/components/ui/button.tsx";
import {DialogFooter, DialogHeader, DialogTitle} from "@/components/ui/dialog.tsx";
import {useState} from "react";
import {Input} from "@/components/ui/input.tsx";
import {Activity, Trip} from "@/types.ts";
import {Textarea} from "@/components/ui/textarea.tsx";

import {Form, FormField} from "@/components/ui/form"
import {z} from "zod"
import {useForm} from "react-hook-form";
import {zodResolver} from "@hookform/resolvers/zod";
import DateInput from "@/components/dialog/DateInput.tsx";
import AmountInput from "@/components/dialog/AmountInput.tsx";
import AddressInput from "@/components/dialog/AddressInput.tsx";
import {dateFromString} from "@/components/util.ts";
import {RowContainer, useDialogContext} from "@/components/dialog/Dialog.tsx";
import {isoDate, optionalLocation, optionalString} from "@/schemas";
import {toast} from "sonner";

const formSchema = z.object({
  name: z.string().nonempty("Required"),
  description: optionalString(),
  date: isoDate("Required"),
  price: z.number().optional(),
  address: optionalString(),
  location: optionalLocation()
})

export default function ActivityDialogContent({
  trip, activity
}: {
  trip: Trip
  activity?: Activity | null
}) {
  const [edit, setEdit] = useState<boolean>(activity == null)
  const {onClose} = useDialogContext()

  const form = useForm<z.input<typeof formSchema>, unknown, z.output<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: activity?.name ?? "",
      description: activity?.description ?? "",
      date: activity?.date ? dateFromString(activity.date) : undefined,
      price: activity?.price ?? undefined,
      address: activity?.address ?? "",
      location: {
        latitude: activity?.location?.latitude.toString() ?? "",
        longitude: activity?.location?.longitude.toString() ?? ""
      }
    },
    disabled: !edit
  })

  async function onSubmit(values: z.infer<typeof formSchema>) {
    let response
    if (activity != null) {
      response = await fetch("/api/v1/trips/" + trip.id + "/activities/" + activity?.id, {
        method: "PUT",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(values)
      })
    } else {
      response = await fetch("/api/v1/trips/" + trip.id + "/activities", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(values)
      })
    }

    if (response.ok)
      onClose(true)
    else toast("Error upserting Activity", {
      description: await response.text()
    })
  }

  async function onDeleteButtonClick() {
    const response = await fetch("/api/v1/trips/" + trip.id + "/activities/" + activity!.id, {method: "DELETE"})

    if (response.ok)
      onClose(true)
    else toast("Error deleting Activity", {
      description: await response.text()
    })
  }

  return (
    <>
      <DialogHeader>
        <DialogTitle>
          {edit ? (activity != null ? "Edit" : "New") : "View"} Activity
        </DialogTitle>
      </DialogHeader>
      <Form id="activity-form" form={form} onSubmit={form.handleSubmit(onSubmit)}>
        <FormField control={form.control} name="name" label="Name"
                   render={({field}) =>
                       <Input placeholder="My new Activity" {...field} />
                   }
        />
        <FormField control={form.control} name="description" label="Description"
                   render={({field}) =>
                       <Textarea id="description" {...field}/>
                   }
        />
        <RowContainer>
          <FormField control={form.control} name="date" label="Date"
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
        <FormField control={form.control} name="address" label="Address"
                   render={({field}) =>
                       <AddressInput
                           {...field}
                           updateLatitude={lat => form.setValue("location.latitude", String(lat))}
                           updateLongitude={lon => form.setValue("location.longitude", String(lon))}
                       />
                   }
        />
        <RowContainer>
          <FormField control={form.control} name="location.latitude" label="Latitude"
                     render={({field}) =>
                         <Input {...field} readOnly={true}/>
                     }
          />
          <FormField control={form.control} name="location.longitude" label="Longitude"
                     render={({field}) =>
                         <Input {...field} readOnly={true}/>
                     }
          />
        </RowContainer>
      </Form>
      <DialogFooter>
        {edit ?
          <Button form="activity-form" type="submit" className="w-full text-base">
            Save
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