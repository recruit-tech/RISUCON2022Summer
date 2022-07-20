import CloseIcon from "@mui/icons-material/Close";
import IconButton from "@mui/material/IconButton";
import { type SnackbarProps } from "@mui/material/Snackbar";
import React, { useCallback, useMemo, useState } from "react";

export function useSnackbar() {
  const [isOpen, setIsOpen] = useState(false);
  const [message, setMessage] = useState<SnackbarProps["message"]>(undefined);

  const open = useCallback(
    (message?: SnackbarProps["message"]) => {
      setIsOpen(true);
      setMessage(message);
    },
    [setIsOpen, setMessage]
  );

  const onClose = useCallback<NonNullable<SnackbarProps["onClose"]>>(
    (_, reason) => {
      if (reason === "clickaway") {
        return;
      }
      setIsOpen(false);
    },
    [setIsOpen]
  );

  const action = useMemo<NonNullable<SnackbarProps["action"]>>(
    () => (
      <IconButton
        size="small"
        aria-label="close"
        color="inherit"
        onClick={() => setIsOpen(false)}
      >
        <CloseIcon fontSize="small" />
      </IconButton>
    ),
    [setIsOpen]
  );

  return [
    open,
    {
      open: isOpen,
      onClose,
      children: message,
      action,
    } as SnackbarProps,
  ] as const;
}
