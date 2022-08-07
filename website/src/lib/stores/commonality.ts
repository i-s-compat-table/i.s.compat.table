import { identity, urlParam } from "./searchParams";
export type Commonality = "db-specific" | "shared" | "universal" | "any";

const validate = (val: string | null): val is Commonality => {
  return (
    typeof val === "string" && ["db-specific", "shared", "universal", "any"].includes(val)
  );
};
export const key = "commonality";
export const fallback = "any";
export const parse = (val: string | null): Commonality | null =>
  validate(val) ? val : null;
const commonality = urlParam<Commonality>(key, fallback, parse, identity);
export default commonality;
