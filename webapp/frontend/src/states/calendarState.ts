import { selectorFamily, useRecoilValueLoadable } from "recoil";
import { api } from "../apiClient";
import { Calendar } from "../apiClient/__generated__";
import {
  arrangeCalendarLayout,
  LayoutCalendar,
} from "../utils/arrangeCalendarLayout";
import { SelectorKeys } from "./keys";

const calendarSelector = selectorFamily<Calendar, [id: string, date: number]>({
  key: SelectorKeys.CALENDAR,
  get:
    ([id, date]) =>
    async () => {
      const response = await api.getCalendarId(id, date);
      return response.data;
    },
});

const calendarLayoutSelector = selectorFamily<
  LayoutCalendar,
  [id: string, date: number]
>({
  key: SelectorKeys.CALENDAR_LAYOUT,
  get:
    ([id, date]) =>
    async ({ get }) => {
      const calendar = get(calendarSelector([id, date]));
      return arrangeCalendarLayout(calendar, id);
    },
});

export function useCalendarLayout(id: string, date: number) {
  const { state, contents } = useRecoilValueLoadable(
    calendarLayoutSelector([id, date])
  );

  switch (state) {
    case "hasError":
      return {
        error: contents instanceof Error ? contents : new Error(contents),
        isLoading: false,
        calendarLayout: null,
      } as const;

    case "loading":
      return {
        error: null,
        isLoading: true,
        calendarLayout: null,
      } as const;

    case "hasValue":
      return {
        error: null,
        isLoading: false,
        calendarLayout: contents,
      } as const;
  }
}
