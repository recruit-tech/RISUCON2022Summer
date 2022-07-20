# ローカル用のVM競技環境作成方法

Ansibleを使ってmultipass上に競技環境を立てる方法を記述します

下記の環境ができます。
- 競技1台 (CPU: 1 core, memory: 2 GB)
- ベンチ1台 (CPU: 1 core, memory: 2 GB)

想定動作環境は ArmベースのMacです。

## prerequirement

### preparation
- ssh-keygen

## installation

### install ansible
```shell
brew install ansible
# ansible [core 2.12.5]
```

### install multipass
```shell
brew install --cask multipass
# multipass   1.8.1+mac
# multipassd  1.8.1+mac
```

## quick start

```shell
cd infra/virtual-machine

## launch VM
make reset-multipass-web
make reset-multipass-bench
```

## usage

### login

```shell
#for app
multipass shell r-calendar-web

# for bench
multipass shell r-calendar-bench
```

### start bench

```shell
make start-bench
```

