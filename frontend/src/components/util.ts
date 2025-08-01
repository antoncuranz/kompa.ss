
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

export function formatDate(date: Date) {
  return new Intl.DateTimeFormat(undefined, {
    dateStyle: "full"
  }).format(date)
  // return date.toLocaleDateString("de-DE")
}

export function formatTime(time: Date) {
  return time.getHours() + ":" + time.getMinutes()
}
