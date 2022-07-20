import Box from "@mui/material/Box";
import Skeleton from "@mui/material/Skeleton";
import Typography from "@mui/material/Typography";
import { type VFC } from "react";
import { User } from "../../../apiClient/__generated__";
import { useCalendarLayout } from "../../../states/calendarState";
import { theme } from "../../../styles/theme";
import { DATE } from "../../../utils/dateTime";
import { ErrorBox } from "../../atoms/ErrorBox";
import { ScheduleCell } from "./ScheduleCell";

interface Props {
  readonly userId: User["id"];
  readonly date: number;
  readonly showHour?: boolean;
}

const MINIMUM_BLOCK_HEIGHT = "2rem";
const HEIGHT = `calc(${MINIMUM_BLOCK_HEIGHT} * 48)`;

const getDateFormat = (date: number) => {
  const d = new Date(date * DATE);
  return `${d.getUTCFullYear()}-${`${d.getUTCMonth() + 1}`.padStart(
    2,
    "0"
  )}-${`${d.getUTCDate()}`.padStart(2, "0")}`;
};

export const CalendarColumn: VFC<Props> = ({ userId, date, showHour }) => {
  const { error, isLoading, calendarLayout } = useCalendarLayout(userId, date);

  if (error !== null) return <ErrorBox error={error} sx={{ height: HEIGHT }} />;

  if (isLoading) return <Skeleton variant="rectangular" height={HEIGHT} />;

  return (
    <>
      <Typography sx={{ textAlign: "center" }}>
        {getDateFormat(date)}
      </Typography>
      <Box
        sx={{
          position: "relative",
          height: HEIGHT,
          borderColor: "gray",
          borderWidth: 1,
          borderStyle: "solid",
        }}
      >
        {new Array(47).fill(0).map((_, i) => (
          <Box
            key={`ruled-${i}`}
            sx={{
              position: "absolute",
              top: `calc(${MINIMUM_BLOCK_HEIGHT} * ${i})`,
              left: 0,
              width: "100%",
              height: MINIMUM_BLOCK_HEIGHT,
              border: "none",
              borderBottomColor: "gray",
              borderBottomWidth: 1,
              borderBottomStyle: i % 2 === 0 ? "dashed" : "solid",
              "&::before":
                showHour && i % 2 === 0
                  ? {
                      display: "block",
                      position: "absolute",
                      content: `'${i / 2}:00'`,
                      transform: "translate(-120%, -50%)",
                      fontSize: theme.typography.caption,
                    }
                  : undefined,
            }}
          />
        ))}
        {calendarLayout.schedules.map((schedule) => (
          <ScheduleCell key={schedule.id} schedule={schedule} />
        ))}
      </Box>
    </>
  );
};
