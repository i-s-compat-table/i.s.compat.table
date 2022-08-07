import type { AllDbs } from "$lib/types";
import { urlParam } from "./searchParams";
const allowedValues: AllDbs[] = [
  // order matters for table presentation
  "mysql",
  "mariadb",
  "tidb",
  "postgres",
  "cockroachdb",
  "mssql",
];
const fallback = [...allowedValues];
const isValid = (val: any): val is AllDbs => allowedValues.includes(val);
const parse = (v: string | null): AllDbs[] | null => {
  if (!v) return null;
  const val = v.split(",");
  return val.every(isValid) ? val : null;
};
const serialize = (val: AllDbs[]) => val.join(",");

const dbs = urlParam<AllDbs[]>("dbs", fallback, parse, serialize);
export default dbs;
