@echo off
cd E:\go_example\wsprotGame
protoc --go_out=./proto/gen --proto_path=./proto/intr ./proto/intr/*.proto
pause

