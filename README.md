# OpenApi
```shell
cd open-api
go mod init github.com/milman2/go-api

```

# asdf
```shell
# asdf 초기화
# asdf가 설치된 경우 shims와 bin을 PATH에 추가
export PATH="$HOME/.asdf/shims:$HOME/.asdf/bin:$PATH"

# asdf 초기화 스크립트 로드 (있을 경우)
if [ -f "$HOME/.asdf/asdf.sh" ]; then
  . "$HOME/.asdf/asdf.sh"
fi
if [ -f "$HOME/.asdf/completions/asdf.bash" ]; then
  . "$HOME/.asdf/completions/asdf.bash"
fi

# Go CLI 도구 공통 설치 경로 (버전과 무관하게 사용)
export GOBIN="$HOME/.go-tools/bin"
export PATH="$PATH:$GOBIN"

# 기타 개인 bin 디렉토리
export PATH="$HOME/.local/bin:$PATH"
export PATH="$HOME/bin:$PATH"

```

# [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)
```shell
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
```

# generate api
```shell
mkdir -p pkg/api
oapi-codegen --config=oapi-codegen.yaml openapi.yaml
# OpenAPI 검증 도구 사용
npx @redocly/cli lint open-api/openapi.yaml

#
npm install -g swagger-cli
swagger-cli validate open-api/openapi.yaml
```


# Run
```shell
go mod tidy
go run .
curl http://localhost:8080/devices/01KATMQXHNQSCYVZEHD1ABX3BJ

curl http://localhost:8080/swagger/swagger/
curl http://localhost:8080/swagger/openapi.yaml
```

# [openapi-generator-cli](https://github.com/OpenAPITools/openapi-generator-cli)
```shell
# install java
sudo apt update && sudo apt install -y default-jdk

npx openapi-generator-cli generate -i openapi.yaml -o client -g typescript-fetch
```