import * as core from "zod/v4/core";
import {dateTimeToString, dateToString} from "@/components/util.ts";
import {co, z} from "jazz-tools";

export function optionalString(params?: string | core.$ZodStringParams) {
  return z.string(params).trim().transform(x => x || undefined).optional()
}

export function isoDate(params?: string | core.$ZodDateParams) {
  return z.date(params).transform(dateToString)
}

export function isoDateTime(params?: string | core.$ZodDateParams) {
  return z.date(params).transform(dateTimeToString)
}

export function optionalNumberString(params?: string | core.$ZodStringParams) {
  return z.string(params).transform(x => x ? Number(x) : undefined)
}

export function location(params?: string | core.$ZodObjectParams) {
  return z.object({
    latitude: z.number(),
    longitude: z.number()
  }, params).transform(x => (x.latitude && x.longitude) ? x : undefined)
}

export function trainStation(params?: string | core.$ZodObjectParams) {
  return z.object({
    id: z.string().nonempty(),
    name: z.string().nonempty(),
  }, params).transform(station => station.id)
}

export const Band = co.map({
  name: z.string(), // Zod primitive type
});

export const Festival = co.list(Band);

export const JazzFestAccountRoot = co.map({
  myFestival: Festival,
});

export const JazzFestAccount = co
    .account({
      root: JazzFestAccountRoot,
      profile: co.profile(),
    })
    .withMigration((account) => {
      if (!account.$jazz.has('root')) {
        account.$jazz.set('root', {
          myFestival: [],
        });
      }
    });