export const getRelationPath = (base: string, relation: string) =>
  `${base}/relation/${relation}`;
export const getColumnPath = (base: string, relation: string, column: string) =>
  `${base}/column/${relation}.${column}`;
export const issues = "https://github.com/i-s-compat-table/i.s.compat.table/issues";
