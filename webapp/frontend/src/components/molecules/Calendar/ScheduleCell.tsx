import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import Link from "next/link";
import { type VFC } from "react";
import { LayoutSchedule } from "../../../utils/arrangeCalendarLayout";

interface Props {
  readonly schedule: LayoutSchedule;
}

export const ScheduleCell: VFC<Props> = ({ schedule }) => (
  // eslint-disable-next-line @next/next/link-passhref
  <Link
    href={{
      pathname: "/schedule/[id]",
      query: { id: schedule.id },
    }}
  >
    <Box
      sx={{
        position: "absolute",
        top: schedule.layout.top,
        left: schedule.layout.left,
        width: schedule.layout.width,
        height: schedule.layout.height,
        backgroundColor: schedule.layout.backgroundColor,
        borderColor: "white",
        borderWidth: 1,
        borderStyle: "solid",
        overflow: "hidden",
        cursor: "pointer",
      }}
    >
      <Typography color="white">{schedule.title}</Typography>
    </Box>
  </Link>
);
