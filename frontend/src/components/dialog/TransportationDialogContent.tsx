import { Button } from "@/components/ui/button.tsx"
import {
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog.tsx"
import { useState } from "react"
import { Input } from "@/components/ui/input.tsx"
import { getTransportationTypeEmoji, TransportationType } from "@/types.ts"
import { Form, FormField } from "@/components/ui/form"
import { z } from "zod"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import AmountInput from "@/components/dialog/input/AmountInput.tsx"
import AddressInput from "@/components/dialog/input/AddressInput.tsx"
import { dateFromString, titleCase } from "@/components/util.ts"
import { RowContainer, useDialogContext } from "@/components/dialog/Dialog.tsx"
import { GenericTransportation, Trip } from "@/schema"
import { isoDateTime, location, optionalString } from "@/formschema"
import LocationInput from "@/components/dialog/input/LocationInput.tsx"
import { Spinner } from "@/components/ui/shadcn-io/spinner"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectPositioner,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Separator } from "@/components/ui/separator.tsx"
import DateTimeInput from "@/components/dialog/input/DateTimeInput.tsx"

const formSchema = z.object({
  name: z.string().nonempty("Required"),
  genericType: z.string().nonempty("Required"),
  price: z.number().optional(),
  departureDateTime: isoDateTime("Required"),
  arrivalDateTime: isoDateTime("Required"),
  origin: location("Required"),
  destination: location("Required"),
  originAddress: optionalString(),
  destinationAddress: optionalString(),
})

export default function TransportationDialogContent({
  trip,
  transportation,
}: {
  trip: Trip
  transportation?: GenericTransportation
}) {
  const [edit, setEdit] = useState<boolean>(transportation == null)
  const { onClose } = useDialogContext()

  const form = useForm<
    z.input<typeof formSchema>,
    unknown,
    z.output<typeof formSchema>
  >({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: transportation?.name ?? "",
      genericType: transportation?.genericType ?? "",
      price: transportation?.price ?? undefined,
      departureDateTime: transportation?.departureDateTime
        ? dateFromString(transportation.departureDateTime)
        : undefined,
      arrivalDateTime: transportation?.arrivalDateTime
        ? dateFromString(transportation.arrivalDateTime)
        : undefined,
      origin: transportation?.origin ?? undefined,
      destination: transportation?.destination ?? undefined,
      originAddress: transportation?.originAddress ?? "",
      destinationAddress: transportation?.destinationAddress ?? "",
    },
    disabled: !edit,
  })
  const { isSubmitting } = form.formState

  async function onSubmit(values: z.infer<typeof formSchema>) {
    if (transportation) {
      transportation.$jazz.applyDiff(values)
      if (!transportation.origin) {
        transportation.$jazz.set("origin", values.origin)
      }
      if (!transportation.destination) {
        transportation.$jazz.set("destination", values.destination)
      }
    } else {
      trip.transportation.$jazz.push({
        type: "generic",
        ...values,
      })
    }
    onClose()
  }

  async function onDeleteButtonClick() {
    if (transportation === undefined) {
      return
    }

    trip.transportation.$jazz.remove(
      a => a?.$jazz.id == transportation.$jazz.id,
    )
    onClose()
  }

  return (
    <>
      <DialogHeader>
        <DialogTitle>
          {edit ? (transportation != null ? "Edit" : "New") : "View"}{" "}
          Transportation
        </DialogTitle>
      </DialogHeader>
      <Form
        id="transportation-form"
        form={form}
        onSubmit={form.handleSubmit(onSubmit)}
      >
        <FormField
          control={form.control}
          name="name"
          label="Name"
          render={({ field }) => (
            <Input
              data-1p-ignore
              placeholder="My new Transportation"
              {...field}
            />
          )}
        />
        <RowContainer>
          <FormField
            control={form.control}
            name="genericType"
            label="Type"
            render={({ field }) => (
              <Select
                name={field.name}
                onValueChange={field.onChange}
                value={field.value ?? ""}
                disabled={field.disabled}
              >
                <SelectTrigger className="w-full">
                  <SelectValue placeholder="Select type" />
                </SelectTrigger>
                <SelectPositioner>
                  <SelectContent>
                    {[
                      TransportationType.Bus,
                      TransportationType.Ferry,
                      TransportationType.Boat,
                      TransportationType.Bike,
                      TransportationType.Car,
                      TransportationType.Hike,
                      TransportationType.Other,
                    ].map(type => (
                      <SelectItem key={type} value={type}>
                        {getTransportationTypeEmoji(type)} {titleCase(type)}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </SelectPositioner>
              </Select>
            )}
          />
          <FormField
            control={form.control}
            name="price"
            label="Price"
            render={({ field }) => <AmountInput {...field} />}
          />
        </RowContainer>
        <Separator className="mt-4 mb-2" />
        <FormField
          control={form.control}
          name="departureDateTime"
          label="Departure Time"
          render={({ field }) => (
            <DateTimeInput
              startDate={trip.startDate}
              endDate={trip.endDate}
              {...field}
            />
          )}
        />
        <FormField
          control={form.control}
          name="originAddress"
          label="Origin Address"
          render={({ field }) => (
            <AddressInput
              {...field}
              updateCoordinates={coords => form.setValue("origin", coords)}
            />
          )}
        />
        <FormField
          control={form.control}
          name="origin"
          label="Origin Coordinates"
          render={({ field }) => <LocationInput {...field} />}
        />
        <Separator className="mt-4 mb-2" />
        <FormField
          control={form.control}
          name="arrivalDateTime"
          label="Arrival Time"
          render={({ field }) => (
            <DateTimeInput
              startDate={trip.startDate}
              endDate={trip.endDate}
              {...field}
            />
          )}
        />
        <FormField
          control={form.control}
          name="destinationAddress"
          label="Destination Address"
          render={({ field }) => (
            <AddressInput
              {...field}
              updateCoordinates={coords => form.setValue("destination", coords)}
            />
          )}
        />
        <FormField
          control={form.control}
          name="destination"
          label="Destination Coordinates"
          render={({ field }) => <LocationInput {...field} />}
        />
      </Form>
      <DialogFooter>
        {edit ? (
          <Button
            form="transportation-form"
            type="submit"
            className="w-full"
            disabled={isSubmitting}
          >
            {isSubmitting ? <Spinner variant="pinwheel" /> : "Save"}
          </Button>
        ) : (
          <>
            <Button
              variant="destructive"
              className="w-full"
              onClick={onDeleteButtonClick}
            >
              Delete
            </Button>
            <Button
              variant="secondary"
              className="w-full"
              onClick={() => setEdit(true)}
            >
              Edit
            </Button>
          </>
        )}
      </DialogFooter>
    </>
  )
}
