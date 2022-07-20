import { AxiosError, AxiosResponse } from "axios";
import { isAxiosError } from "../apiClient";

const UNEXPECTED_NETWORK_ERROR_MESSAGE = "Unexpected network error";

interface ErrorMessagePayload {
  readonly message: string;
}

const isErrorMessagePayload = (value: unknown): value is ErrorMessagePayload =>
  typeof value === "object" &&
  value !== null &&
  "message" in value &&
  typeof (value as { readonly message: unknown }).message === "string";

const getDataMessage = (data: unknown): string | undefined => {
  if (typeof data === "string") return data;

  if (isErrorMessagePayload(data)) {
    return data.message;
  }

  return undefined;
};

const getResponseMessage = ({
  data,
  statusText,
}: AxiosResponse<unknown, unknown>): string =>
  getDataMessage(data) || statusText;

const getAxiosErrorMessage = ({
  response,
}: AxiosError<unknown, unknown>): string =>
  (response && getResponseMessage(response)) ||
  UNEXPECTED_NETWORK_ERROR_MESSAGE;

export const getNetworkErrorMessage = (error: unknown) =>
  isAxiosError(error)
    ? getAxiosErrorMessage(error)
    : UNEXPECTED_NETWORK_ERROR_MESSAGE;
