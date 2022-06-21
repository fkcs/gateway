package common

import (
	"github.com/fkcs/gateway/internal/infrastructure/logger"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"mime/multipart"
	"reflect"
)

const (
	MethodPost = "POST"
	MethodGet  = "GET"
)

type Param struct {
	Key          string
	Required     bool
	DefaultValue string
}

type ParamDoc struct {
	Param
	Doc string
}

type FileDoc struct {
	Key string
	Doc string
}

func Required(key, doc string) ParamDoc {
	return ParamDoc{
		Param: Param{
			Key:          key,
			Required:     true,
			DefaultValue: "",
		},
		Doc: doc,
	}
}

func Optional(key, defaultValue, doc string) ParamDoc {
	return ParamDoc{
		Param: Param{
			Key:          key,
			Required:     false,
			DefaultValue: defaultValue,
		},
		Doc: doc,
	}
}

type Handler struct {
	Tags          []string
	Path          string
	Method        string
	Name          string
	ReqType       reflect.Type
	RspType       reflect.Type
	ParamDocs     map[string]ParamDoc
	ImportFileDoc FileDoc
	Processor     func(map[string]string, interface{}, multipart.File) (int, interface{})
}

func MakeFileDoc(key, doc string) FileDoc {
	return FileDoc{
		Key: key,
		Doc: doc,
	}
}

func (x FileDoc) IsNil() bool {
	return x.Key == ""
}

func NewHandler(name, method string) *Handler {
	return &Handler{
		Tags:          make([]string, 0),
		Method:        method,
		Name:          name,
		ReqType:       nil,
		RspType:       nil,
		ParamDocs:     make(map[string]ParamDoc),
		ImportFileDoc: FileDoc{},
		Processor:     nil,
	}
}

func (x *Handler) WithTags(tags ...string) *Handler {
	x.Tags = append(x.Tags, tags...)
	return x
}

func (x *Handler) WithParam(paramDoc ParamDoc) *Handler {
	x.ParamDocs[paramDoc.Key] = paramDoc
	return x
}

func (x *Handler) WithImportFile(key string) *Handler {
	x.ImportFileDoc = MakeFileDoc(key, "")
	return x
}

func (x *Handler) WithReqType(typ reflect.Type) *Handler {
	x.ReqType = typ
	return x
}

func (x *Handler) WithRspType(typ reflect.Type) *Handler {
	x.RspType = typ
	return x
}

func (x *Handler) WithProcessor(p func(map[string]string, interface{}, multipart.File) (int, interface{})) *Handler {
	x.Processor = p
	return x
}

func (x *Handler) ParseRequest(ctx *fasthttp.RequestCtx) (int, interface{}) {
	params := make(map[string]string, 0)
	errMsgs := make([]string, 0)
	for key, param := range x.ParamDocs {
		if valueBytes := ctx.FormValue(key); len(valueBytes) > 0 {
			value := string(valueBytes)
			params[key] = value
		} else {
			if param.Required {
				msg := fmt.Sprintf("param key %s is missing", key)
				errMsgs = append(errMsgs, msg)
			} else {
				logger.Logger().Debugf("param key %s is missing, default to '%s'", key, param.DefaultValue)
				params[key] = param.DefaultValue
			}
		}
	}

	var file multipart.File
	if !x.ImportFileDoc.IsNil() {
		fileHeader, err := ctx.FormFile(x.ImportFileDoc.Key)
		if err != nil {
			msg := fmt.Sprintf("error fetching file %s: %s", x.ImportFileDoc.Key, err)
			errMsgs = append(errMsgs, msg)
		}
		file, err = fileHeader.Open()
		defer file.Close()
		if err != nil {
			msg := fmt.Sprintf("error fetching file %s: %s", x.ImportFileDoc.Key, err)
			errMsgs = append(errMsgs, msg)
		}
	}
	if len(errMsgs) > 0 {
		return fasthttp.StatusBadRequest, fmt.Errorf("%v", errMsgs)
	}

	var reqInterface interface{}
	if x.ReqType != nil {
		reqBytes := ctx.PostBody()
		req := reflect.New(x.ReqType)
		if err := json.Unmarshal(reqBytes, req.Interface()); err != nil {
			return fasthttp.StatusBadRequest, fmt.Errorf("%v", err)
		}
		reqInterface = req.Elem().Interface()
	}
	return x.Processor(params, reqInterface, file)
}
