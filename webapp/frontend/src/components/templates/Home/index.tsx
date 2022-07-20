import Error from "next/error";
import { useCurrentUser } from "../../../states/currentUserState";
import { Profile } from "../../organisms/Profile";
import { Loading } from "../Loading";

export const Home = () => {
  const { error, isLoading, currentUser } = useCurrentUser();

  if (error !== null) return <Error statusCode={0} />;

  if (isLoading) return <Loading />;

  return <Profile user={currentUser} editable editPagePath="/me/edit" />;
};
