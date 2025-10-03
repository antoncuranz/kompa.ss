import {Button} from "@/components/ui/button.tsx";
import {DialogFooter, DialogHeader, DialogTitle} from "@/components/ui/dialog.tsx";
import {useState} from "react";
import {Input} from "@/components/ui/input.tsx";
import {Trip} from "@/types.ts";
import {Textarea} from "@/components/ui/textarea.tsx";
import {dateFromString} from "@/components/util.ts";
import {RowContainer, useDialogContext} from "@/components/dialog/Dialog.tsx";
import {toast} from "sonner";
import {Form, FormField} from "@/components/ui/form.tsx";
import {z} from "zod"
import {isoDate, optionalString} from "@/schemas.ts";
import {useForm} from "react-hook-form";
import {zodResolver} from "@hookform/resolvers/zod";
import DateInput from "@/components/dialog/DateInput.tsx";

const formSchema = z.object({
  name: z.string().nonempty("Required"),
  description: optionalString(),
  startDate: isoDate("Required"),
  endDate: isoDate("Required"),
  imageUrl: optionalString()
})

export default function TripDialogContent({
  trip
}: {
  trip?: Trip | null
}) {
  const [edit, setEdit] = useState<boolean>(trip == null)
  const {onClose} = useDialogContext()

  const form = useForm<z.input<typeof formSchema>, unknown, z.output<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: trip?.name ?? "",
      description: trip?.description ?? "",
      startDate: trip?.startDate ? dateFromString(trip.startDate) : undefined,
      endDate: trip?.endDate ? dateFromString(trip.endDate) : undefined,
      imageUrl: trip?.imageUrl ?? ""
    },
    disabled: !edit
  })

  async function onSubmit(values: z.infer<typeof formSchema>) {
    let response
    if (trip != null) {
      response = await fetch("/api/v1/trips/" + trip.id, {
        method: "PUT",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(values)
      })
    } else {
      response = await fetch("/api/v1/trips", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(values)
      })
    }

    if (response.ok)
      onClose(true)
    else toast("Error upserting Trip", {
      description: await response.text()
    })
  }

  async function onDeleteButtonClick() {
    const response = await fetch("/api/v1/trips/" + trip!.id, {method: "DELETE"})

    if (response.ok)
      onClose(true)
    else toast("Error deleting Trip", {
      description: await response.text()
    })
  }

  return (
    <>
      <DialogHeader>
        <DialogTitle>
          {edit ? (trip != null ? "Edit" : "New") : "View"} Trip
        </DialogTitle>
      </DialogHeader>
      <Form id="trip-form" form={form} onSubmit={form.handleSubmit(onSubmit)}>
        <FormField control={form.control} name="name" label="Name"
                   render={({field}) =>
                       <Input placeholder="My new Activity" {...field} />
                   }
        />
        <RowContainer>
          <FormField control={form.control} name="startDate" label="Start Date"
                     render={({field}) =>
                         <DateInput {...field}/>
                     }
          />
          <FormField control={form.control} name="endDate" label="End Date"
                     render={({field}) =>
                         <DateInput {...field}/>
                     }
          />
        </RowContainer>
        <FormField control={form.control} name="description" label="Description"
                   render={({field}) =>
                       <Textarea id="description" {...field}/>
                   }
        />
        <FormField control={form.control} name="imageUrl" label="Image URL"
                   render={({field}) =>
                       <Input {...field}/>
                   }
        />
      </Form>
      <DialogFooter>
        {edit ?
          <Button form="trip-form" type="submit" className="w-full text-base">
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