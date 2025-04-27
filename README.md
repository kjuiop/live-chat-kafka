# live-chat-kafka

---

`live-chat-kafka`는 **실시간 채팅 시스템**을 위한 서버 애플리케이션으로, Kafka와 WebSocket을 활용하여 확장성 있는 채팅 서비스 구축을 목표로 합니다.


## 🧩 시스템 구성

![image](https://github.com/user-attachments/assets/b64af16b-e320-49c9-8187-6cadc1b12c3c)


- 사용자는 HTTP API를 통해 채팅방을 생성하거나 삭제할 수 있습니다.
- 채팅방에 입장한 사용자는 WebSocket을 통해 실시간으로 메시지를 주고받을 수 있습니다.
- 서버는 Redis를 사용하여 메시지 브로커 역할을 수행하며, 채팅 메시지를 효율적으로 관리합니다.

## ⚙️ 기술 스택

- **언어**: Go 1.23
- **웹 프레임워크**: Gin
- **WebSocket**: Gorilla WebSocket
- **데이터 저장소**: Redis
- **빌드 도구**: Makefile

## 🚀 시작하기


### TEST

```shell
make test
```


### BUILD

```shell
make build
```


### RUN
```shell
./live-chat-server
```

<br />


## 📄 API Spec Document


### ws join

```shell
ws://localhost:8091/ws/chat/join/rooms/N1-01JSVD2N05RD0F4GPGDHR5C73J/user/jake
```


### ws chat

```shell
{
    "Method":"chat",
    "SendUserId": "jake",
    "Message": "hello"
}
```