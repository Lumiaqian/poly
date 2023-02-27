# Poly

基于 [Wails2](https://wails.io/) 开发的观看纯净直播的跨平台桌面软件

## 技术栈

- 基于 [Wails2](https://wails.io/) 开发的观看纯净直播的跨平台桌面软件
- 后端技术 Go
- 前端技术
  - Quasar V2 (Vue 3)
  - Typescript
  - Quasar App CLI with Vite
  - Composition API with `<script setup>`
  - Sass with SCSS syntax
  - ESLint + Pinia
  - Prettier ESLint preset

## 特性

* 跨平台 Mac、Win11、Linux
* 没有广告、打赏等乱七八糟的东西
* 无服务端

- 支持的直播平台

  - [X] 虎牙直播
  - [ ] 斗鱼直播
  - [ ] Bilibili直播
- 后续功能

  - [ ] 关注列表持久化
  - [ ] 平台
  - [ ] 分类
  - [ ] 电视类似TV-BOX
  - [ ] ...

## 在本地运行

Clone 这个 project

```bash
  git clone https://github.com/Lumiaqian/poly.git
```

前往项目目录

```bash
  cd ploy
```

调试

```bash
  wails dev
```

执行编译

```bash
  wails build
```
