// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/items": {
            "get": {
                "description": "Gets the most popular items",
                "produces": [
                    "application/json"
                ],
                "summary": "Get most popular items",
                "operationId": "get-popular-items",
                "parameters": [
                    {
                        "type": "string",
                        "description": "number of results size",
                        "name": "size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Items returned",
                        "schema": {
                            "$ref": "#/definitions/model.Items"
                        }
                    },
                    "400": {
                        "description": "Invalid payload",
                        "schema": {
                            "$ref": "#/definitions/model.ErrResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/model.ErrResponse"
                        }
                    }
                }
            }
        },
        "/api/items/user/{user}": {
            "get": {
                "description": "Gets user items",
                "produces": [
                    "application/json"
                ],
                "summary": "Get User Items",
                "operationId": "get-user-item",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user ID",
                        "name": "user",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Items returned",
                        "schema": {
                            "$ref": "#/definitions/model.Items"
                        }
                    },
                    "400": {
                        "description": "Invalid payload",
                        "schema": {
                            "$ref": "#/definitions/model.ErrResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/model.ErrResponse"
                        }
                    }
                }
            }
        },
        "/api/probability/user/{user}/item/{item}": {
            "get": {
                "description": "Gets probability for a given user and item",
                "produces": [
                    "application/json"
                ],
                "summary": "Get probability",
                "operationId": "get-probability",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user ID",
                        "name": "user",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "item ID",
                        "name": "item",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ItemProbability returned",
                        "schema": {
                            "$ref": "#/definitions/model.ItemProbability"
                        }
                    },
                    "400": {
                        "description": "Invalid payload",
                        "schema": {
                            "$ref": "#/definitions/model.ErrResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/model.ErrResponse"
                        }
                    }
                }
            }
        },
        "/api/rate": {
            "post": {
                "description": "Adds rating for a given user with an item",
                "produces": [
                    "application/json"
                ],
                "summary": "Create rating for a gien user with an item",
                "operationId": "post-rate",
                "parameters": [
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Rate"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Rating created",
                        "schema": {
                            "$ref": "#/definitions/model.Rate"
                        }
                    },
                    "400": {
                        "description": "Invalid payload",
                        "schema": {
                            "$ref": "#/definitions/model.ErrResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/model.ErrResponse"
                        }
                    }
                }
            }
        },
        "/api/recommendation/user/{user}": {
            "get": {
                "description": "Gets recommendations for a given user",
                "produces": [
                    "application/json"
                ],
                "summary": "Get recommendations",
                "operationId": "get-recommendations",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user ID",
                        "name": "user",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Recommendation returned",
                        "schema": {
                            "$ref": "#/definitions/model.Recommendations"
                        }
                    },
                    "400": {
                        "description": "Invalid payload",
                        "schema": {
                            "$ref": "#/definitions/model.ErrResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/model.ErrResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.ErrResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "error": {
                    "type": "string"
                }
            }
        },
        "model.Item": {
            "type": "object",
            "properties": {
                "item": {
                    "type": "string"
                },
                "score": {
                    "type": "number"
                }
            }
        },
        "model.ItemProbability": {
            "type": "object",
            "properties": {
                "item": {
                    "type": "string"
                },
                "propability": {
                    "type": "number"
                },
                "user": {
                    "type": "string"
                }
            }
        },
        "model.Items": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Item"
                    }
                },
                "user": {
                    "type": "string"
                }
            }
        },
        "model.Rate": {
            "type": "object",
            "properties": {
                "item": {
                    "type": "string"
                },
                "score": {
                    "type": "number"
                },
                "user": {
                    "type": "string"
                }
            }
        },
        "model.Recommendation": {
            "type": "object",
            "properties": {
                "item": {
                    "type": "string"
                },
                "score": {
                    "type": "number"
                }
            }
        },
        "model.Recommendations": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Recommendation"
                    }
                },
                "user": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
