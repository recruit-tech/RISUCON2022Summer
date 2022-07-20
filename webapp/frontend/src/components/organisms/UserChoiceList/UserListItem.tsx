import ListItem from "@mui/material/ListItem";
import ListItemButton from "@mui/material/ListItemButton";
import { useCallback, type ComponentProps, type VFC } from "react";
import { User } from "../../../apiClient/__generated__";
import { UserCard } from "./UserCard";

interface Props extends ComponentProps<typeof ListItem> {
  readonly checked: boolean;
  readonly user: User;
  readonly onToggleUser: (userId: string, toBeChecked: boolean) => void;
}

export const UserListItem: VFC<Props> = ({
  checked,
  user,
  onToggleUser,
  ...props
}) => {
  const onToggleCheck = useCallback(
    (toBeChecked) => {
      onToggleUser(user.id, toBeChecked);
    },
    [onToggleUser, user.id]
  );

  return (
    <ListItem divider disablePadding {...props}>
      <ListItemButton disableGutters sx={{ py: 0 }}>
        <UserCard
          checked={checked}
          user={user}
          onToggleCheck={onToggleCheck}
          sx={{ width: "100%" }}
        />
      </ListItemButton>
    </ListItem>
  );
};
