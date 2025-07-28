import moment, {Moment} from "moment";
import SuperJSON from "superjson";

export function getDaysBetween(startDate: Moment, endDate: Moment) {
  const dates = [];
  const currentDate = moment(startDate);

  if (endDate < startDate) {
    return [];
  }

  while (currentDate <= endDate) {
    dates.push(moment(currentDate));
    currentDate.add(1, "day")
  }

  return dates;
}

export function isSameLocalDay(dateA: Moment, dateB: Moment) {
  return dateA.year() === dateB.year() &&
    dateA.month() === dateB.month() &&
    dateA.date() === dateB.date();
}

export function durationString(dateA: Moment, dateB: Moment) {
  const durationMinutes = dateA.diff(dateB, "minutes")
  return Math.floor(durationMinutes / 60) + "h " + durationMinutes % 60 + "min"
}

export function registerMomentSerde() {
  SuperJSON.registerCustom<Moment, string>(
    {
      isApplicable: (v): v is Moment => moment.isMoment(v),
      serialize: v => v.toJSON(),
      deserialize: v => moment(v),
    },
    'moment.js'
  )
}
