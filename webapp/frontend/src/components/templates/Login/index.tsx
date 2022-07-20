import Alert from "@mui/material/Alert";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";
import { useCallback, useState, type FormEventHandler, type VFC } from "react";
import { api } from "../../../apiClient";
import { getNetworkErrorMessage } from "../../../utils/getNetworkErrorMessage";

interface Props {
  readonly errorMessage: string | null;
  readonly isSubmitting: boolean;
  readonly onSubmit: FormEventHandler<HTMLFormElement>;
}

function useLogin(): Props {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [errorMessage, setErrorMessage] = useState<Props["errorMessage"]>(null);
  const onSubmit = useCallback<Props["onSubmit"]>(
    (event) => {
      event.preventDefault();
      const data = new FormData(event.currentTarget);
      const email = data.get("email");
      const password = data.get("password");

      if (typeof email !== "string" || typeof password !== "string") {
        setErrorMessage("入力内容が不正です");
        return;
      }

      if (email === "") {
        setErrorMessage("メールアドレスは必須です");
        return;
      }

      if (password === "") {
        setErrorMessage("パスワードは必須です");
        return;
      }

      setIsSubmitting(true);
      setErrorMessage(null);
      api
        .postLogin({ email, password })
        .then(() => { window.location.href = "/" })
        .catch((error) => {
          console.error(error);
          setErrorMessage(getNetworkErrorMessage(error));
        })
        .finally(() => {
          setIsSubmitting(false);
        });
    },
    [setIsSubmitting, setErrorMessage]
  );

  return {
    errorMessage,
    isSubmitting,
    onSubmit,
  };
}

const Presentation: VFC<Props> = ({ errorMessage, isSubmitting, onSubmit }) => (
  <Box
    sx={{
      marginTop: 8,
      display: "flex",
      flexDirection: "column",
      alignItems: "center",
    }}
  >
    <Typography component="h1" variant="h5">
      ログイン
    </Typography>
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
    <Box
      component="form"
      onSubmit={onSubmit}
      noValidate
      sx={{ mt: 1, width: "100%" }}
    >
      <TextField
        margin="normal"
        required
        fullWidth
        id="email"
        label="メールアドレス"
        name="email"
        type="email"
        autoComplete="email"
        autoFocus
      />
      <TextField
        margin="normal"
        required
        fullWidth
        name="password"
        label="パスワード"
        type="password"
        id="password"
        autoComplete="current-password"
      />
      <Button
        type="submit"
        disabled={isSubmitting}
        fullWidth
        variant="contained"
        sx={{ mt: 3, mb: 2 }}
      >
        ログイン
      </Button>
    </Box>
  </Box>
);

export const Login = () => {
  const props = useLogin();
  return <Presentation {...props} />;
};
