# GRPC - ToyProject

- 목표
    - Golang 을 통해 개발을 진행한다.
    - GRPC Server와 Gateway를 분리하여 개발한다.
    - Oauth 인증을 통한 로그인 서비스 제공.
    - 파일 업로드 시스템 구축.

## Tech Stack

- Golang
    - webServer
        - golang-Gin
    - CLI
        - Cobra
    - Config
        - Viper
    - GRPC
        - Protobuf
        - GRPC server
        - GRPC Gateway
    - ORM
        - ENT
- DB
    - Postgresql
    - Redis

## User Strory

- [ ]  User는 외부 인증을 통한 회원가입이 가능하다.
- [ ]  User는 File 업로드가 가능하다.
- [ ]  User는 File 다운로드가 가능하다.

## 설계 및 설정

(Diagram 추가)

1. 외부와의 REST 통신을 위한 GRPC Gateway를 구축
    1. GRPC Gateway 자체를 사용하지 않고 GIn web Framework를 사용

       Reason

        1. Google GRPC Gateway와 종속성 해체.
        2. 좀 더 유연하게 gateway 구성이 가능.