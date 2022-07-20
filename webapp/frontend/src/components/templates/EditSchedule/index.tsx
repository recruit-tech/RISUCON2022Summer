import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import Error from "next/error";
import { useRouter } from "next/router";
import { useCallback, type ComponentProps } from "react";
import { api } from "../../../apiClient";
import { useCurrentUser } from "../../../states/currentUserState";
import { useSchedule } from "../../../states/scheduleState";
import { convertScheduleFormDataToRequestBody } from "../../../utils/convertFormDataToRequestBody";
import { getQuery } from "../../../utils/getQuery";
import { ScheduleForm } from "../../organisms/ScheduleForm";
import { Loading } from "../Loading";

export const EditSchedule = () => {
  const router = useRouter();
  const id = getQuery(router, "id");
  const { error, isLoading, schedule } = useSchedule(id);
  const { currentUser, isLoading: isLoadingCurrentUser, error: currentUserError } = useCurrentUser();

  const onSubmit = useCallback<ComponentProps<typeof ScheduleForm>["onSubmit"]>(
    (scheduleFormData) => {
      if (id === undefined) return Promise.reject();

      return api
        .updateScheduleId(
          id,
          convertScheduleFormDataToRequestBody(scheduleFormData)
        )
        .then(() => {
          router.push({ pathname: "/schedule/[id]", query: { id } });
        });
    },
    [id, router]
  );

  if (error || currentUserError) return <Error statusCode={0} title={error?.message} />;

  if (id == undefined || isLoading || isLoadingCurrentUser || !schedule) return <Loading />;

  return (
    <Box sx={{ pt: 8, display: "flex", flexDirection: "column" }}>
      <Typography component="h1" variant="h4">
        スケジュール編集
      </Typography>
      <ScheduleForm type="edit" onSubmit={onSubmit} defaultValue={schedule} currentUser={currentUser} />
    </Box>
  );
};
