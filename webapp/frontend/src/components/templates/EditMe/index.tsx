import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import Error from "next/error";
import { useRouter } from "next/router";
import { useCallback, type ComponentProps } from "react";
import { api } from "../../../apiClient";
import {
  useCurrentUser,
  useCurrentUserRefresher,
} from "../../../states/currentUserState";
import { UserForm } from "../../organisms/UserForm";
import { Loading } from "../Loading";

export const EditMe = () => {
  const { error, isLoading, currentUser } = useCurrentUser();
  const router = useRouter();
  const refreshCurrentUser = useCurrentUserRefresher();

  const onSubmit = useCallback<ComponentProps<typeof UserForm>["onSubmit"]>(
    (userFormData) => {
      const { icon, ...userProperties } = userFormData;
      return api
        .updateMe(userProperties)
        .then(async () => {
          if (icon === null) return;
          await api.updateMeIcon(icon);
        })
        .then(() => {
          refreshCurrentUser();
          router.push("/");
        });
    },
    [refreshCurrentUser, router]
  );

  if (error) return <Error statusCode={0} title={error?.message} />;

  if (isLoading || !currentUser) return <Loading />;

  return (
    <Box sx={{ pt: 8, display: "flex", flexDirection: "column" }}>
      <Typography component="h1" variant="h4">
        マイプロフィール編集
      </Typography>
      <UserForm type="edit" onSubmit={onSubmit} defaultValue={currentUser} />
    </Box>
  );
};
