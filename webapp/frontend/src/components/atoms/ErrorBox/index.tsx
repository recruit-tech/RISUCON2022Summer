import ErrorIcon from "@mui/icons-material/ErrorOutline";
import Box from "@mui/material/Box";
import { red } from "@mui/material/colors";
import Tooltip from "@mui/material/Tooltip";
import { type ComponentProps, type VFC } from "react";

interface Props extends ComponentProps<typeof Box> {
  readonly error: Error;
}

export const ErrorBox: VFC<Props> = ({ error, sx, ...props }) => (
  <Tooltip title={error.message} arrow followCursor placement="top">
    <Box
      sx={{
        width: "100%",
        height: "100%",
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        backgroundColor: red[50],
        ...sx,
      }}
      {...props}
    >
      <ErrorIcon color="error" />
    </Box>
  </Tooltip>
);
