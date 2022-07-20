import { useCallback } from "react";
import {
  atomFamily,
  selectorFamily,
  useRecoilState,
  useRecoilValueLoadable,
} from "recoil";
import { api, type ResponseBody } from "../apiClient";
import { AtomKeys, SelectorKeys } from "./keys";

const abortControllerAtom = atomFamily<AbortController, string>({
  key: AtomKeys.USER_LIST_ABORT_CONTROLLER,
  default: () => new AbortController(),
});

const userListSelector = selectorFamily<ResponseBody["getUser"], string>({
  key: SelectorKeys.USER_LIST,
  get:
    (query) =>
    async ({ get }) => {
      if (query === "") return { users: [] };
      const abortController = get(abortControllerAtom(query));
      const response = await api.getUser(query, {
        signal: abortController.signal,
      });
      if (response.status !== 200) return { users: [] };
      return response.data;
    },
});

export function useUserList(query: string) {
  const [abortController, setAbortController] = useRecoilState(
    abortControllerAtom(query)
  );
  const { state, contents } = useRecoilValueLoadable(userListSelector(query));
  const abort = useCallback(() => {
    abortController.abort();
    setAbortController(new AbortController());
  }, [abortController, setAbortController]);

  switch (state) {
    case "hasError":
      return {
        error: contents instanceof Error ? contents : new Error(contents),
        isLoading: false,
        userList: null,
        abort,
      } as const;

    case "loading":
      return {
        error: null,
        isLoading: true,
        userList: null,
        abort,
      } as const;

    case "hasValue":
      return {
        error: null,
        isLoading: false,
        userList: contents.users,
        abort,
      } as const;
  }
}
