import { NextRouter } from "next/router";

export function getQuery(router: NextRouter, param: string) {
  const value = router.query[param];
  return Array.isArray(value) ? value[0] : value;
}
