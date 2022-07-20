import ListItem from "@mui/material/ListItem";
import ListItemButton from "@mui/material/ListItemButton";
import Link from "next/link";
import { type ComponentProps, type VFC } from "react";
import { User } from "../../../apiClient/__generated__";
import { UserCard } from "../../molecules/UserCard";

interface Props extends ComponentProps<typeof ListItem> {
  readonly user: User;
}

export const UserListItem: VFC<Props> = ({ user, ...props }) => (
  // eslint-disable-next-line @next/next/link-passhref
  <Link href={{ pathname: "/user/[id]", query: { id: user.id } }}>
    <ListItem divider disablePadding {...props}>
      <ListItemButton disableGutters sx={{ py: 0 }}>
        <UserCard user={user} sx={{ width: "100%" }} />
      </ListItemButton>
    </ListItem>
  </Link>
);
