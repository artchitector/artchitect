import moment from "moment";

export function format(date, locale) {
  let format = "YYYY MMM Do HH:mm:ss";
  if (locale !== undefined) {
    moment.locale(locale)
    if (locale === "ru") {
      format = "D MMMM YYYY HH:mm:ss"
    }
  }
  return !!date ? moment(date).format(format) : '---'
}
