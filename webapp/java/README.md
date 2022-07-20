# Java実装

author: @kawasima

## 前提条件

JDK17以上とMavenが必要です。

## 動かし方

orecoco-reserve, r-calendarともに、プロジェクト配下に移動し、

```shell
mvn compile quarkus:build
```

でビルドします。

```shell
java -jar target/quarkus-app/quarkus-run.jar
```

で起動します。

