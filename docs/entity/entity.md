## 简介

> 物联网世界里的操作对象，以及这些对象组合抽象出来的对象，包括网关，设备，设备的聚合抽象等等。


- 实体具有零个或多个属性。
- 实体属性的类型是任意类型的。
- 实体应当具有如下必须字段：
    - entity_id: 实体id
    - tag：实体标签
    - user_id: 实体拥有者
    - source：实体的创建插件
    - version：实体的版本




实体区分转台存储和运行时实体actor，实体状态在actor维护和更新，并采用一定策略落盘。
![entity-store](../images/entity-store.png)


















