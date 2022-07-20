import { selectorFamily, useRecoilValueLoadable, waitForAll } from "recoil";
import { api } from "../apiClient";
import { User } from "../apiClient/__generated__";
import { SelectorKeys } from "./keys";

const userSelector = selectorFamily<User, User["id"] | undefined>({
  key: SelectorKeys.USER,
  get: (id) => async () => {
    if (id === undefined) throw new Error("id is undefined");
    const response = await api.getUserId(id);
    return response.data;
  },
});

export function useUser(id: User["id"] | undefined) {
  const { state, contents } = useRecoilValueLoadable(userSelector(id));

  switch (state) {
    case "hasError":
      return {
        error: contents instanceof Error ? contents : new Error(contents),
        isLoading: false,
        user: null,
      } as const;

    case "loading":
      return {
        error: null,
        isLoading: true,
        user: null,
      } as const;

    case "hasValue":
      return {
        error: null,
        isLoading: false,
        user: contents,
      } as const;
  }
}

export function useUsers(ids: User["id"][]) {
  const { state, contents } = useRecoilValueLoadable(
    waitForAll(ids.map((id) => userSelector(id)))
  );

  switch (state) {
    case "hasError":
      return {
        error: contents instanceof Error ? contents : new Error(contents),
        isLoading: false,
        users: null,
      } as const;

    case "loading":
      return {
        error: null,
        isLoading: true,
        users: null,
      } as const;

    case "hasValue":
      return {
        error: null,
        isLoading: false,
        users: contents,
      } as const;
  }
}
