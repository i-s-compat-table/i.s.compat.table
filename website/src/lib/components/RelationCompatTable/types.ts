import type { AllDbs, SupportRange } from "$lib/types";

export type TableSupportRow = {
  name: string;
  specific: boolean;
  universal: boolean;
  support: Partial<Record<AllDbs, SupportRange>>;
};
