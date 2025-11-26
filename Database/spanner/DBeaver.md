# env variable
- export SPANNER_EMULATOR_HOST=192.168.50.135:9010

# 실제 cloud spanner(not emulator)를 사용하는 경우
gcloud auth application-default login

# New Connection
- Connection Settings 
    - Main - [Driver Settings]
        - URL Template: jdbc:cloudspanner:/projects/school-live-local/instances/school-app-instance/databases/school-app;usePlainText=true
    - Driver properties
        - useplaintext: true

# Community 버전이 Driver 버전이 낮아서 안 되는 듯.
