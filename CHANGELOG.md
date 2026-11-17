# Changelog

## v0.0.5
_发布日期：2026-10-23_

- feat: harden event serialization and document event types
- feat: 新增语义事件类型 Expired/Busy/Renewed/Lost/WaitChanged 及 String

## v0.0.4
_发布日期：2026-10-12_

- build: 依赖改为已发布的 go-storage v0.0.3，移除本地 replace
- 更新 go-events

## v0.0.3
_发布日期：2023-08-12_

- Merge pull request #3 from storage-lock/dev
- feat: Event Fork的时候继承OwnerId和WatchDogId

## v0.0.2
_发布日期：2023-08-12_

- Merge pull request #2 from storage-lock/dev
- Merge pull request #1 from storage-lock/dev
- ci: GitHub action修改为只在push时执行
- add Eventable interface
- ci: add github ci flow

## v0.0.1
_发布日期：2023-08-07_

- go mod tidy

## Unreleased

- build: 升级依赖到 go-storage/go-storage-lock v0.1.0
