import AdapterDateFns from "@mui/lab/AdapterDateFns";
import DateTimePicker, {
  type DateTimePickerProps
} from "@mui/lab/DateTimePicker";
import LocalizationProvider from "@mui/lab/LocalizationProvider";
import Alert from "@mui/material/Alert";
import MuiAvatar from "@mui/material/Avatar";
import AvatarGroup from "@mui/material/AvatarGroup";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import { red } from "@mui/material/colors";
import MenuItem from "@mui/material/MenuItem";
import Modal from "@mui/material/Modal";
import Select from "@mui/material/Select";
import { SelectInputProps } from "@mui/material/Select/SelectInput";
import Skeleton from "@mui/material/Skeleton";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";
import Error from "next/error";
import {
  FormEventHandler,
  useCallback,
  useState,
  type ComponentProps
} from "react";
import { ScheduleWithID, User } from "../../../apiClient/__generated__";
import { useUsers } from "../../../states/userState";
import { calculateDate } from "../../../utils/dateTime";
import { getNetworkErrorMessage } from "../../../utils/getNetworkErrorMessage";
import { MEETING_ROOM_LIST } from "../../../utils/meetingRoom";
import {
  MaybeScheduleFormData,
  ScheduleFormData,
  validateScheduleData,
  ValidationError
} from "../../../utils/validates";
import { Avatar } from "../../atoms/Avatar";
import { UserChoiceList } from "../../organisms/UserChoiceList";
import { Loading } from "../../templates/Loading";

type Props = {
  readonly onSubmit: (scheduleFormData: ScheduleFormData) => Promise<void>;
  readonly currentUser: User;
} & (
  | { readonly type: "new"; readonly defaultValue?: undefined }
  | { readonly type: "edit"; readonly defaultValue: ScheduleWithID }
);

