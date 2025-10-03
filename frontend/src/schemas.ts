import * as core from "zod/v4/core";
import {z} from "zod"
import {dateToString} from "@/components/util.ts";

export function optionalString(params?: string | core.$ZodStringParams) {
  return z.string(params).trim().transform(x => x || undefined).optional()
}

export function isoDate(params?: string | core.$ZodDateParams) {
  return z.date(params).transform(dateToString)
}

export function optionalNumberString(params?: string | core.$ZodStringParams) {
  return z.string(params).transform(x => x ? Number(x) : undefined)
}

export function optionalLocation() {
  return z.object({
    latitude: optionalNumberString(),
    longitude: optionalNumberString()
  }).transform(x => (x.latitude && x.longitude) ? x : undefined)
}

export function trainStation(params?: string | core.$ZodObjectParams) {
  return z.object({
    id: z.string().nonempty(),
    name: z.string().nonempty(),
  }, params).transform(station => station.id)
}

