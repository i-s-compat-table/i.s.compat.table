export const asc = <T>(a: T, b: T) => (a > b ? 1 : b > a ? -1 : 0);
export const objEntries = <T>(obj: Record<string, T>): [string, T][] => {
  return Object.entries(obj).sort(([a], [b]) => asc(a, b));
};
export const byKey = <T>(key: keyof T) => {
  return (a: T, b: T) => asc(a[key], b[key]);
};
