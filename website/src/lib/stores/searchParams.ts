import { browser } from "$app/env";
import type { Writable } from "svelte/store";
import { writable } from "svelte/store";

export const identity = <T>(v: T) => v;

const getInitialValue = <T>(
  store: Writable<T>,
  key: string,
  fallback: T,
  parse: (v: string | null) => T | null,
) => {
  const url = new URL(window.location.toString());
  const val = parse(url.searchParams.get(key)) ?? fallback;
  store.set(val);
};
const updateUrl = <T>(key: string, fallback: T, serialize: (v: T) => string) => {
  const _fallback = serialize(fallback);
  console.log({ _fallback });
  return (value: T) => {
    const searchParam = serialize(value);
    const url = new URL(window.location.toString());
    if (searchParam === _fallback) {
      url.searchParams.delete(key);
    } else {
      url.searchParams.set(key, searchParam);
    }
    window.history.replaceState(window.history.state, "", url);
  };
};
export const urlParam = <T>(
  key: string,
  fallback: T,
  parse: (v: string | null) => T | null,
  serialize: (v: T) => string,
) => {
  const store = writable(fallback);
  if (browser) {
    getInitialValue(store, key, fallback, parse);
    store.subscribe(updateUrl(key, fallback, serialize));
  }
  return store;
};
