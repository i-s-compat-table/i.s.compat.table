import { asc } from "$lib/sort";
import { renderTypeInfo, type AllDbs, type RangeRef, type TypeInfo } from "$lib/types";
import { tsvParse } from "d3-dsv";
import semver from "semver";
export type RawData = Record<
  | "id"
  | "table_name"
  | "column_name"
  | "column_type"
  | "nullable"
  | "url"
  | "note"
  | "versions"
  | "license"
  | "license_url"
  | "attribution",
  string
> & { db_name: AllDbs };

type IntermediateData = Record<
  | "id"
  | "table_name"
  | "column_name"
  | "column_type"
  | "nullable"
  | "url"
  | "note"
  | "license"
  | "license_url"
  | "attribution",
  string
> & { db_name: AllDbs } & { versions: Array<string | number> };

export type UniqueVersionRange = RangeRef & { versions: Set<string | number> };
export type Munged<DbNames extends string = AllDbs> = {
  [table_name in string]: {
    [column_name in string]: Partial<{
      [db_name in DbNames]: {
        [kind in string]: Array<UniqueVersionRange>;
      };
    }>;
  };
};

export type TableSupport<DbNames extends string = AllDbs> = {
  [table_name in string]: {
    [db_name in DbNames]?: { range: string; isCurrent: boolean };
  };
};

export function sortByVersion(a: string | number, b: string | number) {
  if (
    typeof a === "string" &&
    semver.valid(a) &&
    typeof b === "string" &&
    semver.valid(b)
  )
    return semver.gt(a, b) ? 1 : semver.gt(b, a) ? -1 : 0;
  return asc(a, b);
}
export const munge = (
  tsv: string,
  filter: (r: RawData) => boolean = () => true,
): [Munged, TableSupport] => {
  const parsed = (tsvParse(tsv) as Array<RawData>).filter(filter);
  const intermediate: IntermediateData[] = parsed.map((d: RawData) => {
    const versions = d.versions
      .split(",")
      .map((i) => {
        const f = parseFloat(i);
        if (Number.isNaN(f)) return i;
        else return f + 0.0;
      })
      .sort();
    const result: IntermediateData = { ...d, versions }; // break ref to original obj
    return result;
  });

  const uniqueVersions: Partial<Record<AllDbs, Set<string | number>>> =
    intermediate.reduce((a: Partial<Record<AllDbs, Set<string | number>>>, d) => {
      const s = a[d.db_name] || new Set();
      d.versions.forEach((v) => s.add(v));
      a[d.db_name] = s;
      return a;
    }, {});

  const _versions = Object.entries(uniqueVersions).reduce(
    (a: Record<string, (string | number)[]>, [k, v]) => {
      a[k] = Array.from(v).sort(asc);
      return a;
    },
    {},
  );
  function getRange(db: string, vs: (string | number)[]) {
    const allVersions = _versions[db]; // already sorted in ascending order
    vs = [...new Set(vs)].sort(sortByVersion);
    let range = "";
    const vns = vs.map((v) => allVersions.indexOf(v));
    vs.forEach((v, i) => {
      const current = vns[i];
      const prev = vns[i - 1];
      const next = vns[i + 1];
      if (prev === undefined) {
        // v is start
        range += String(v);
      } else if (next === undefined) {
        // v is end
        if (range && !range.endsWith("-")) {
          range += current - 1 === prev ? "-" : ", ";
        }
        range += String(v);
      } else if (current - 1 !== prev) {
        // start new sub-range
        if (range) range += ", ";
        range += String(v);
      } else if (current + 1 !== next) {
        // end sub-range
        range += `-${v}`;
      } else {
        // mid-sub-range
        return;
      }
    });
    return range;
  }
  function getCurrent(db: string, vs: (string | number)[]) {
    return vs.some((v) => {
      if (Number.isFinite(v)) {
        return v === last(_versions[db]);
      } else if (typeof v === "string") {
        return /current/i.test(v);
      } else {
        return false;
      }
    });
  }
  const last = <T>(arr: T[]) => arr[arr.length - 1];
  const munged: Munged = intermediate.reduce((tree: Munged, r: IntermediateData) => {
    const {
      table_name,
      column_name,
      column_type,
      nullable,
      versions,
      db_name,
      note,
      url,
      license,
      license_url,
      attribution,
    } = r;
    const tableInfo = tree[table_name] ?? {};
    const columnInfo = tableInfo[column_name] ?? {};
    const dbInfo = columnInfo[db_name] ?? {};

    const nullDisplay = { true: true, false: false }[nullable] ?? null;
    const typeInfo: TypeInfo = { typeName: column_type, nullable: nullDisplay };
    const kind = renderTypeInfo(typeInfo);

    const dbColTypeVersions = dbInfo[kind] ?? [];
    const rr: RangeRef & { versions: Set<string | number> } = {
      range: "",
      versions: new Set(),
      isCurrent: getCurrent(db_name, versions),
      url,
      note,
      license,
      license_url,
      attribution,
    };
    versions.forEach((v) => rr.versions.add(v));
    rr.range = getRange(db_name, versions);
    // persist updates
    dbColTypeVersions.push(rr);
    dbInfo[kind] = dbColTypeVersions;
    columnInfo[db_name] = dbInfo;
    tableInfo[column_name] = columnInfo;
    tree[table_name] = tableInfo;
    return tree;
  }, {});

  const tableSupport = Object.entries(munged).reduce(
    (acc: TableSupport, [tableName, columns]) => {
      const dbVersions: Partial<Record<AllDbs, Set<string | number>>> = {};
      Object.values(columns).forEach((colSupport) => {
        Object.entries(colSupport).forEach(([db, kinds]) => {
          const versions = dbVersions[db as AllDbs] ?? new Set();
          Object.values(kinds).forEach((ranges) => {
            ranges.forEach((range) => {
              range.versions.forEach((v) => versions.add(v));
            });
          });
          dbVersions[db as AllDbs] = versions;
        });
      });
      const dbSupport = Object.entries(dbVersions).reduce(
        (acc: TableSupport[string], [dbName, versions]) => {
          const vs = Array.from(versions);
          acc[dbName as AllDbs] = {
            range: getRange(dbName, vs.sort()),
            isCurrent: getCurrent(dbName, vs),
          };
          return acc;
        },
        {},
      );
      acc[tableName] = dbSupport;
      return acc;
    },
    {},
  );
  return [munged, tableSupport];
};
