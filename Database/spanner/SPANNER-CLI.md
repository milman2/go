# Docker emulator
```shell
docker pull gcr.io/cloud-spanner-emulator/emulator
docker run -p 9010:9010 -p 9020:9020 gcr.io/cloud-spanner-emulator/emulator
```

# emulator 사용
```shell
export SPANNER_EMULATOR_HOST=localhost:9020
```

# 사전 준비(설치)
```shell
gcloud components install cloud-spanner-emulator
gcloud components update
```

# instance/database가 아직 없다면 생성
```shell
#gcloud auth revoke milman2@voyagergames.gg
#gcloud auth application-default revoke

# emulator 구성
gcloud config configurations create emulator
gcloud config set auth/disable_credentials true
gcloud config set project test-project # project
#gcloud config set account myuser@example.com
gcloud config set api_endpoint_overrides/spanner http://localhost:9020/

# instance
gcloud spanner instances create test-instance --config=emulator-config --description="Test Instance" --nodes=1 --project=test-project
gcloud spanner instances list
gcloud spanner instances delete test-instance2 --project=test-project

# database
gcloud spanner databases create test-db --instance=test-instance
gcloud spanner databases list --instance=test-instance
gcloud spanner databases delete test-db2 --instance=test-instance
```

# 에뮬레이터와 기본 구성 간에 전환
```shell
gcloud config list
gcloud config list project
gcloud config configurations list
#gcloud projects list # 실제 GCP 프로젝트만 보여준다. 

gcloud config configurations activate [emulator | default]
```

# Spanner CLI로 Emulator 연결하기
```shell
# 방법 1
docker exec -it school-live-api-spanner-cli-1 /bin/bash
spanner-cli -p test-project -i test-instance -d test-db
spanner-cli -p school-live-local -i school-app-instance -d school-app

# 방법 2
# docker exec -it <spanner-cli-container> spanner-cli -p <project> -i <instance> -d <database>
docker exec -it school-live-api-spanner-cli-1 spanner-cli \
  -p test-project \
  -i test-instance \
  -d test-db

curl http://localhost:9020/v1/projects/test-project/instances
curl http://localhost:9020/v1/projects/test-project/instances/test-instance/databases
curl http://localhost:9020/v1/projects/test-project/instances/test-instance/databases/test-db/sessions
```

# Query
```sql
-- 현재 데이터베이스에 존재하는 테이블 목록
SELECT table_name
FROM information_schema.tables
WHERE table_catalog = '' AND table_schema = '';

-- 현재 데이터베이스에 존재하는 인덱스 목록
SELECT index_name, table_name
FROM information_schema.indexes;

-- 현재 데이터베이스에 존재하는 컬럼 목록
SELECT table_name, column_name, spanner_type
FROM information_schema.columns;
```

# DDL example
```sql
CREATE TABLE Users (
  UserId STRING(36) NOT NULL,
  Name   STRING(1024),
  Age    INT64,
) PRIMARY KEY(UserId);

INSERT INTO Users (UserId, Name, Age)
VALUES ("u1", "Alice", 20);

SELECT * FROM Users;

-- schema/metadata 조회
SELECT t.table_name
FROM information_schema.tables AS t
WHERE t.table_catalog = '' AND t.table_schema = '';

UPDATE Users SET Age = 21 WHERE UserId = "u1";
DELETE FROM Users WHERE UserId = "u1";
```

# 실제 Cloud Spanner로 전환
```shell
# emulator 환견 변수 제거
unset SPANNER_EMULATOR_HOST
# 올바른 프로젝트 설정
gcloud config set project <실제-프로젝트-ID>
# 계정 인증
gcloud auth login
# 또는 Application Default Credentials 설정
gcloud auth application-default login

# instance/db 확인
gcloud spanner instances list
gcloud spanner databases list --instance=<인스턴스-ID>
```