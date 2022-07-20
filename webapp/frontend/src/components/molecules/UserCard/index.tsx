import EmailOutlinedIcon from "@mui/icons-material/EmailOutlined";
import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";
import { type ComponentPropsWithoutRef, type VFC } from "react";
import { User } from "../../../apiClient/__generated__";
import { Avatar } from "../../atoms/Avatar";

interface Props extends ComponentPropsWithoutRef<typeof Card> {
  readonly user: User;
}

export const UserCard: VFC<Props> = ({ user, sx, ...props }) => (
  <Card sx={{ borderRadius: 0, ...sx }} {...props}>
    <CardContent
      sx={{
        display: "flex",
        justifyContent: "flex-start",
        alignItems: "center",
      }}
    >
      <Avatar userId={user.id} name={user.name} />
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
