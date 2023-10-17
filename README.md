# blubooks-cli
gin --appPort 4080 --port 3000 --all run main.go 


npx -p typescript tsc node_modules/pagedjs/src/*.js --declaration --allowJs --emitDeclarationOnly --outDir node_modules/pagedjs/src
npx -p typescript tsc node_modules/pagedjs/src/**/*.js --declaration --allowJs --emitDeclarationOnly --outDir node_modules/pagedjs/src


npx -p typescript tsc node_modules/pagedjs/src/**/*.js --declaration --allowJs --emitDeclarationOnly --outDir src/types/pagedjs


protoc --go_out=. --go_opt=paths=source_relative   --go-grpc_out=. --go-grpc_opt=paths=source_relative  person.proto
npx protoc --ts_out src/models --proto_path src/proto  person.proto

 npx protoc --ts_out src/models --proto_path ../../internal/app/  person.proto


protoc --js_out=library=myproto_libs,binary:. person.proto

 protoc person.proto \
    --proto_path=. \
    --js_out=import_style=commonjs:. \
    --grpc-web_out=import_style=typescript,mode=grpcwebtext:.


 protoc person.proto \
    --proto_path=. \
    --js_out=import_style=commonjs:. \
    --protobuf-ts=import_style=protobuf-ts,mode=grpcwebtext:.



npx protoc --ts_out . --proto_path protos  person.proto

npm install -g  protoc-gen-grpc-web