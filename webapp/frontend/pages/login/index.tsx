import { Login } from "../../src/components/templates/Login";
import { type Page } from "../_app";
(Login as Page).skipCurrentUserChecking = true;

export default Login;
