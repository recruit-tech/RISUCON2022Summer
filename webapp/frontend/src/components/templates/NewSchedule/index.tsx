import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import Error from "next/error";
import { useRouter } from "next/router";
import { useCallback, type ComponentProps } from "react";
import { api } from "../../../apiClient";
import { useCurrentUser } from "../../../states/currentUserState";
import { convertScheduleFormDataToRequestBody } from "../../../utils/convertFormDataToRequestBody";
import { ScheduleForm } from "../../organisms/ScheduleForm";
import { Loading } from "../Loading";

export const NewSchedule = () => {
  const router = useRouter();
  const { currentUser, isLoading: isLoadingCurrentUser, error: currentUserError } = useCurrentUser();
  const onSubmit = useCallback<ComponentProps<typeof ScheduleForm>["onSubmit"]>(
    (scheduleFormData) =>
      api
        .postSchedule(convertScheduleFormDataToRequestBody(scheduleFormData))
        .then(({ data: { id } }) => {
          router.push({ pathname: "/schedule/[id]", query: { id } });
        }),
    [router]
  );

  if (currentUserError) return <Error statusCode={0} title={currentUserError?.message} />;

  if (isLoadingCurrentUser) return <Loading />;

  return (
    <Box sx={{ pt: 8, display: "flex", flexDirection: "column" }}>
      <Typography component="h1" variant="h4">
        新規スケジュール作成
      </Typography>
      <ScheduleForm type="new" onSubmit={onSubmit} currentUser={currentUser} />
    </Box>
  );
};
