FROM node:18.2.0-alpine

RUN apk update && apk add git bash mysql-client

WORKDIR /webapp/nodejs/orecoco-reserve

COPY nodejs/orecoco-reserve/package.json nodejs/orecoco-reserve/package-lock.json ./

RUN npm ci

COPY nodejs/orecoco-reserve/app.ts nodejs/orecoco-reserve/tsconfig.json ./
RUN npm run build

EXPOSE 3003

COPY sql ../../sql

COPY scripts/wait-for-db.sh ./

COPY nodejs/orecoco-reserve/meeting_room.txt ./

CMD ["./wait-for-db.sh", "npm run dev"]
