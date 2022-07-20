import { CacheProvider, type EmotionCache } from "@emotion/react";
import NoSsr from "@mui/base/NoSsr";
import Container from "@mui/material/Container";
import CssBaseline from "@mui/material/CssBaseline";
import { ThemeProvider } from "@mui/material/styles";
import type { NextPage } from "next";
import type { AppProps as NextAppProps } from "next/app";
import Head from "next/head";
import { useRouter } from "next/router";
import { useEffect, type PropsWithChildren } from "react";
import { RecoilRoot } from "recoil";
import { LoggedInAppBar } from "../src/components/molecules/LoggedInAppBar";
import { NotLoggedInAppBar } from "../src/components/molecules/NotLoggedInAppBar";
import { Loading } from "../src/components/templates/Loading";
import { useCurrentUser } from "../src/states/currentUserState";
import { theme } from "../src/styles/theme";
import { createEmotionCache } from "../src/utils/createEmotionCache";

const clientSideEmotionCache = createEmotionCache();

export type Page = NextPage & {
  /**
   * if specify true, this page will be full-SSR and skip checking logged in
   */
  skipCurrentUserChecking?: boolean;
};

interface AppProps extends NextAppProps {
  readonly emotionCache?: EmotionCache;
  readonly Component: Page;
}

function AppBar({ children }: PropsWithChildren<{}>) {
  const { error, isLoading, currentUser } = useCurrentUser();
  const router = useRouter();
  useEffect(() => {
    if (error !== null) {
      router.push("/login");
    }
  }, [error, router]);

  if (isLoading || error !== null) return <Loading />;
  return (
    <>
      <LoggedInAppBar currentUser={currentUser} />
      {children}
    </>
  );
}

function App({
  Component,
  emotionCache = clientSideEmotionCache,
  pageProps,
}: AppProps) {
  const { skipCurrentUserChecking = false } = Component;

  return (
    <>
      <Head>
        <title>R-Calendar</title>
        <meta
          name="viewport"
          content="minimum-scale=1, initial-scale=1, width=device-width"
        />
      </Head>
      <CacheProvider value={emotionCache}>
        <ThemeProvider theme={theme}>
          <CssBaseline />
          <Container
            component="main"
            disableGutters
            maxWidth="md"
            sx={{ pt: 8, paddingX: 4, minHeight: "100vh" }}
          >
            {skipCurrentUserChecking ? (
              <>
                <NotLoggedInAppBar />
                <Component {...pageProps} />
              </>
            ) : (
              <NoSsr fallback={<Loading />}>
                <RecoilRoot>
                  <AppBar>
                    <Component {...pageProps} />
                  </AppBar>
                </RecoilRoot>
              </NoSsr>
            )}
          </Container>
        </ThemeProvider>
      </CacheProvider>
    </>
  );
}

export default App;
