export type TypeInfo = {
  typeName: string; // should always be uppercase
  nullable: boolean | null;
};

export const renderTypeInfo = (t: TypeInfo) => {
  let result = t.typeName.toUpperCase();
  if (t.nullable === true) result += " NULL";
  else if (t.nullable === false) result += " NOT NULL";
  return result;
};
export type VersionInfo = {
  versionRange: RangeRef;
  typeInfo: TypeInfo;
};
export type EachDb<DbNames extends string, T> = Record<DbNames, T>;

export type SupportRange = { range: string; isCurrent: boolean };

export type RangeRef = SupportRange & {
  url: string | null;
  note: string | null;
  license: string | null;
  license_url: string | null;
  attribution: string | null;
};
export type TableCompatibility<DbNames extends string = AllDbs> = {
  tableName: string;
  columns: Array<ColumnCompatibility<DbNames>>;
} & EachDb<DbNames, RangeRef>;

export type ColumnCompatibility<DbNames extends string = AllDbs> = {
  id: number;
  columnName: string;
} & EachDb<DbNames, VersionInfo[]>;

export type CompatibilityRow<DbNames extends string = AllDbs> =
  TableCompatibility<DbNames>;
export type AllDbs = "postgres" | "mysql" | "mariadb" | "tidb" | "cockroachdb" | "mssql";
