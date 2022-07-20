import Alert from "@mui/material/Alert";
import Box from "@mui/material/Box";
import CircularProgress from "@mui/material/CircularProgress";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";
import { useCallback, useState, type ComponentProps } from "react";
import { useUserList } from "../../../states/userListState";
import { UserList } from "../../organisms/UserList";

export const UserSearch = () => {
  const [query, setQuery] = useState("");
  const { error, isLoading, userList, abort } = useUserList(query);
  const onChange = useCallback<
    NonNullable<ComponentProps<typeof TextField>["onChange"]>
  >(
    (event) => {
      if (isLoading) abort();
      setQuery(event.target.value);
    },
    [isLoading, abort, setQuery]
  );

  return (
    <Box
      sx={{
        marginTop: 8,
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
      }}
    >
      <Typography component="h1" variant="h5">
        User Search
      </Typography>
      <Box sx={{ mt: 1, width: "100%" }}>
        <TextField
          margin="normal"
          fullWidth
          id="query"
          label="検索ワード"
          name="query"
          autoFocus
          onChange={onChange}
        />
      </Box>
      <Box sx={{ width: "100%" }}>
        {error !== null ? (
          <Alert
            variant="outlined"
            severity="error"
            sx={{
              marginTop: 4,
              width: "100%",
            }}
          >
            {error.toString()}
          </Alert>
        ) : isLoading ? (
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
        ) : query === "" ? null : userList.length === 0 ? (
          "該当するユーザーが見つかりませんでした"
        ) : (
          <UserList users={userList} />
        )}
      </Box>
    </Box>
  );
};
