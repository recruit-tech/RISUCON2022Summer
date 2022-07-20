import EmailOutlinedIcon from "@mui/icons-material/EmailOutlined";
import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Checkbox from "@mui/material/Checkbox";
import Typography from "@mui/material/Typography";
import { useCallback, type ComponentPropsWithoutRef, type VFC } from "react";
import { User } from "../../../apiClient/__generated__";
import { Avatar } from "../../atoms/Avatar";

interface Props extends ComponentPropsWithoutRef<typeof Card> {
  readonly checked: boolean;
  readonly user: User;
  readonly onToggleCheck: (toBeChecked: boolean) => void;
}

export const UserCard: VFC<Props> = ({
  checked,
  user,
  sx,
  onToggleCheck,
  ...props
}) => {
  const onChange = useCallback(
    () => onToggleCheck(!checked),
    [checked, onToggleCheck]
  );

  return (
    <Card sx={{ borderRadius: 0, ...sx }} {...props}>
      <CardContent
        sx={{
          display: "flex",
          justifyContent: "flex-start",
          alignItems: "center",
        }}
      >
        <Checkbox checked={checked} onChange={onChange} />
        <Avatar sx={{ ml: 1 }} userId={user.id} name={user.name} />
        <Box sx={{ ml: 4 }}>
          <Typography sx={{ fontSize: 16 }} color="text.secondary" gutterBottom>
            {user.name}
          </Typography>
          <Typography
            sx={{
              fontSize: 14,
              display: "flex",
              justifyContent: "center",
              alignItems: "center",
            }}
            color="text.secondary"
          >
            <EmailOutlinedIcon fontSize="small" sx={{ mr: 1 }} />
            <code>{user.email}</code>
          </Typography>
        </Box>
      </CardContent>
    </Card>
  );
};
