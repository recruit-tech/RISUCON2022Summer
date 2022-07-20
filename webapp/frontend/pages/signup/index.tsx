import { Signup } from "../../src/components/templates/Signup";
import { type Page } from "../_app";
(Signup as Page).skipCurrentUserChecking = true;

export default Signup;
