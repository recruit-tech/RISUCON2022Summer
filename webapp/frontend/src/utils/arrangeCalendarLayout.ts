import { Calendar, ScheduleWithID } from "../apiClient/__generated__";
import { convertIdToHexColor } from "./convertIdToHexColor";
import { calculateTimeInDate, DATE } from "./dateTime";

type Percentage = `${number}%`;

export interface Layout {
  readonly backgroundColor: string;
  readonly top: Percentage;
  readonly left: Percentage;
  readonly width: Percentage;
  readonly height: Percentage;
}

export interface LayoutSchedule extends ScheduleWithID {
  readonly layout: Layout;
}

export interface LayoutCalendar extends Omit<Calendar, "schedules"> {
  readonly schedules: LayoutSchedule[];
}

/**
 * sort schedules immutably
 * sort by start_at (ascent), end_at (decent), id (ascent)
 * @param schedules
 * @returns new schedules array which is sorted
 */
const getSortedSchedules = (schedules: ScheduleWithID[]): ScheduleWithID[] =>
  // This process will now take place on the API Server
  /*
  [...schedules].sort((a, b) => {
    if (a.start_at !== b.start_at) return a.start_at - b.start_at;
    if (a.end_at !== b.end_at) return b.end_at - a.end_at;
    return a.id > b.id ? 1 : -1;
  });
   */
  schedules;

type ScheduleColumn = ScheduleWithID[];

/**
 * arrange schedules to schedule columns
 * @param schedules sorted schedule list
 * @returns schedule columns which is arranged
 */
const arrangeScheduleColumns = (
  schedules: ScheduleWithID[]
): ScheduleColumn[] => {
  if (schedules.length === 0) return [];
  const column: ScheduleColumn = [];
  const rest: ScheduleWithID[] = [];

  let beforeEndAt = -1;
  for (const schedule of schedules) {
    if (schedule.start_at > beforeEndAt) {
      beforeEndAt = schedule.end_at;
      column.push(schedule);
    } else {
      rest.push(schedule);
    }
  }

  return [column, ...arrangeScheduleColumns(rest)];
};

/**
 * convert `ScheduleColumn` list to `LayoutSchedule` list
 * @param scheduleColumns schedule columns which is arranged
 * @param date date counts of schedules (date from 1 January 1970 based on UNIX time)
 * @param ownerId ID of calendar owner
 * @returns schedules with layout
 */
const convertScheduleColumnsToLayouts = (
  scheduleColumns: ScheduleColumn[],
  date: Calendar["date"],
  ownerId: string
): LayoutSchedule[] => {
  const backgroundColor = convertIdToHexColor(ownerId);
  const width: Percentage = `${100 / scheduleColumns.length}%`;
  const layouts: LayoutSchedule[] = [];
  for (let i = 0; i < scheduleColumns.length; i++) {
    const left: Percentage = `${(i * 100) / scheduleColumns.length}%`;
    for (const schedule of scheduleColumns[i]) {
      const startAt = calculateTimeInDate(schedule.start_at, date);
      const endAt = calculateTimeInDate(schedule.end_at, date);
      layouts.push({
        ...schedule,
        layout: {
          backgroundColor,
          top: `${(startAt * 100) / DATE}%`,
          left,
          width,
          height: `${((endAt - startAt) / DATE) * 100}%`,
        },
      });
    }
  }
  return layouts;
};

/**
 * arrange calendar to schedule layouts
 * @param calendar
 * @returns `Calendar` include schedule layouts
 */
export const arrangeCalendarLayout = (
  calendar: Calendar,
  ownerId: string
): LayoutCalendar => {
  const schedules = getSortedSchedules(calendar.schedules);
  const scheduleColumns = arrangeScheduleColumns(schedules);
  const scheduleLayouts = convertScheduleColumnsToLayouts(
    scheduleColumns,
    calendar.date,
    ownerId
  );
  return { ...calendar, schedules: scheduleLayouts };
};
