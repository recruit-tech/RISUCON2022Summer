import Error from "next/error";
import { useRouter } from "next/router";
import { useUser } from "../../../states/userState";
import { getQuery } from "../../../utils/getQuery";
import { Profile } from "../../organisms/Profile";
import { Loading } from "../Loading";

export const UserInfo = () => {
  const router = useRouter();
  const userId = getQuery(router, "id");
  const { error, isLoading, user } = useUser(userId);

  if (userId == undefined || isLoading) return <Loading />;

  if (error !== null) return <Error statusCode={0} title={error.message} />;

  return <Profile user={user} />;
};
