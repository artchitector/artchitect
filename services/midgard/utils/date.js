import moment from "moment";

export function format(date) {
  return !!date ? moment(date).format("YYYY MMM Do HH:mm:ss") : '---'
}
