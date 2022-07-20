import Alert from "@mui/material/Alert";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import { red } from "@mui/material/colors";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";
import { FormEventHandler, useCallback, useState } from "react";
import { User, UserProperties } from "../../../apiClient/__generated__";
import { getNetworkErrorMessage } from "../../../utils/getNetworkErrorMessage";
import { ValidationError } from "../../../utils/validates";

export type UserFormData = UserProperties & { readonly icon: File | null };
type MaybeUserFormData = Record<keyof UserFormData, FormDataEntryValue | null>;

function validateUserData(
  data: MaybeUserFormData
): asserts data is UserFormData {
  if (
    typeof data.name !== "string" ||
    typeof data.email !== "string" ||
    typeof data.password !== "string"
  ) {
    throw new ValidationError("入力内容が不正です");
  }

  if (data.name === "") {
    throw new ValidationError("名前は必須です");
  }

  if (data.email === "") {
    throw new ValidationError("メールアドレスは必須です");
  }

  if (data.password === "") {
    throw new ValidationError("パスワードは必須です");
  }

  if (!(data.icon instanceof File || data.icon === null)) {
    throw new ValidationError("アイコンに不正な値が指定されています");
  }
}

type Props = {
  readonly onSubmit: (userFormData: UserFormData) => Promise<void>;
} & (
  | { readonly type: "new"; readonly defaultValue?: undefined }
  | { readonly type: "edit"; readonly defaultValue: User }
);

export const UserForm = ({ onSubmit, type, defaultValue }: Props) => {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const handleSubmit = useCallback<FormEventHandler<HTMLFormElement>>(
    (event) => {
      event.preventDefault();
      const data = new FormData(event.currentTarget);
      const name = data.get("name");
      const email = data.get("email");
      const password = data.get("password");
      const icon = data.get("icon");

      try {
        const userFormData: MaybeUserFormData = {
          name,
          email,
          password,
          icon,
        };
        validateUserData(userFormData);

        setIsSubmitting(true);
        setErrorMessage(null);
        onSubmit(userFormData)
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
    [onSubmit]
  );

  return (
    <Box component="form" onSubmit={handleSubmit}>
      <Box sx={{ mt: 2 }}>
        <Typography component="h2" variant="h6">
          アイコン
        </Typography>
        <input accept="image/*" name="icon" type="file" />
      </Box>
      <Box sx={{ mt: 2 }}>
        <Typography component="h2" variant="h6">
          名前
          <Typography component="span" color={red[500]}>
            *
          </Typography>
        </Typography>
        <TextField
          name="name"
          sx={{ mt: 1 }}
          variant="filled"
          required
          fullWidth
          margin="none"
          defaultValue={defaultValue?.name}
          type="text"
          autoComplete="nickname"
          autoFocus
        />
      </Box>
      <Box sx={{ mt: 2 }}>
        <Typography component="h2" variant="h6">
          メールアドレス
          <Typography component="span" color={red[500]}>
            *
          </Typography>
        </Typography>
        <TextField
          name="email"
          sx={{ mt: 1 }}
          variant="filled"
          required
          fullWidth
          margin="none"
          defaultValue={defaultValue?.email}
          type="email"
          autoComplete="email"
          autoFocus
        />
      </Box>
      <Box sx={{ mt: 2 }}>
        <Typography component="h2" variant="h6">
          パスワード
          <Typography component="span" color={red[500]}>
            *
          </Typography>
        </Typography>
        <TextField
          name="password"
          sx={{ mt: 1 }}
          variant="filled"
          required
          fullWidth
          margin="none"
          type="password"
          autoComplete="new-password"
          autoFocus
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
