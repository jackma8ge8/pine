protoc --go_out=proto --proto_path=proto proto/*.proto  --experimental_allow_proto3_optional

# ts-node ./proto/tools/tool.ts server > ./proto/jsonproto/server.json
# ts-node ./proto/tools/tool.ts client > ./proto/jsonproto/client.json