import { selectorFamily, useRecoilValueLoadable } from "recoil";
import { api } from "../apiClient";
import { ScheduleWithID } from "../apiClient/__generated__";
import { SelectorKeys } from "./keys";

const scheduleSelector = selectorFamily<
  ScheduleWithID,
  ScheduleWithID["id"] | undefined
>({
  key: SelectorKeys.USER,
  get: (id) => async () => {
    if (id === undefined) throw new Error("id is undefined");
    const response = await api.getScheduleId(id);
    return response.data;
  },
});

export function useSchedule(id: ScheduleWithID["id"] | undefined) {
  const { state, contents } = useRecoilValueLoadable(scheduleSelector(id));

  switch (state) {
    case "hasError":
      return {
        error: contents instanceof Error ? contents : new Error(contents),
        isLoading: false,
        schedule: null,
      } as const;

    case "loading":
      return {
        error: null,
        isLoading: true,
        schedule: null,
      } as const;

    case "hasValue":
      return {
        error: null,
        isLoading: false,
        schedule: contents,
      } as const;
  }
}
