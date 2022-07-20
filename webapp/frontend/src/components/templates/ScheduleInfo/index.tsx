import EditIcon from "@mui/icons-material/Edit";
import AvatarGroup from "@mui/material/AvatarGroup";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import IconButton from "@mui/material/IconButton";
import Modal from "@mui/material/Modal";
import Typography from "@mui/material/Typography";
import Error from "next/error";
import Link from "next/link";
import { useRouter } from "next/router";
import { useCallback, useState } from "react";
import { useSchedule } from "../../../states/scheduleState";
import { useUsers } from "../../../states/userState";
import { theme } from "../../../styles/theme";
import { getQuery } from "../../../utils/getQuery";
import { Avatar } from "../../atoms/Avatar";
import { Datetime } from "../../atoms/Datetime";
import { UserList } from "../../organisms/UserList";
import { Loading } from "../Loading";

export const ScheduleInfo = () => {
  const router = useRouter();
  const id = getQuery(router, "id");
  const {
    error: scheduleError,
    isLoading: isLoadingSchedule,
    schedule,
  } = useSchedule(id);
  const {
    error: usersError,
    isLoading: isLoadingUsers,
    users,
  } = useUsers(schedule !== null ? schedule.attendees.map(({ id }) => id) : []);

  const [open, setOpen] = useState(false);
  const onOpenModal = useCallback(() => setOpen(true), [setOpen]);
  const onCloseModal = useCallback(() => setOpen(false), [setOpen]);

  if (id == undefined || isLoadingSchedule || isLoadingUsers)
    return <Loading />;

  if (scheduleError !== null || usersError !== null)
    return (
      <Error statusCode={0} title={(scheduleError ?? usersError)?.message} />
    );

  return (
    <Box sx={{ pt: 8 }}>
      <Typography paragraph sx={{ display: "flex" }}>
        <Typography component="h1" variant="h4">
          {schedule.title}
        </Typography>
        {/* eslint-disable-next-line @next/next/link-passhref */}
        <Link href={{ pathname: "/schedule/[id]/edit", query: { id } }}>
          <IconButton color="default" aria-label="edit">
            <EditIcon />
          </IconButton>
        </Link>
      </Typography>
      <Typography paragraph sx={{ mt: 2 }}>
        <Typography component="span" sx={{ mr: 1 }}>
          日付:
        </Typography>
        <Datetime time={schedule.start_at} />
        {" ~ "}
        <Datetime time={schedule.end_at} />
      </Typography>
      <Box
        sx={{
          mt: 2,
          width: "fit-content",
          display: "flex",
          alignItems: "center",
        }}
      >
        <Typography component="span" sx={{ mr: 1 }}>
          出席者:
        </Typography>
        <AvatarGroup max={5} total={schedule.attendees.length}>
          {users.map((user) => (
            <Avatar key={user.id} userId={user.id} name={user.name} />
          ))}
        </AvatarGroup>
        <Button sx={{ ml: 1 }} onClick={onOpenModal}>
          詳細を見る
        </Button>
        <Modal open={open} onClose={onCloseModal}>
          <Box
            sx={{
              position: "absolute",
              left: "50%",
              top: "50%",
              width: { xs: "100%", sm: 720 },
              height: "80vh",
              padding: 4,
              borderRadius: 2,
              backgroundColor: theme.palette.background.default,
              transform: "translate(-50%, -50%)",
              overflowX: "hidden",
              overflowY: "scroll",
            }}
          >
            <UserList users={users} />
          </Box>
        </Modal>
      </Box>
      <Typography paragraph sx={{ mt: 2 }}>
        <Typography component="span" sx={{ mr: 1 }}>
          会議室:
        </Typography>
        {schedule.meeting_room || "未指定"}
      </Typography>
      <Typography paragraph sx={{ mt: 2 }}>
        <Typography component="span" sx={{ mr: 1 }}>
          概要:
        </Typography>
        {schedule.description}
      </Typography>
    </Box>
  );
};
