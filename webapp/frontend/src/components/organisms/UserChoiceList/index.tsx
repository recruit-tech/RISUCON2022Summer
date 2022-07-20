import Alert from "@mui/material/Alert";
import Box from "@mui/material/Box";
import CircularProgress from "@mui/material/CircularProgress";
import List from "@mui/material/List";
import TextField from "@mui/material/TextField";
import { useCallback, useState, type ComponentProps, type VFC } from "react";
import { useUserList } from "../../../states/userListState";
import { useUsers } from "../../../states/userState";
import { theme } from "../../../styles/theme";
import { UserListItem } from "./UserListItem";

interface Props extends ComponentProps<typeof List> {
  readonly checkedSet: Set<string>;
  readonly onToggleUser: ComponentProps<typeof UserListItem>["onToggleUser"];
}

export const UserChoiceList: VFC<Props> = ({ checkedSet, onToggleUser }) => {
  const [query, setQuery] = useState("");
  const {
    error: userListError,
    isLoading: isLoadingUserList,
    userList,
    abort,
  } = useUserList(query);
  const {
    error: checkedUsersError,
    isLoading: isLoadingCheckedUsers,
    users: checkedUsers,
  } = useUsers(Array.from(checkedSet));
  const onChange = useCallback<
    NonNullable<ComponentProps<typeof TextField>["onChange"]>
  >(
    (event) => {
      if (isLoadingUserList) abort();
      setQuery(event.target.value);
    },
    [isLoadingUserList, abort, setQuery]
  );

  return (
    <Box
      sx={{
        padding: 4,
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        backgroundColor: theme.palette.background.paper,
        height: "100%",
        overflowY: "scroll",
      }}
    >
      <Box sx={{ mt: 1, width: "100%" }}>
        <TextField
          margin="none"
          fullWidth
          id="query"
          label="Query"
          name="query"
          autoFocus
          onChange={onChange}
        />
      </Box>
      <Box sx={{ width: "100%", mt: 4 }}>
        {userListError !== null || checkedUsersError !== null ? (
          <Alert
            variant="outlined"
            severity="error"
            sx={{
              marginTop: 4,
              width: "100%",
            }}
          >
            {(userListError ?? checkedUsersError)?.message}
          </Alert>
        ) : isLoadingUserList ? (
          <Box
            sx={{
              pt: 16,
              display: "flex",
              justifyContent: "center",
              alignItems: "center",
            }}
          >
            <CircularProgress />
          </Box>
        ) : query === "" ? (
          isLoadingCheckedUsers ? (
            <Box
              sx={{
                pt: 16,
                display: "flex",
                justifyContent: "center",
                alignItems: "center",
              }}
            >
              <CircularProgress />
            </Box>
          ) : (
            <List>
              {checkedUsers.map((user) => (
                <UserListItem
                  key={user.id}
                  checked={checkedSet.has(user.id)}
                  user={user}
                  onToggleUser={onToggleUser}
                />
              ))}
            </List>
          )
        ) : userList.length === 0 ? (
          "該当するユーザーが見つかりませんでした"
        ) : (
          <List>
            {userList.map((user) => (
              <UserListItem
                key={user.id}
                checked={checkedSet.has(user.id)}
                user={user}
                onToggleUser={onToggleUser}
              />
            ))}
          </List>
        )}
      </Box>
    </Box>
  );
};
