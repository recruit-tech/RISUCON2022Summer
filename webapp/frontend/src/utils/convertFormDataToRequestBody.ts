import type { ScheduleWithAttendeeIDs } from "../apiClient/__generated__";
import { calculateUnixTime } from "./dateTime";
import type { ScheduleFormData } from "./validates";

/**
 * MUI.DateTimePicker set time as local time zone value.
 * `calculateUnixTimeFromSelectedDate` function convert it to UNIX Timezone as UTC
 */
const calculateUnixTimeFromSelectedDate = (selectedDate: Date) =>
  calculateUnixTime(
    selectedDate.getFullYear(),
    selectedDate.getMonth(),
    selectedDate.getDate(),
    selectedDate.getHours(),
    selectedDate.getMinutes()
  );

export const convertScheduleFormDataToRequestBody = ({
  title,
  description,
  attendeeSet,
  startAt,
  endAt,
  meetingRoom,
}: ScheduleFormData): ScheduleWithAttendeeIDs => ({
  title,
  description,
  attendees: Array.from(attendeeSet),
  start_at: calculateUnixTimeFromSelectedDate(startAt),
  end_at: calculateUnixTimeFromSelectedDate(endAt),
  meeting_room: meetingRoom,
});
