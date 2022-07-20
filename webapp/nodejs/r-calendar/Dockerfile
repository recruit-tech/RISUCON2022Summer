FROM node:18.2.0-alpine

RUN apk update && apk add git bash file mysql-client

WORKDIR /webapp/nodejs/r-calendar

COPY nodejs/r-calendar/package.json nodejs/r-calendar/package-lock.json ./

RUN npm ci

COPY nodejs/r-calendar/app.ts nodejs/r-calendar/time_in.js nodejs/r-calendar/tsconfig.json ./

EXPOSE 3000

COPY sql ../../sql

COPY scripts/wait-for-db.sh ./

CMD ["./wait-for-db.sh", "npm run dev"]
