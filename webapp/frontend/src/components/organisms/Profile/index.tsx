import EditIcon from "@mui/icons-material/Edit";
import EmailOutlinedIcon from "@mui/icons-material/EmailOutlined";
import Box from "@mui/material/Box";
import IconButton from "@mui/material/IconButton";
import Link from "@mui/material/Link";
import Typography from "@mui/material/Typography";
import NextLink from "next/link";
import {
  type ComponentProps,
  type ComponentPropsWithoutRef,
  type VFC,
} from "react";
import { type User } from "../../../apiClient/__generated__";
import { Avatar } from "../../atoms/Avatar";
import { Calendar } from "../../molecules/Calendar";

type Props = ComponentPropsWithoutRef<typeof Box> & {
  readonly user: User;
} & (
    | {
        readonly editable?: false;
        readonly editPagePath?: undefined;
      }
    | {
        readonly editable: true;
        readonly editPagePath: ComponentProps<typeof NextLink>["href"];
      }
  );

export const Profile: VFC<Props> = ({
  user,
  editable,
  editPagePath,
  ...props
}) => (
  <Box {...props}>
    <Box
      sx={{
        display: "flex",
        justifyContent: "flex-start",
        alignItems: "center",
        pt: 8,
        pl: 2,
      }}
    >
      <Avatar
        userId={user.id}
        name={user.name}
        sx={{ width: 56, height: 56 }}
      />
      <Box sx={{ ml: 4 }}>
        <Typography sx={{ display: "flex", alignItems: "center" }}>
          <Typography sx={{ fontSize: 24 }} color="text.primary">
            {user.name}
          </Typography>
          {editable ? (
            // eslint-disable-next-line @next/next/link-passhref
            <NextLink href={editPagePath}>
              <IconButton color="default" aria-label="edit">
                <EditIcon />
              </IconButton>
            </NextLink>
          ) : null}
        </Typography>
        <Typography
          sx={{
            fontSize: 16,
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
          }}
          color="text.secondary"
        >
          <EmailOutlinedIcon fontSize="small" sx={{ mr: 1 }} />
          <address>
            <Link href={`mailto:${user.email}`} underline="hover">
              <code>{user.email}</code>
            </Link>
          </address>
        </Typography>
      </Box>
    </Box>
    <Box
      sx={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        pt: 8,
      }}
    >
      <Calendar userId={user.id} />
    </Box>
  </Box>
);
