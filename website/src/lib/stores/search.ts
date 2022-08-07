import { urlParam } from "./searchParams";
const query = urlParam(
  "q",
  "",
  (s) => (s ? decodeURIComponent(s) : null),
  (s) => encodeURIComponent(s),
);
export default query;
