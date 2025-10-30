import {Button} from "@/components/ui/button.tsx";
import {DialogFooter, DialogHeader, DialogTitle} from "@/components/ui/dialog.tsx";
import {useState} from "react";
import {Input} from "@/components/ui/input.tsx";
import {Textarea} from "@/components/ui/textarea.tsx";

import {Form, FormField} from "@/components/ui/form"
import {z} from "zod"
import {useForm} from "react-hook-form";
import {zodResolver} from "@hookform/resolvers/zod";
import DateInput from "@/components/dialog/input/DateInput.tsx";
import AmountInput from "@/components/dialog/input/AmountInput.tsx";
import AddressInput from "@/components/dialog/input/AddressInput.tsx";
import {dateFromString} from "@/components/util.ts";
import {RowContainer, useDialogContext} from "@/components/dialog/Dialog.tsx";
import {Activity, isoDate, optionalLocation, optionalString, Trip} from "@/schema";
import LocationInput from "@/components/dialog/input/LocationInput.tsx";
import {Spinner} from "@/components/ui/shadcn-io/spinner";

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
  activity?: Activity
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
      location: activity?.location ?? undefined,
    },
    disabled: !edit
  })
  const { isSubmitting } = form.formState;

  async function onSubmit(values: z.infer<typeof formSchema>) {
    if (activity) {
      activity.$jazz.applyDiff(values)
      if (!activity.location) {
        activity.$jazz.set("location", values.location)
      }
    } else {
      trip.activities.$jazz.push(values)
    }
    onClose()
  }

  async function onDeleteButtonClick() {
    if (activity === undefined) {
      return
    }

    trip.activities.$jazz.remove(a => a && a.$jazz.id == activity.$jazz.id)
    onClose()
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
                       <Input data-1p-ignore placeholder="My new Activity" {...field} />
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
          <Button form="activity-form" type="submit" className="w-full" disabled={isSubmitting}>
            {isSubmitting ? <Spinner variant="pinwheel"/> : "Save"}
          </Button>
        :
          <>
            <Button variant="destructive" className="w-full" onClick={onDeleteButtonClick}>
              Delete
            </Button>
            <Button variant="secondary" className="w-full" onClick={() => setEdit(true)}>
              Edit
            </Button>
          </>
        }
      </DialogFooter>
    </>
  )
}