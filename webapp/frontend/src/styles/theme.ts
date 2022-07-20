import { blue, orange, red } from "@mui/material/colors";
import { createTheme } from "@mui/material/styles";

export const theme = createTheme({
  palette: {
    primary: {
      main: blue.A400,
    },
    secondary: {
      main: orange.A200,
    },
    error: {
      main: red.A400,
    },
  },
});
