import { identity, urlParam } from "./searchParams";
export type Commonality = "db-specific" | "shared" | "universal" | "any";
const validate = (val: string | null): val is Commonality => {
  return (
    typeof val === "string" && ["db-specific", "shared", "universal", "any"].includes(val)
  );
};
const parse = (val: string | null): Commonality | null => (validate(val) ? val : null);
const commonality = urlParam<Commonality>("commonality", "any", parse, identity);
export default commonality;
