import CircularProgress from "@mui/material/CircularProgress";
import Container from "@mui/material/Container";

export const Loading = () => (
  <Container
    disableGutters
    sx={{
      width: "auto",
      position: "absolute",
      top: "50%",
      left: "50%",
      transform: "translate(-50%, -50%)",
    }}
  >
    <CircularProgress />
  </Container>
);