export const ScheduleForm = ({ onSubmit, type, defaultValue, currentUser }: Props) => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const onOpenModal = useCallback(() => setIsModalOpen(true), [setIsModalOpen]);
  const onCloseModal = useCallback(
    () => setIsModalOpen(false),
    [setIsModalOpen]
  );

  const [attendeeSet, setAttendeeSet] = useState(
    new Set([currentUser.id, ...(defaultValue ? defaultValue.attendees.map((attendee) => attendee.id, currentUser.id) : [])])
  );
  const onToggleUser = useCallback<
    ComponentProps<typeof UserChoiceList>["onToggleUser"]
  >(
    (userId, toBeChecked) => {
      if (userId === currentUser.id) return;
      if (toBeChecked) {
        setAttendeeSet((set) => new Set(set).add(userId));
      } else {
        setAttendeeSet((set) => {
          const result = new Set(set);
          result.delete(userId);
          return result;
        });
      }
    },
    [setAttendeeSet, currentUser.id]
  );
  const { error, isLoading, users } = useUsers(Array.from(attendeeSet));

  const [startAt, setStartAt] = useState<Date | null>(
    defaultValue ? calculateDate(defaultValue.start_at) : null
  );
  const [endAt, setEndAt] = useState<Date | null>(
    defaultValue ? calculateDate(defaultValue.end_at) : null
  );
  const onChangeStartAt: DateTimePickerProps<Date>["onChange"] = useCallback(
    (date) => setStartAt(date),
    [setStartAt]
  );
  const onChangeEndAt: DateTimePickerProps<Date>["onChange"] = useCallback(
    (date) => setEndAt(date),
    [setEndAt]
  );

  const [meetingRoom, setMeetingRoom] = useState(
    defaultValue?.meeting_room ?? ""
  );
  const onChangeMeetingRoom = useCallback<
    NonNullable<SelectInputProps<string>["onChange"]>
  >((e) => {
    setMeetingRoom(e.target.value);
  }, []);

  const [isSubmitting, setIsSubmitting] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const handleSubmit = useCallback<FormEventHandler<HTMLFormElement>>(
    (event) => {
      event.preventDefault();
      const data = new FormData(event.currentTarget);
      const title = data.get("title");
      const description = data.get("description") ?? "";

      try {
        const scheduleFormData: MaybeScheduleFormData = {
          title,
          description,
          attendeeSet,
          startAt,
          endAt,
          meetingRoom,
        };
        validateScheduleData(scheduleFormData);

        setIsSubmitting(true);
        setErrorMessage(null);
        onSubmit(scheduleFormData)
          .catch((error) => {
            console.error(error);
            setErrorMessage(getNetworkErrorMessage(error));
          })
          .finally(() => {
            setIsSubmitting(false);
          });
      } catch (e) {
        if (e instanceof ValidationError) {
          setErrorMessage(e.message);
        }
      }
    },
    [attendeeSet, startAt, endAt, meetingRoom, onSubmit]
  );

  if (!isModalOpen && isLoading) return <Loading />;

  if (error !== null) return <Error statusCode={0} title={error.message} />;

  return (
    <Box component="form" onSubmit={handleSubmit}>
      <Box sx={{ mt: 2 }}>
        <Typography component="h2" variant="h6">
          タイトル
          <Typography component="span" color={red[500]}>
            *
          </Typography>
        </Typography>
        <TextField
          name="title"
          sx={{ mt: 1 }}
          variant="filled"
          required
          fullWidth
          margin="none"
          defaultValue={defaultValue?.title}
        />
      </Box>
      <Box sx={{ mt: 4 }}>
        <Typography component="h2" variant="h6">
          日付
          <Typography component="span" color={red[500]}>
            *
          </Typography>
        </Typography>
        <LocalizationProvider dateAdapter={AdapterDateFns}>
          <Box sx={{ mt: 2, pl: 1, display: "flex", alignItems: "center" }}>
            <DateTimePicker
              renderInput={(props) => <TextField required {...props} />}
              label="From"
              value={startAt}
              onChange={onChangeStartAt}
            />
            <Typography sx={{ mx: 2 }}>~</Typography>
            <DateTimePicker
              renderInput={(props) => <TextField required {...props} />}
              label="To"
              value={endAt}
              onChange={onChangeEndAt}
            />
          </Box>
        </LocalizationProvider>
      </Box>
      <Box sx={{ mt: 4 }}>
        <Typography component="h2" variant="h6">
          出席者
          <Typography component="span" color={red[500]}>
            *
          </Typography>
        </Typography>
        <Box
          sx={{
            mt: 1,
            pl: 1,
            display: "flex",
            flexDirection: "row",
            alignItems: "center",
          }}
        >
          {isLoading ? (
            <Skeleton variant="circular">
              <AvatarGroup>
                <MuiAvatar />
              </AvatarGroup>
            </Skeleton>
          ) : users.length > 0 ? (
            <AvatarGroup max={5}>
              {users.map((user) => (
                <Avatar key={user.id} userId={user.id} name={user.name} />
              ))}
            </AvatarGroup>
          ) : (
            <Typography>参加者なし</Typography>
          )}
          <Button sx={{ ml: 1 }} onClick={onOpenModal}>
            編集
          </Button>
          <Modal open={isModalOpen} onClose={onCloseModal}>
            <Box
              sx={{
                position: "absolute",
                left: "50%",
                top: "50%",
                width: { xs: "100%", sm: 720 },
                height: "80vh",
                transform: "translate(-50%, -50%)",
                overflowX: "hidden",
                overflowY: "scroll",
              }}
            >
              <UserChoiceList
                checkedSet={attendeeSet}
                onToggleUser={onToggleUser}
              />
            </Box>
          </Modal>
        </Box>
      </Box>
      <Box sx={{ mt: 4 }}>
        <Typography component="h2" variant="h6">
          会議室
        </Typography>
        <Select
          sx={{ minWidth: 128 }}
          value={meetingRoom}
          onChange={onChangeMeetingRoom}
        >
          <MenuItem value="">未指定</MenuItem>
          {MEETING_ROOM_LIST.map((meetingRoom) => (
            <MenuItem key={meetingRoom} value={meetingRoom}>
              {meetingRoom}
            </MenuItem>
          ))}
        </Select>
      </Box>
      <Box sx={{ mt: 4 }}>
        <Typography component="h2" variant="h6">
          概要
        </Typography>
        <TextField
          name="description"
          sx={{ mt: 1 }}
          variant="filled"
          fullWidth
          multiline
          rows={3}
          margin="none"
          defaultValue={defaultValue?.description}
        />
      </Box>
      <Alert
        variant="outlined"
        severity="error"
        sx={{
          marginTop: 4,
          width: "100%",
          visibility: errorMessage !== null ? "visible" : "hidden",
        }}
      >
        {errorMessage}
      </Alert>
      <Button
        type="submit"
        disabled={isSubmitting}
        fullWidth
        variant="contained"
        sx={{ mt: 3, mb: 2 }}
      >
        {type === "new" ? "作成" : "更新"}
      </Button>
    </Box>
  );
};
