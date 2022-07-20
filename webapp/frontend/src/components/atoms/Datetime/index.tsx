import Typography from "@mui/material/Typography";
import { type ComponentProps, type VFC } from "react";

interface Props extends ComponentProps<typeof Typography> {
  /**
   * UNIX Timestamp
   */
  readonly time: number;
}

export const Datetime: VFC<Props> = ({ time, ...props }) => {
  const date = new Date(time);
  return (
    <Typography component="time" dateTime={date.toISOString()} {...props}>
      {`${date.getUTCFullYear()}-${(date.getUTCMonth() + 1)
        .toString()
        .padStart(2, "0")}-${`${date.getUTCDate()}`.padStart(2, "0")} ${date
        .getUTCHours()
        .toString()
        .padStart(2, "0")}:${date
        .getUTCMinutes()
        .toString()
        .padStart(2, "0")}:${date.getUTCSeconds().toString().padStart(2, "0")}`}
    </Typography>
  );
};
