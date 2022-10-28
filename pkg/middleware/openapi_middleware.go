package middleware

import (
	"fmt"
	"io"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	legacyrouter "github.com/getkin/kin-openapi/routers/legacy"
	"github.com/gofiber/fiber/v2"
)

func yamlBodyDecoder(r io.Reader, h http.Header, s *openapi3.SchemaRef, fn openapi3filter.EncodingFn) (interface{}, error) {
	return "", nil
}

func OpenapiInputValidator(openApiFile string) func(c *fiber.Ctx) {
	return func(c *fiber.Ctx) {
		openapi3filter.RegisterBodyDecoder("text/yaml", yamlBodyDecoder)
		ctx := c.Context()

		loader := openapi3.Loader{Context: ctx}
		doc, err := loader.LoadFromFile(openApiFile)
		if err != nil {
			fmt.Println(err.Error())
		}

		err = doc.Validate(ctx)
		if err != nil {
			fmt.Println(err.Error())
		}

		router, err := legacyrouter.NewRouter(doc) //.WithSwaggerFromFile(openapiFile)
		if err != nil {
			fmt.Println(err.Error())
		}
		httpReq := &http.Request{}
		route, pathParams, _ := router.FindRoute(&http.Request{})

		// Validate Request
		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    httpReq,
			PathParams: pathParams,
			Route:      route,
		}

		if erro := openapi3filter.ValidateRequest(ctx, requestValidationInput); erro != nil {
			// panic(err)
			fmt.Println("Fail")
		}

		c.Next()
	}
}
