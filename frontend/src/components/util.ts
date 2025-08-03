
export function getDaysBetween(startDate: Date, endDate: Date) {
  const dates = [];
  const currentDate = new Date(startDate);

  if (endDate < startDate) {
    return [];
  }

  while (currentDate <= endDate) {
    dates.push(new Date(currentDate));
    currentDate.setDate(currentDate.getDate() + 1);
  }

  return dates;
}

export function isSameDay(dateA: Date, dateB: Date) {
  return dateA.getFullYear() === dateB.getFullYear() &&
    dateA.getMonth() === dateB.getMonth() &&
    dateA.getDate() === dateB.getDate();
}

export function formatDuration(dateA: Date, dateB: Date) {
  const difference = dateB.getTime() - dateA.getTime();
  return formatDurationMinutes(Math.floor(difference / (1000 * 60)))
}

export function formatDurationMinutes(minutes: number) {
  return Math.floor(minutes / 60) + "h " + minutes % 60 + "min"
}

export function formatDateLong(date: Date) {
  return new Intl.DateTimeFormat(undefined, {
    dateStyle: "full"
  }).format(date)
}

export function formatDateShort(date: Date) {
  return new Intl.DateTimeFormat("de-DE", {
    dateStyle: "short"
  }).format(date)
}

export function formatTime(time: Date) {
  return time.getHours() + ":" + time.getMinutes().toString().padStart(2, "0")
}

export function formatAmount(amount: number|null, decimals = 2): string {
  if (amount == null) return ""
  const sign = amount < 0 ? "-" : ""
  let result = sign + Math.abs(amount / Math.pow(10, decimals) >> 0).toString().padStart(1, "0")
  if (decimals != 0)
    result += "," + Math.abs(amount % Math.pow(10, decimals)).toString().padStart(decimals, "0")
  return result
}

export function getDateString(date: Date) {
  const offset = date.getTimezoneOffset()
  const offsetDate = new Date(date.getTime() - (offset*60*1000))
  return offsetDate.toISOString().split('T')[0]
}

export function nullIfEmpty(value: string) {
  const trimmedValue = value.trim()
  return trimmedValue.length != 0 ? trimmedValue : null;
}
