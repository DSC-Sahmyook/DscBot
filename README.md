# DscBot
DSC-Sahmyook 내 에서 프로젝트를 진행 할 때 프로젝트 관리를 도와줄 봇을 만드는 프로젝트

# 변수명 (계속 수정중) 좋은 의견 받는 중
1. 기본적으로 변수명은 camelcase(낙타형)으로 한다. 예) messageCreate => 띄어쓰기가 필요한 곳에 대문자를 이용한다.
2. model의 structure 이름은 복수형 명사를 이용한다. 예) Messages, Users, Todos
3. 관계형 structure 이름은 Rel + a테이블명 + b테이블명으로 한다. 예) RelChannelUser, RelProjectUser 등
4. controller 함수명은 명사 + 동사로 한다. 예)MessageCreator, MessageNotifier, ProjectUpdater 등

# 나눈 폴더에 따라 작업을 진행 하세요.
## config
프로젝트를 위한 설정 파일이 들어갈 폴더   
예) db 설정, 토큰 설정 등

## api
트렐로 api 통신을 할 때 작성한 파일이 들어갈 폴더   
api로 부터 정보를 받는 함수를 이곳에 작성합니다.

## model
정보를 받아오거나 db와 통신을 할 때 structure가 필요 합니다.   
구조체를 작성할 때 와 DB로 부터 데이터를 가져오는 코드를 작성할 때 이곳에 작성합니다.

## controller
디스코드와 통신하여 실행할 함수를 만든 파일이 들어갈 폴더   
messageCreate와 같이 디스코드 요청에 대한 응답에 직접적으로 연관이 있는 함수들을 모아 둡니다.

