import { type AxiosError, type AxiosResponse } from "axios";
import { Configuration, DefaultApiFactory } from "./__generated__";

const config = new Configuration({
  basePath: "/api",
});

export const api = DefaultApiFactory(config);

// ref. https://zenn.dev/takepepe/articles/openapi-generator-cli-ts
type InferSecondFromBack<T> = T extends [...any[], infer U, any] ? U : never;
type InferFunctionArgs<T> = T extends (...arg: infer U) => any
  ? Required<U>
  : never;
type InferRequestBody<T> = {
  [K in keyof T]: InferSecondFromBack<InferFunctionArgs<T[K]>>;
};
type InferPromise<T> = T extends (...arg: any) => Promise<infer I> ? I : never;
type InferAxiosResponse<T> = T extends AxiosResponse<infer I> ? I : never;
type InferResponseBody<T> = {
  [K in keyof T]: InferAxiosResponse<InferPromise<T[K]>>;
};

export type RequestBody = InferRequestBody<typeof api>;
export type ResponseBody = InferResponseBody<typeof api>;

export function isAxiosError(value: any): value is AxiosError {
  return (
    typeof value === "object" &&
    "isAxiosError" in value &&
    value.isAxiosError === true
  );
}
