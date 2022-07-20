import {
  selector,
  useRecoilRefresher_UNSTABLE,
  useRecoilValueLoadable,
} from "recoil";
import { api } from "../apiClient";
import { SelectorKeys } from "./keys";

const currentUserSelector = selector({
  key: SelectorKeys.CURRENT_USER,
  get: async () => {
    const response = await api.getMe();
    return response.data;
  },
});

export function useCurrentUser() {
  const { state, contents } = useRecoilValueLoadable(currentUserSelector);

  switch (state) {
    case "hasError":
      return {
        error: contents instanceof Error ? contents : new Error(contents),
        isLoading: false,
        currentUser: null,
      } as const;

    case "loading":
      return {
        error: null,
        isLoading: true,
        currentUser: null,
      } as const;

    case "hasValue":
      return {
        error: null,
        isLoading: false,
        currentUser: contents,
      } as const;
  }
}

export function useCurrentUserRefresher() {
  return useRecoilRefresher_UNSTABLE(currentUserSelector);
}
