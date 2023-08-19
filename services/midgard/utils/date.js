import moment from "moment";

export function format(date, locale) {
  let format = "YYYY MMM Do HH:mm";
  if (locale !== undefined) {
    moment.locale(locale)
    if (locale === "ru") {
      format = "D MMMM YYYY HH:mm"
    }
  }
  return !!date ? moment(date).format(format) : '---'
}

export function duration(start, end, locale) {
  if (locale !== undefined) {
    moment.locale(locale)
  }
  let a = moment(start)
  let b = moment(end)
  return moment.duration(b.diff(a)).humanize(true)
}
