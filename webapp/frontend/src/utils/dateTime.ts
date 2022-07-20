export const MILLISECOND = 1;
export const SECOND = 1000 * MILLISECOND;
export const MINUTE = 60 * SECOND;
export const HOUR = 60 * MINUTE;
export const DATE = 24 * HOUR;

export const getToday = () =>
  Date.UTC(/* 2022/07/08 */ 2022, 6, 8) / 24 / 60 / 60 / 1000;

export const calculateUnixTime = (
  year: number,
  month: number,
  date: number,
  hour: number,
  minute: number
) => Date.UTC(year, month, date, hour, minute) / 1000;

export const calculateDate = (unixTime: number) => new Date(unixTime * 1000);

export const calculateTimeInDate = (time: number, date: number): number =>
  Math.max(0, Math.min(time * SECOND - date * DATE, DATE));
