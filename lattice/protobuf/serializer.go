package protobuf

import (
	"github.com/bufbuild/protocompile/parser"
	"github.com/bufbuild/protocompile/reporter"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	pref "google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
	"io"
)

// MakeFileDescriptor 生成proto的文件描述
//
// Parameters:
//   - reader io.Reader
//
// Returns:
//   - pref.FileDescriptor
func MakeFileDescriptor(reader io.Reader) pref.FileDescriptor {
	errHandler := reporter.NewHandler(nil)
	ast, err := parser.Parse("example.proto", reader, errHandler)
	if err != nil {
		panic(err)
	}

	result, err := parser.ResultFromAST(ast, true, errHandler)
	if err != nil {
		panic(err)
	}

	fdp := result.FileDescriptorProto()

	// get FileDescriptor
	fd, err := protodesc.NewFile(fdp, nil)
	if err != nil {
		panic(err)
	}
	return fd
}

// MarshallMessage 序列化
//
// Parameters:
//   - fd pref.FileDescriptor
//   - json string
//
// Returns:
//   - []byte
//   - error
func MarshallMessage(fd pref.FileDescriptor, json string) ([]byte, error) {
	messageDescriptor := fd.Messages().Get(0)
	// messageDescriptor := fd.Messages().ByName(pref.Name(name))
	message := dynamicpb.NewMessage(messageDescriptor)

	if err := protojson.Unmarshal([]byte(json), message); err != nil {
		return nil, err
	}

	bytes, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	return bytes, err
}

// UnmarshallMessage 反序列化
//
// Parameters:
//   - fd pref.FileDescriptor
//   - data []byte
//
// Returns:
//   - string
//   - error
func UnmarshallMessage(fd pref.FileDescriptor, data []byte) (string, error) {
	messageDescriptor := fd.Messages().Get(0)
	// messageDescriptor := fd.Messages().ByName(pref.Name(name))
	message := dynamicpb.NewMessage(messageDescriptor)

	err := proto.Unmarshal(data, message)
	if err != nil {
		return "", err
	}

	jsonBytes, err := protojson.Marshal(message)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
