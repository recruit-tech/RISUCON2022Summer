import AppBar from "@mui/material/AppBar";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";
import Link from "next/link";
import { type ComponentPropsWithoutRef, type VFC } from "react";

interface Props extends ComponentPropsWithoutRef<typeof AppBar> {}

export const NotLoggedInAppBar: VFC<Props> = (props) => (
  <Box sx={{ flexGrow: 1 }}>
    <AppBar color="inherit" {...props}>
      <Toolbar>
        <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
          R-Calendar
        </Typography>
        {/* eslint-disable-next-line @next/next/link-passhref */}
        <Link href="/login">
          <Button color="primary">ログイン</Button>
        </Link>
        {/* eslint-disable-next-line @next/next/link-passhref */}
        <Link href="/signup">
          <Button sx={{ ml: 2 }} color="primary">
            新規ユーザー作成
          </Button>
        </Link>
      </Toolbar>
    </AppBar>
  </Box>
);
