import MuiAvatar from "@mui/material/Avatar";
import { type ComponentProps, type VFC } from "react";
import { convertIdToHexColor } from "../../../utils/convertIdToHexColor";

interface Props extends ComponentProps<typeof MuiAvatar> {
  readonly userId: string;
  readonly name: string;
}

const shortenName = (name: string) => {
  const parts = name.trim().split(/\s+/g).filter(Boolean);
  switch (parts.length) {
    case 0:
      return "??";
    case 1:
      return parts[0].substring(0, 2);
    default:
      return parts[0][0] + parts[1][0];
  }
};

export const Avatar: VFC<Props> = ({ userId, name, sx, ...props }) => {
  const bgcolor = convertIdToHexColor(userId);

  return (
    <MuiAvatar
      sx={{
        bgcolor,
        wordBreak: "keep-all",
        ...sx,
      }}
      src={`/api/user/icon/${userId}`}
      alt={name}
      {...props}
    >
      {shortenName(name)}
    </MuiAvatar>
  );
};
