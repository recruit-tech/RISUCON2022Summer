import AddIcon from "@mui/icons-material/Add";
import LogoutIcon from "@mui/icons-material/Logout";
import PersonOutlineIcon from "@mui/icons-material/PersonOutline";
import SearchIcon from "@mui/icons-material/Search";
import AppBar from "@mui/material/AppBar";
import Box from "@mui/material/Box";
import IconButton from "@mui/material/IconButton";
import ListItemIcon, { ListItemIconProps } from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import Menu from "@mui/material/Menu";
import MenuItem from "@mui/material/MenuItem";
import Toolbar from "@mui/material/Toolbar";
import Tooltip from "@mui/material/Tooltip";
import Typography from "@mui/material/Typography";
import Link, { LinkProps } from "next/link";
import { useRouter } from "next/router";
import {
  useCallback,
  useState,
  type ComponentPropsWithoutRef,
  type MouseEvent,
  type VFC
} from "react";
import { api } from "../../../apiClient";
import { User } from "../../../apiClient/__generated__";
import { useCurrentUserRefresher } from "../../../states/currentUserState";
import { Avatar } from "../../atoms/Avatar";

interface Props extends ComponentPropsWithoutRef<typeof AppBar> {
  readonly currentUser: User;
}

type MenuItem = {
  readonly title: string;
  readonly icon: ListItemIconProps["children"];
} & (
  | { readonly type: "link"; readonly href: LinkProps["href"] }
  | { readonly type: "action"; readonly onClick: (event: MouseEvent) => void }
);

const MENU_ITEMS: MenuItem[] = [
  {
    type: "link",
    title: "マイページ",
    icon: <PersonOutlineIcon fontSize="small" />,
    href: "/",
  },
  {
    type: "link",
    title: "ユーザー検索",
    icon: <SearchIcon fontSize="small" />,
    href: "/user/search",
  },
  {
    type: "link",
    title: "新規スケジュール作成",
    icon: <AddIcon fontSize="small" />,
    href: "/schedule/new",
  },
  // logout action is added in NotLoggedInAppBar to use hooks
];

export const LoggedInAppBar: VFC<Props> = ({ currentUser, ...props }) => {
  const [anchorEl, setAnchorEl] = useState<
    MouseEvent<HTMLButtonElement>["currentTarget"] | null
  >(null);

  const onAvatarClick = useCallback(
    (event: MouseEvent<HTMLButtonElement>) => {
      setAnchorEl(event.currentTarget);
    },
    [setAnchorEl]
  );

  const onCloseMenu = useCallback(() => {
    setAnchorEl(null);
  }, [setAnchorEl]);

  const router = useRouter();
  const refreshCurrentUser = useCurrentUserRefresher();
  const onClickLogoutMenuItem = useCallback(() => {
    onCloseMenu();
    api.postLogout().then(() => {
      refreshCurrentUser();
      router.push("/login");
    });
  }, [onCloseMenu, refreshCurrentUser, router]);

  return (
    <Box>
      <AppBar color="inherit" {...props}>
        <Toolbar sx={{ display: "flex", justifyContent: "space-between" }}>
          <Link href="/" passHref>
            <a style={{ textDecoration: "none" }}>
              <Typography
                color="WindowText"
                variant="h6"
                component="h1"
                sx={{ flexGrow: 1 }}
              >
                R-Calendar
              </Typography>
            </a>
          </Link>

          <Box sx={{ flexGrow: 0 }}>
            <Tooltip title="メニューを開く">
              <IconButton onClick={onAvatarClick} sx={{ p: 0 }}>
                <Avatar userId={currentUser.id} name={currentUser.name} />
              </IconButton>
            </Tooltip>
            <Menu
              sx={{ mt: "45px" }}
              id="menu-appbar"
              anchorEl={anchorEl}
              anchorOrigin={{
                vertical: "top",
                horizontal: "right",
              }}
              keepMounted
              transformOrigin={{
                vertical: "top",
                horizontal: "right",
              }}
              open={anchorEl !== null}
              onClose={onCloseMenu}
            >
              {[
                ...MENU_ITEMS,
                {
                  type: "action",
                  title: "ログアウト",
                  icon: <LogoutIcon fontSize="small" />,
                  onClick: onClickLogoutMenuItem,
                } as MenuItem,
              ].map((item) =>
                item.type === "link" ? (
                  // eslint-disable-next-line @next/next/link-passhref
                  <Link key={item.title} href={item.href}>
                    <MenuItem onClick={onCloseMenu}>
                      <ListItemIcon>{item.icon}</ListItemIcon>
                      <ListItemText>{item.title}</ListItemText>
                    </MenuItem>
                  </Link>
                ) : item.type === "action" ? (
                  <MenuItem key={item.title} onClick={item.onClick}>
                    <ListItemIcon>{item.icon}</ListItemIcon>
                    <ListItemText>{item.title}</ListItemText>
                  </MenuItem>
                ) : null
              )}
            </Menu>
          </Box>
        </Toolbar>
      </AppBar>
    </Box>
  );
};
