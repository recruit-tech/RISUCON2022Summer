import List from "@mui/material/List";
import { type ComponentProps, type VFC } from "react";
import { User } from "../../../apiClient/__generated__";
import { UserListItem } from "./UserListItem";

interface Props extends ComponentProps<typeof List> {
  readonly users: User[];
}

export const UserList: VFC<Props> = ({ users, ...props }) => (
  <List {...props}>
    {users.map((user) => (
      <UserListItem key={user.id} user={user} />
    ))}
  </List>
);
