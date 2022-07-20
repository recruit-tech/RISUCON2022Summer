import http from "http";
import { NextApiRequest, NextApiResponse } from "next";
import httpProxyMiddleware from "next-http-proxy-middleware";

export const config = {
  api: {
    bodyParser: false,
  },
};

export default function proxy(
  req: NextApiRequest,
  res: NextApiResponse
): Promise<any> {
  return httpProxyMiddleware(req, res, {
    target: "http://localhost:3000",
    changeOrigin: true,
    cookieDomainRewrite: "http://localhost:3000",
    pathRewrite: [
      {
        patternStr: "^/api/",
        replaceStr: "/",
      },
    ],
    agent: new http.Agent(),
  });
}
