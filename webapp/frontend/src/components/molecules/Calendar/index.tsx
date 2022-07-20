import AddIcon from "@mui/icons-material/Add";
import ArrowBackIosNewIcon from "@mui/icons-material/ArrowBackIosNew";
import ArrowForwardIosIcon from "@mui/icons-material/ArrowForwardIos";
import ButtonGroup from "@mui/material/ButtonGroup";
import Grid from "@mui/material/Grid";
import IconButton from "@mui/material/IconButton";
import Link from "next/link";
import { useCallback, useState, type VFC } from "react";
import { User } from "../../../apiClient/__generated__";
import { getToday } from "../../../utils/dateTime";
import { CalendarColumn } from "./CalendarColumn";

interface Props {
  readonly userId: User["id"];
}

export const Calendar: VFC<Props> = ({ userId }) => {
  const [date, setDate] = useState(getToday());
  const clickBackwardButton = useCallback(
    () => setDate((d) => d - 1),
    [setDate]
  );
  const clickForwardButton = useCallback(
    () => setDate((d) => d + 1),
    [setDate]
  );

  return (
    <Grid container spacing={0.5}>
      <Grid item xs={12} key="button-group">
        <ButtonGroup sx={{ display: "flex", justifyContent: "flex-end" }}>
          <IconButton
            color="primary"
            aria-label="backward"
            onClick={clickBackwardButton}
          >
            <ArrowBackIosNewIcon />
          </IconButton>
          {/* eslint-disable-next-line @next/next/link-passhref */}
          <Link href="/schedule/new">
            <IconButton color="primary" aria-label="new">
              <AddIcon />
            </IconButton>
          </Link>
          <IconButton
            color="primary"
            aria-label="forward"
            onClick={clickForwardButton}
          >
            <ArrowForwardIosIcon />
          </IconButton>
        </ButtonGroup>
      </Grid>
      <Grid item xs={4} key={date}>
        <CalendarColumn userId={userId} date={date} showHour />
      </Grid>
      <Grid item xs={4} key={date + 1}>
        <CalendarColumn userId={userId} date={date + 1} />
      </Grid>
      <Grid item xs={4} key={date + 2}>
        <CalendarColumn userId={userId} date={date + 2} />
      </Grid>
    </Grid>
  );
};
